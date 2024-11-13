package ratelimiter

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

type FixedWindowRateLimiterJWT struct {
	client *redis.Client
	limit  int
	window time.Duration
}

func NewFixedWindowLimiterJWT(client *redis.Client, limit int, window time.Duration) *FixedWindowRateLimiterJWT {
	return &FixedWindowRateLimiterJWT{
		client: client,
		limit:  limit,
		window: window,
	}
}

func (rl *FixedWindowRateLimiterJWT) Allow(key interface{}) (bool, time.Duration, error) {
	token := fmt.Sprintf("%v", key)

	ctx := context.Background()

	count, err := rl.client.Incr(ctx, token).Result()
	if err != nil {
		return false, 0, fmt.Errorf("failed to increment count: %w", err)
	}

	if count == 1 {
		err := rl.client.Expire(ctx, token, rl.window).Err()
		if err != nil {
			return false, 0, fmt.Errorf("failed to set expiration: %w", err)
		}
	}

	if count > int64(rl.limit) {
		ttl, err := rl.client.TTL(ctx, token).Result()
		if err != nil {
			return false, 0, fmt.Errorf("failed to get TTL: %w", err)
		}
		return false, ttl, nil
	}

	return true, 0, nil
}
