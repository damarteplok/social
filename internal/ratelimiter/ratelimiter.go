package ratelimiter

import "time"

type Limiter interface {
	Allow(key interface{}) (bool, time.Duration, error)
}

type Config struct {
	RequestPerTimeFrame int
	TimeFrame           time.Duration
	Enabled             bool
}
