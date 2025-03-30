package pkg

import (
	"github.com/sirius1b/go-rate-limit/internal"
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
		limiter = internal.NewFixedWindowLimiter(option.toInternal())
	case TokenBucket:
		limiter = internal.NewTokenBucketLimiter(option.toInternal())
	case SlidingWindowLog:
		limiter = internal.NewSlidingWindowLimiter(option.toInternal())
	}

	return limiter, nil
}
