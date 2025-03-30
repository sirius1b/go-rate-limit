package pkg

import (
	"errors"

	"github.com/sirius1b/go-rate-limit/internal"
)

type LimiterType int

const (
	FixedWindow LimiterType = iota
	TokenBucket
)

type IRateLimiter interface {
	Allow(string) bool
	Wait(string) bool
	Limit() int
	Token(string) int
}

func Require(limiterType LimiterType, option Options) (IRateLimiter, error) {
	var limiter IRateLimiter
	switch limiterType {
	case FixedWindow:
		limiter = internal.NewFixedWindowLimiter(option.toInternal())

	case TokenBucket:
		internal.NewTokenBucketLimiter(option.toInternal())
	default:
		return nil, errors.New("Umplemented")
	}
	return limiter, nil
}
