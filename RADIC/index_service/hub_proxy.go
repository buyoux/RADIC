package indexservice

import (
	"RADIC/util"
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	etcdv3 "go.etcd.io/etcd/client/v3"
	"golang.org/x/time/rate"
)

// 将etcd服务的方法抽象为接口，在需要使用ServiceHub或HubProxy时使用IServiceHub接口
// 这样就可以选择使用ServiceHub或HubProxy
// 代理HubProxy功能和ServiceHub一致，只是多了两个功能：缓存和限流

type IServiceHub interface {
	// 注册服务
	Regist(service string, endpoint string, leaseID etcdv3.LeaseID) (etcdv3.LeaseID, error)
	// 注销服务
	UnRegist(service string, endpoint string) error
	// 服务发现
	GetServiceEndpoints(service string) []string
	// 选择服务的一台endpoint，做负载均衡
	GetServiceEndpoint(service string) string
	Close()
}

// 代理模式。对ServiceHub做一层代理，想访问endpoints时需要通过代理
// 代理提供了两个功能：缓存和限流保护
type HubProxy struct {
	// hub           *ServiceHub
	// 匿名成员变量，这个匿名变量继承ServiceHub的方法，也可重写其方法
	// 没有名称直接使用其类型名称当成变量名称使用
	*ServiceHub
	endpointCache sync.Map // 维护每一个service下的所有servers
	limiter       *rate.Limiter
	loadBalancer  LoadBalancer
}

var (
	proxy     *HubProxy
	proxyOnce sync.Once
)

// HubProxy的构造函数，单例模式
// qps:一秒钟最多允许请求多少次
func GetServiceHubProxy(etcdServers []string, heartbeatFrequency int64, qps int) *HubProxy {
	if proxy == nil {
		proxyOnce.Do(func() {
			serviceHub := GetServiceHub(etcdServers, heartbeatFrequency)
			if serviceHub != nil {
				proxy = &HubProxy{
					ServiceHub:    serviceHub,
					endpointCache: sync.Map{},
					// time.Duration(1e9/qps)*time.Nanosecond: 每多少纳秒生成一个令牌
					// 1s = 1e9ns
					limiter:      rate.NewLimiter(rate.Every(time.Duration(1e9/qps)*time.Nanosecond), qps),
					loadBalancer: &RandomSelect{},
				}
			}
		})
	}
	return proxy
}

// // 注册服务
// func (proxy *HubProxy) Regist(service string, endpoint string, leaseID etcdv3.LeaseID) (
// 	etcdv3.LeaseID, error) {
// 	return proxy.hub.Regist(service, endpoint, leaseID)
// }

// // 注销服务
// func (proxy *HubProxy) UnRegist(service string, endpoint string) error {
// 	return proxy.hub.UnRegist(service, endpoint)
// }

// 监听etcd的数据变化，及时更新本地缓存
func (proxy *HubProxy) watchEndpointsOfService(service string) {
	if _, exists := proxy.ServiceHub.watched.LoadOrStore(service, true); exists {
		// 监听过了，不必重复监听
		return
	}
	ctx := context.Background()
	prefix := strings.TrimRight(SERVICE_ROOT_PATH, "/") + "/" + service + "/"
	// 根据前缀监听服务节点变化，每一个修改都会放入管道ch
	ch := proxy.ServiceHub.client.Watch(ctx, prefix, etcdv3.WithPrefix())
	util.Log.Printf("监听服务%s的节点变化", service)
	go func() {
		// 遍历管道。这是死循环，除非关闭管道退出循环
		for response := range ch {
			// 每次从ch里取出的是事件的集合
			for _, event := range response.Events {
				path := strings.Split(string(event.Kv.Key), "/")
				if len(path) > 2 {
					service := path[len(path)-2]
					// 与etcd进行一次全量同步
					endpoints := proxy.ServiceHub.GetServiceEndpoints(service)
					if len(endpoints) > 0 {
						// 查询etcd的结果放入本地缓存
						proxy.endpointCache.Store(service, endpoints)
					} else {
						// 该service下已经没有endpoint
						proxy.endpointCache.Delete(service)
					}
				}
			}
		}
	}()
}

// 服务发现
// 把第一次查询etcd的结果缓存起来，然后安装一个Watcher，仅etcd数据变化时更新本地缓存
// 这样可以减轻etcd的访问压力，同时加上限流保护
func (proxy *HubProxy) GetServiceEndpoints(service string) []string {
	// 阻塞，直到桶中有一个令牌或超时
	// ctx, cancel := context.WithTimeout(context.Background(), 100 * time.Millisecond)
	// defer cancel()
	// proxy.limiter.Wait(ctx)
	fmt.Println("HubProxy GetServiceEndpoints")
	if !proxy.limiter.Allow() {
		// 不阻塞，如果桶中没有1个令牌，则函数返回空，即没有可用的endpoints
		return nil
	}
	// 监听etcd的数据变化，及时更新本地缓存
	proxy.watchEndpointsOfService(service)
	// 若本地缓存直接命中，直接返回endpoints，不需要再走etcd
	if endpoints, exists := proxy.endpointCache.Load(service); exists {
		return endpoints.([]string)
	} else {
		endpoints := proxy.ServiceHub.GetServiceEndpoints(service)
		if len(endpoints) > 0 {
			// 查询etcd的结果放入本地缓存
			proxy.endpointCache.Store(service, endpoints)
		}
		return endpoints
	}
}

func (proxy *HubProxy) GetServiceEndpoint(service string) string {
	return proxy.loadBalancer.Take(proxy.GetServiceEndpoints(service))
}
