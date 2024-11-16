package indexservice

import (
	"RADIC/types"
	"RADIC/util"
	"context"
	"fmt"
	"strconv"
	"time"
)

const (
	INDEX_SERVICE = "index_service"
)

type IndexServiceWorker struct {
	Indexer *Indexer // Indexer将正排索引和倒排索引放到一起
	// 服务注册相关的配置
	hub      *ServiceHub
	selfAddr string
}

// 初始化索引
func (service *IndexServiceWorker) Init(DocNumEstimate int, dbtype int, DataDir string) error {
	service.Indexer = new(Indexer)
	return service.Indexer.Init(DocNumEstimate, dbtype, DataDir)
}

func (service *IndexServiceWorker) Regist(etcdServers []string, servicePort int) error {
	// 向注册中心注册自己
	if len(etcdServers) > 0 {
		if servicePort <= 1024 {
			return fmt.Errorf("invalid listen port %d, should more than 1024", servicePort)
		}
		selfLocalIP, err := util.GetLocalIP()
		if err != nil {
			panic(err)
		}
		// 单机模拟分布式时，把selfLocalIP写死为127.0.0.1
		selfLocalIP = "127.0.0.1"
		service.selfAddr = selfLocalIP + ":" + strconv.Itoa(servicePort)
		var heartBeat int64 = 3                      // 每隔3秒上报一次心跳，即周期性续约
		hub := GetServiceHub(etcdServers, heartBeat) // 单例模式
		leaseId, err := hub.Regist(INDEX_SERVICE, service.selfAddr, 0)
		if err != nil {
			panic(err)
		}
		service.hub = hub
		// 周期性注册自己（续约）
		go func() {
			for {
				hub.Regist(INDEX_SERVICE, service.selfAddr, leaseId)
				// 根据心跳时长设置休眠时间
				time.Sleep(time.Duration(heartBeat)*time.Second - 100*time.Millisecond)
			}
		}()
	}
	return nil
}

// 系统重启时，直接从索引文件里加载数据
func (service *IndexServiceWorker) LoadFromIndexFile() int {
	return service.Indexer.LoadFromIndexFile()
}

// 关闭索引
func (service *IndexServiceWorker) Close() error {
	if service.hub != nil {
		service.hub.UnRegist(INDEX_SERVICE, service.selfAddr)
	}
	return service.Indexer.Close()
}

// 从索引上删除文档
func (service *IndexServiceWorker) DeleteDoc(ctx context.Context, docId *DocId) (*AffectedCount, error) {
	return &AffectedCount{int32(service.Indexer.DeleteDoc(docId.DocId))}, nil
}

// 向索引中添加文档(如果已存在，会先删除)
func (service *IndexServiceWorker) AddDoc(ctx context.Context, doc *types.Document) (*AffectedCount, error) {
	n, err := service.Indexer.AddDoc(*doc)
	return &AffectedCount{int32(n)}, err
}

// 检索，返回文档列表
func (service *IndexServiceWorker) Search(ctx context.Context, request *SearchRequest) (*SearchResult, error) {
	result := service.Indexer.Search(request.Query, request.OnFlag, request.OffFlag, request.OrFlags)
	return &SearchResult{Results: result}, nil
}

func (service *IndexServiceWorker) Count(ctx context.Context, request *SearchRequest) (*AffectedCount, error) {
	return &AffectedCount{int32(service.Indexer.Count())}, nil
}
