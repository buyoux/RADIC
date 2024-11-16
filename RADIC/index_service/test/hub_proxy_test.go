package test

import (
	indexservice "RADIC/index_service"
	"fmt"
	"testing"
	"time"
)

func TestGetServiceEndpointsByProxy(t *testing.T) {
	const qps = 10 // qps限制为10
	proxy := indexservice.GetServiceHubProxy(etcdServers, 3, qps)

	endpoint := "127.0.0.1:5000"
	proxy.Regist(serviceName, endpoint, 0)
	defer proxy.UnRegist(serviceName, endpoint)
	endpoints := proxy.GetServiceEndpoints(serviceName)
	fmt.Printf("endpoints %v\n", endpoints)

	endpoint = "127.0.0.2:5000"
	proxy.Regist(serviceName, endpoint, 0)
	defer proxy.UnRegist(serviceName, endpoint)
	endpoints = proxy.GetServiceEndpoints(serviceName)
	fmt.Printf("endpoints %v\n", endpoints)

	endpoint = "127.0.0.3:5000"
	proxy.Regist(serviceName, endpoint, 0)
	defer proxy.UnRegist(serviceName, endpoint)
	endpoints = proxy.GetServiceEndpoints(serviceName)
	fmt.Printf("endpoints %v\n", endpoints)

	time.Sleep(1 * time.Second)  //暂停1秒钟，把令牌桶的容量打满
	for i := 0; i < qps+5; i++ { //桶里面有10个令牌，从第11次开始就拒绝访问了
		ep := proxy.GetServiceEndpoint(serviceName)
		fmt.Printf("%d endpoints %v\n", i, ep)
	}

	time.Sleep(1 * time.Second)  //暂停1秒钟，把令牌桶的容量打满
	for i := 0; i < qps+5; i++ { //桶里面有10个令牌，从第11次开始就拒绝访问了
		ep := proxy.GetServiceEndpoint(serviceName)
		fmt.Printf("%d endpoints %v\n", i, ep)
	}
}

// go test -v ./index_service/test -run=^TestGetServiceEndpointsByProxy$ -count=1