package pkg

import (
	"errors"

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

type ABCD struct {
}

func Require(limiterType LimiterType, option Options) (IRateLimiter, error) {
	var limiter IRateLimiter
	switch limiterType {
	case FixedWindow:
		limiter = internal.NewFixedWindowLimiter(option.toInternal())
	default:
		panic("TODO")
	}
	return limiter, errors.New("unimplemented")
}
