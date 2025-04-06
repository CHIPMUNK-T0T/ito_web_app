package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type RateLimiter struct {
	requests map[string][]time.Time
	mu       sync.RWMutex
}

func NewRateLimiter() *RateLimiter {
	limiter := &RateLimiter{
		requests: make(map[string][]time.Time),
	}

	// 古いリクエスト記録をクリーンアップ
	go func() {
		for {
			time.Sleep(time.Minute)
			limiter.cleanup()
		}
	}()

	return limiter
}

func (rl *RateLimiter) RateLimit(maxRequests int, window time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()

		rl.mu.Lock()
		defer rl.mu.Unlock()

		now := time.Now()
		if _, exists := rl.requests[ip]; !exists {
			rl.requests[ip] = []time.Time{now}
			c.Next()
			return
		}

		// 指定時間枠内のリクエストのみを保持
		var validRequests []time.Time
		for _, t := range rl.requests[ip] {
			if now.Sub(t) <= window {
				validRequests = append(validRequests, t)
			}
		}

		if len(validRequests) >= maxRequests {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "リクエスト制限を超えました",
			})
			c.Abort()
			return
		}

		rl.requests[ip] = append(validRequests, now)
		c.Next()
	}
}

func (rl *RateLimiter) cleanup() {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	for ip, times := range rl.requests {
		var validTimes []time.Time
		for _, t := range times {
			if now.Sub(t) <= time.Hour {
				validTimes = append(validTimes, t)
			}
		}
		if len(validTimes) > 0 {
			rl.requests[ip] = validTimes
		} else {
			delete(rl.requests, ip)
		}
	}
} 