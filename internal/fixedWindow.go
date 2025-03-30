package internal

import (
	"sync"
	"time"
)

type FixedWindowLimiter struct {
	limit     int
	limits    map[string]int
	startTime map[string]time.Time
	window    time.Duration
	mus       map[string]*sync.Mutex
	globalMu  *sync.Mutex
}

func (s *FixedWindowLimiter) getMutex(key string) *sync.Mutex {
	s.globalMu.Lock()
	defer s.globalMu.Unlock()
	if _, exists := s.mus[key]; !exists {
		s.mus[key] = &sync.Mutex{}
	}
	return s.mus[key]
}

func NewFixedWindowLimiter(option Options) *FixedWindowLimiter {
	return &FixedWindowLimiter{
		limit:     option.Limit,
		startTime: make(map[string]time.Time),
		window:    option.Window,
		limits:    make(map[string]int),
		mus:       make(map[string]*sync.Mutex),
		globalMu:  &sync.Mutex{},
	}
}

func (f *FixedWindowLimiter) Allow(key string) bool {
	mu := f.getMutex(key)
	mu.Lock()
	defer mu.Unlock()

	now := time.Now()

	if _, exists := f.limits[key]; !exists {
		f.limits[key] = 0
		f.startTime[key] = now
	}

	if now.Sub(f.startTime[key]) >= f.window {
		f.startTime[key] = now
		f.limits[key] = 0
	}

	if f.limits[key] < f.limit {
		f.limits[key]++
		return true
	}

	return false
}

func (f FixedWindowLimiter) Wait(key string) bool {
	mu := f.getMutex(key)
	mu.Lock()
	defer mu.Unlock()

	now := time.Now()

	if _, exists := f.limits[key]; !exists {
		f.limits[key] = 0
		f.startTime[key] = now
	}

	elapsed := now.Sub(f.startTime[key])

	if f.limits[key] >= f.limit {
		sleepTime := f.window - elapsed
		time.Sleep(sleepTime)

		f.startTime[key] = time.Now()
		f.limits[key] = 0
	}

	return true
}

func (f FixedWindowLimiter) Limit() int {
	return f.limit
}

func (f FixedWindowLimiter) Token(key string) int {
	mu := f.getMutex(key)
	mu.Lock()
	defer mu.Unlock()

	now := time.Now()

	if _, exists := f.limits[key]; !exists {
		f.limits[key] = 0
		f.startTime[key] = now
	}

	return f.limits[key]
}
