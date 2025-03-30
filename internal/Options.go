package internal

import "time"

type Options struct {
	Limit     int
	StartTime time.Time
	Window    time.Duration
	Extra     int
}
