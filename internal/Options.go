package internal

import "time"

type Options struct {
	Limit  int
	Window time.Duration

	Capacity       int
	RefillAmount   int
	RefillDuration time.Duration
}
