package pkg

import (
	fw "github.com/sirius1b/go-rate-limit/internal/fixedWindow"
	sw "github.com/sirius1b/go-rate-limit/internal/slidingWindow"
	tb "github.com/sirius1b/go-rate-limit/internal/tokenBucket"
)

type LimiterType int

const (
	FixedWindow LimiterType = iota
	TokenBucket
	SlidingWindowLog
)

type IRateLimiter interface {
	Allow(string) bool
	Wait(string) bool
	Rate() float64
	Token(string) int
}

func Require(limiterType LimiterType, option Options) (IRateLimiter, error) {
	var limiter IRateLimiter
	switch limiterType {
	case FixedWindow:
		limiter = fw.NewFixedWindowLimiter(option.toInternal())
	case TokenBucket:
		limiter = tb.NewTokenBucketLimiter(option.toInternal())
	case SlidingWindowLog:
		limiter = sw.NewSlidingWindowLimiter(option.toInternal())
	}

	return limiter, nil
}
