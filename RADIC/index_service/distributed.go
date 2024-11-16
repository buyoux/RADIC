package indexservice

import (
	"RADIC/types"
	"RADIC/util"
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
)

type Sentinel struct {
	// 从Hub上获取IndexServiceWorker集合。可能是直接访问ServiceHub,也可能是代理
	hub IServiceHub
	// 与各个IndexServiceWorker建立的连接。把连接缓存起来，避免每次都重建连接
	connPool sync.Map
}

func NewSentinel(etcdServers []string) *Sentinel {
	return &Sentinel{
		// hub: GetServiceHub(etcdServers, 3),   // 直接使用ServiceHub，不走代理
		hub:      GetServiceHubProxy(etcdServers, 3, 100), // 走代理
		connPool: sync.Map{},
	}
}

// 传入endpoint节点，返回grpc连接，并缓存grpc连接
func (sentinel *Sentinel) GetGrpcConn(endpoint string) *grpc.ClientConn {
	if v, exists := sentinel.connPool.Load(endpoint); exists {
		conn := v.(*grpc.ClientConn)
		// 若连接不可用,则从连接缓存中删除
		if conn.GetState() == connectivity.TransientFailure || conn.GetState() == connectivity.Shutdown {
			util.Log.Printf("connection status to endpoint %s is %s", endpoint, conn.GetState().String())
			conn.Close()
			sentinel.connPool.Delete(endpoint)
		} else {
			// 缓存命中该连接，则直接返回
			return conn
		}
	}
	// 连接到服务端
	// 控制连接超时
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()
	conn, err := grpc.DialContext(
		ctx,
		endpoint,
		// Credential即使为空，也必须设置,grpc的加密措施
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		// grpc.Dial本身是异步连接，不会阻塞等待连接就绪。
		// 但设置了grpc.WithBlock()，则阻塞至连接可用后才返回
		grpc.WithBlock(),
	)
	if err != nil {
		util.Log.Printf("dial %s failed: %s", endpoint, err)
		return nil
	}
	util.Log.Printf("connect to grpc server %s", endpoint)
	// 先缓存grpc连接，再返回
	sentinel.connPool.Store(endpoint, conn)
	return conn
}

// 向集群中添加文档(如果已存在，会先删除)
func (sentinel *Sentinel) AddDoc(doc types.Document) (int, error) {
	// 根据负载均衡策略选择一台index worker，将doc添加到上面去
	endpoint := sentinel.hub.GetServiceEndpoint(INDEX_SERVICE)
	if len(endpoint) == 0 {
		return 0, fmt.Errorf("there is no alive index worker")
	}
	conn := sentinel.GetGrpcConn(endpoint)
	if conn == nil {
		return 0, fmt.Errorf("connect to worker %s failed", endpoint)
	}
	client := NewIndexServiceClient(conn)
	affected, err := client.AddDoc(context.Background(), &doc)
	if err != nil {
		return 0, err
	}
	util.Log.Printf("add %d doc to worker %s", affected.Count, endpoint)
	return int(affected.Count), nil
}

// 从集群上删除docId,返回成功删除的doc数(正常情况下不会超过1)
func (sentinel *Sentinel) DeleteDoc(docId string) int {
	endpoints := sentinel.hub.GetServiceEndpoints(INDEX_SERVICE)
	if len(endpoints) == 0 {
		return 0
	}
	var n int32
	wg := sync.WaitGroup{}
	wg.Add(len(endpoints))
	for _, endpoint := range endpoints {
		// 并行到各个IndexServiceWorker上把docId删除，正常情况下只有一个worker上有该doc
		go func(endpoint string) {
			defer wg.Done()
			conn := sentinel.GetGrpcConn(endpoint)
			if conn != nil {
				client := NewIndexServiceClient(conn)
				affected, err := client.DeleteDoc(context.Background(), &DocId{DocId: docId})
				if err != nil {
					util.Log.Printf("delete doc %s from worker %s failed: %s", docId, endpoint, err)
				} else {
					if affected.Count > 0 {
						atomic.AddInt32(&n, affected.Count)
						util.Log.Printf("delete doc %s from worker %s successful", docId, endpoint)
					}
				}

			}
		}(endpoint)
	}
	wg.Wait()
	return int(atomic.LoadInt32(&n))
}

func (sentinel *Sentinel) Search(query *types.TermQuery, onFlag uint64, offFlag uint64, orFlags []uint64) []*types.Document {
	endpoints := sentinel.hub.GetServiceEndpoints(INDEX_SERVICE)
	if len(endpoints) == 0 {
		return nil
	}
	docs := make([]*types.Document, 0, 1000)
	resultCh := make(chan *types.Document, 1000)
	wg := sync.WaitGroup{}
	wg.Add(len(endpoints))
	for _, endpoint := range endpoints {
		go func(endpoint string) {
			defer wg.Done()
			conn := sentinel.GetGrpcConn(endpoint)
			if conn != nil {
				client := NewIndexServiceClient(conn)
				result, err := client.Search(context.Background(), &SearchRequest{
					Query:   query,
					OnFlag:  onFlag,
					OffFlag: offFlag,
					OrFlags: orFlags,
				})
				if err != nil {
					util.Log.Printf("search from worker %s failed: %s", endpoint, err)
				} else {
					if len(result.Results) > 0 {
						util.Log.Printf("search %d doc from worker %s", len(result.Results), endpoint)
						for _, doc := range result.Results {
							resultCh <- doc
						}
					}
				}

			}
		}(endpoint)
	}

	receiveFinish := make(chan struct{})
	go func() {
		for {
			// ok为false跳出循环，需要同时满足 1.管道已经关闭 && 2.管道内无数据
			// 3.完成2之后，这里等待管道内剩余数据消费完毕才退出循环
			doc, ok := <-resultCh
			if !ok {
				break
			}
			docs = append(docs, doc)
		}
		receiveFinish <- struct{}{} // 4.完成3退出循环后，即管道内数据已消费完毕，解除receiveFinish阻塞
	}()
	wg.Wait()       // 1.通过wg.Wait()确认所有index worker都已向管道内写数据完毕(即生产者)
	close(resultCh) // 2.数据写入完毕后就关闭管道，等待管道内剩余数据消费完毕
	<-receiveFinish // 5.等待消费完毕，信号解除阻塞
	return docs
}

// 关闭各个grpc client connection，关闭etcd client connection
func (sentinel *Sentinel) Close() (err error) {
	sentinel.connPool.Range(func(key, value any) bool {
		conn := value.(*grpc.ClientConn)
		err = conn.Close()
		return true
	})
	sentinel.hub.Close()
	return
}
