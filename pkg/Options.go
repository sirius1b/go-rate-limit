package pkg

import (
	"time"

	"github.com/sirius1b/go-rate-limit/internal"
)

type Options struct {
	Limit  int
	Window time.Duration

	Capacity       int
	RefillAmount   int
	RefillDuration time.Duration
}

func (o Options) toInternal() internal.Options {
	return internal.Options{
		Limit:  o.Limit,
		Window: o.Window,

		Capacity:       o.Capacity,
		RefillAmount:   o.RefillAmount,
		RefillDuration: o.RefillDuration,
	}
}
