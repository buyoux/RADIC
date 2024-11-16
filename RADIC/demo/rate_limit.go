package course

import (
	"context"
	"sync/atomic"
	"time"

	"golang.org/x/time/rate"
)

// 接口限流，令牌桶限流算法

var TotalQuery int32

// 模拟常规后端接口执行
func Handler() {
	atomic.AddInt32(&TotalQuery, 1)
	time.Sleep(50 * time.Millisecond)
}

func CallHandler() {
	// 每隔100ms生成一个令牌
	limiter := rate.NewLimiter(rate.Every(100*time.Millisecond), 10)
	n := 3
	for {
		// 取令牌方案1
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()
		// 阻塞，直到桶中有N个令牌，如果超时返回错误。N=1等价于limiter.Wait()
		// 令牌不够也不会取剩下的，而是到时间报超时错误。如果到时之前令牌够了，就直接Handler
		// 若不设置超时时间，就会一直阻塞等待下去
		if err := limiter.WaitN(ctx, n); err == nil {
			Handler()
		}

		// 取令牌方案2
		// 当前时间桶中是否至少还有n个令牌，如果有则返回true
		if limiter.AllowN(time.Now(), n) {
			Handler()
		}

		// 取令牌方案3
		// reserve.Delay()告诉你还需要多久才能有充足的令牌
		reserve := limiter.ReserveN(time.Now(), n)
		time.Sleep(reserve.Delay())
		Handler()
	}
}
