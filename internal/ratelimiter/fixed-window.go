package ratelimiter

import (
	"fmt"
	"sync"
	"time"
)

// TODO: change to redis for multiple purpose
type FixedWindowRateLimiter struct {
	sync.RWMutex
	clients map[string]int
	limit   int
	window  time.Duration
}

func NewFixedWindowLimiter(limit int, window time.Duration) *FixedWindowRateLimiter {
	return &FixedWindowRateLimiter{
		clients: make(map[string]int),
		limit:   limit,
		window:  window,
	}
}

func (rl *FixedWindowRateLimiter) Allow(key interface{}) (bool, time.Duration, error) {
	ip, ok := key.(string)
	if !ok {
		return false, 0, fmt.Errorf("invalid key type")
	}

	rl.RLock()
	count, exists := rl.clients[ip]
	rl.RUnlock()

	if !exists || count < rl.limit {
		rl.Lock()
		defer rl.Unlock()
		if !exists {
			go rl.resetCount(ip)
		}

		rl.clients[ip]++
		return true, 0, nil
	}

	return false, rl.window, nil
}

func (rl *FixedWindowRateLimiter) resetCount(ip string) {
	time.Sleep(rl.window)
	rl.Lock()
	defer rl.Unlock()
	delete(rl.clients, ip)
}
