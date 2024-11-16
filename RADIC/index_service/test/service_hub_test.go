package test

import (
	indexservice "RADIC/index_service"
	"fmt"
	"testing"
)

var serviceName = "test_service"
var etcdServers = []string{"127.0.0.1:2379"}

func TestGetServiceEndpoints(t *testing.T) {
	fmt.Println(etcdServers)
	hub := indexservice.GetServiceHub(etcdServers, 3)

	endpoint := "127.0.0.1:5000"
	hub.Regist(serviceName, endpoint, 0)
	defer hub.UnRegist(serviceName, endpoint)
	endpoints := hub.GetServiceEndpoints(serviceName)
	fmt.Printf("endpoints %v\n", endpoints)

	endpoint = "127.0.0.2:5000"
	hub.Regist(serviceName, endpoint, 0)
	defer hub.UnRegist(serviceName, endpoint)
	endpoints = hub.GetServiceEndpoints(serviceName)
	fmt.Printf("endpoints %v\n", endpoints)

	endpoint = "127.0.0.3:5000"
	hub.Regist(serviceName, endpoint, 0)
	defer hub.UnRegist(serviceName, endpoint)
	endpoints = hub.GetServiceEndpoints(serviceName)
	fmt.Printf("endpoints %v\n", endpoints)
}

// go test -v .\index_service\test\ -run=^TestGetServiceEndpoints$ -count=1
