package pkg

import (
	"github.com/sirius1b/go-rate-limit/internal"
)

type LimiterType int

const (
	FixedWindow LimiterType = iota
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
	default:
		panic("Unimplemnet") // FUTURE TODO:
	}
	return limiter, nil
}
