package indexservice

import (
	"math/rand"
	"sync/atomic"
)

// 策略模式。完成同一个任务可以有多种不同的实现方案

type LoadBalancer interface {
	Take([]string) string
}

// 负载均衡算法--轮询法
type RoundRobin struct {
	acc int64
}

func (b *RoundRobin) Take(endpoints []string) string {
	if len(endpoints) == 0 {
		return ""
	}
	n := atomic.AddInt64(&b.acc, 1) // Take()需要支持并发调用，所以使用原子操作
	index := int(n % int64(len(endpoints)))
	return endpoints[index]
}

// 负载均衡算法--随机法
type RandomSelect struct{}

func (b *RandomSelect) Take(endpoints []string) string {
	if len(endpoints) == 0 {
		return ""
	}
	// 随机选择
	index := rand.Intn(len(endpoints))
	return endpoints[index]
}
