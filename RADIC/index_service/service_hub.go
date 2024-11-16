package indexservice

import (
	"RADIC/util"
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"go.etcd.io/etcd/api/v3/v3rpc/rpctypes"
	etcdv3 "go.etcd.io/etcd/client/v3"
)

const (
	// etcd key的前缀
	SERVICE_ROOT_PATH = "/RADIC/index"
)

type ServiceHub struct {
	client             *etcdv3.Client
	heartbeatFrequency int64 // server每隔几秒钟不动向中心上报一次心跳(即续一次租约)
	watched            sync.Map
	loadBalancer       LoadBalancer // 策略模式。完成同一个任务可以有多种不同的实现方案
}

var (
	serviceHub *ServiceHub // 该全局变量包外不可见，包外使用通过GetServiceHub()获得
	hubOnce    sync.Once   // 单例模式，需要用到once保证只执行一次
)

// ServiceHub的构造函数，单例模式
func GetServiceHub(etcdServers []string, heartbeatFrequency int64) *ServiceHub {
	if serviceHub == nil {
		hubOnce.Do(func() {
			if client, err := etcdv3.New(
				etcdv3.Config{
					Endpoints:   etcdServers,
					DialTimeout: 3 * time.Second,
				},
			); err != nil {
				// 发生log.Fatal时go进程会直接结束退出
				util.Log.Fatalf("connect to etcd server failed: %v", err)
			} else {
				serviceHub = &ServiceHub{
					client:             client,
					heartbeatFrequency: heartbeatFrequency,
					loadBalancer:       &RoundRobin{},
				}
			}
		})
	}
	return serviceHub
}

// 注册服务。第一次注册向etcd写一个key，后续注册仅仅是在续约
// service      微服务的名称
// endpoint     微服务server的地址
// leaseID      租约ID，第一次注册时置为0即可

func (hub *ServiceHub) Regist(service string, endpoint string, leaseID etcdv3.LeaseID) (etcdv3.LeaseID, error) {
	ctx := context.Background()
	if leaseID <= 0 {
		// 创建一个租约，有效期为heartbeatFrequency秒
		if lease, err := hub.client.Grant(ctx, hub.heartbeatFrequency); err != nil {
			util.Log.Printf("创建租约失败：%v", err)
			return 0, err
		} else {
			key := strings.TrimRight(SERVICE_ROOT_PATH, "/") + "/" + service + "/" + endpoint
			// 服务注册,若没有etcdv3.WithLease(leaseID)，即永久有效
			if _, err = hub.client.Put(ctx, key, "", etcdv3.WithLease(leaseID)); err != nil {
				util.Log.Printf("注册服务%s对应的节点%s失败:%v", service, endpoint, err)
				return lease.ID, err
			} else {
				util.Log.Printf("注册服务%s对应的节点%s成功", service, endpoint)
				return lease.ID, nil
			}
		}
	} else {
		// 续租,续约一次，到期继续续约
		if _, err := hub.client.KeepAliveOnce(ctx, leaseID); err == rpctypes.ErrLeaseNotFound {
			// 容错处理，租约找不到，走注册流程（将leaseID置为0）
			return hub.Regist(service, endpoint, 0)
		} else if err != nil {
			util.Log.Printf("续约失败:%v", err)
			return 0, err
		} else {
			util.Log.Printf("服务%s对应的节点%s续约成功", service, endpoint)
			return leaseID, nil
		}
	}
}

// 注销服务
func (hub *ServiceHub) UnRegist(service string, endpoint string) error {
	ctx := context.Background()
	key := strings.TrimRight(SERVICE_ROOT_PATH, "/") + "/" + service + "/" + endpoint
	if _, err := hub.client.Delete(ctx, key); err != nil {
		util.Log.Printf("注销服务%s对应的节点%s失败:%v", service, endpoint, err)
		return err
	} else {
		util.Log.Printf("注销服务%s对应的节点%s", service, endpoint)
		return nil
	}
}

// 服务发现。client每次进行RPC调用之前都查询etcd，
// 获取server集合，然后采用负载均衡算法选择一台server
// 或者可以将负载均衡的功能放在注册中心，即放到getServiceEndpoints函数中，让他只返回一个server
func (hub *ServiceHub) GetServiceEndpoints(service string) []string {
	fmt.Println("ServiceHub GetServiceEndpoints")
	ctx := context.Background()
	prefix := strings.TrimRight(SERVICE_ROOT_PATH, "/") + "/" + service + "/"
	if resp, err := hub.client.Get(ctx, prefix, etcdv3.WithPrefix()); err != nil {
		util.Log.Printf("获取服务%s的节点失败:%v", service, err)
		return nil
	} else {
		endpoints := make([]string, 0, len(resp.Kvs))
		for _, kv := range resp.Kvs {
			// 只需要使用到key，不需要value
			path := strings.Split(string(kv.Key), "/")
			fmt.Println(kv.String(), path[len(path)-1], "120")
			endpoints = append(endpoints, path[len(path)-1])
		}
		util.Log.Printf("刷新%s服务对应的server -- %v\n", service, endpoints)
		return endpoints
	}
}

func (hub *ServiceHub) GetServiceEndpoint(service string) string {
	return hub.loadBalancer.Take(hub.GetServiceEndpoints(service))
}

// 关闭etcd client connection
func (hub *ServiceHub) Close() {
	hub.client.Close()
}

// go test -v .\index_service\test\ -run=^TestGetServiceEndpointsByProxy$ -count=1
