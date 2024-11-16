package test

import (
	indexservice "RADIC/index_service"
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"
)

var (
	balancer  indexservice.LoadBalancer
	endpoints = []string{"127.0.0.1", "127.0.0.2", "127.0.0.3"}
)

func testLB(balancer indexservice.LoadBalancer) {
	// 开启100个协程并发使用balancer
	const P = 100
	const LOOP = 100
	selected := make(chan string, P*LOOP)
	wg := sync.WaitGroup{}
	wg.Add(P)
	for i := 0; i < P; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < LOOP; j++ {
				// 取出一个endpoint
				endpoint := balancer.Take(endpoints)
				selected <- endpoint
				// 模拟使用过程
				time.Sleep(time.Duration(rand.Intn(50)) * time.Millisecond)
			}
		}()
	}
	wg.Wait()
	close(selected)

	cm := make(map[string]int, len(endpoints))
	for {
		endpoint, ok := <-selected
		if !ok {
			break
		}
		value, ok := cm[endpoint]
		if ok {
			cm[endpoint] = value + 1
		} else {
			cm[endpoint] = 1
		}
	}
	for k, v := range cm {
		fmt.Println(k, v)
	}
}

func TestRandomSelect(t *testing.T) {
	balancer = new(indexservice.RandomSelect)
	testLB(balancer)
}

func TestRoundRobinSelect(t *testing.T) {
	balancer = new(indexservice.RoundRobin)
	testLB(balancer)
}

// go test -v .\index_service\test\ -run=^TestRandomSelect -count=1
// go test -v .\index_service\test\ -run=^TestRoundRobinSelect -count=1
