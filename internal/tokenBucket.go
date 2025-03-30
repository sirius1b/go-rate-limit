package internal

import (
	"sync"
	"time"
)

type TokenBucket struct {
	capacity       int
	refillDuration time.Duration
	refillAmount   int
	tokens         map[string]int
	lastRefills    map[string]time.Time
	mus            map[string]*sync.Mutex
	globalMu       *sync.Mutex
}

func (s *TokenBucket) getMutex(key string) *sync.Mutex {
	s.globalMu.Lock()
	defer s.globalMu.Unlock()
	if _, exists := s.mus[key]; !exists {
		s.mus[key] = &sync.Mutex{}
	}
	return s.mus[key]
}

func NewTokenBucketLimiter(option Options) *TokenBucket {
	return &TokenBucket{
		capacity:       option.Capacity,
		refillDuration: option.RefillDuration,
		refillAmount:   option.RefillAmount,
		tokens:         make(map[string]int),
		lastRefills:    make(map[string]time.Time),
		mus:            make(map[string]*sync.Mutex),
		globalMu:       &sync.Mutex{},
	}
}

func (f *TokenBucket) Allow(key string) bool {
	mu := f.getMutex(key)
	mu.Lock()
	defer mu.Unlock()

	now := time.Now()

	f.refill(key, now)

	if f.tokens[key] > 0 {
		f.tokens[key]--
		return true
	}

	return false
}

func (f *TokenBucket) Wait(key string) bool {
	mu := f.getMutex(key)
	mu.Lock()
	defer mu.Unlock()

	now := time.Now()

	f.refill(key, now)

	elapsed := now.Sub(f.lastRefills[key])
	if f.tokens[key] <= 0 {
		sleepTime := f.refillDuration - elapsed

		time.Sleep(sleepTime)
	}

	return true
}

func (f *TokenBucket) Rate() float64 {
	return float64(f.refillAmount) / float64(f.refillDuration.Seconds())
}

func (f *TokenBucket) Token(key string) int {
	mu := f.getMutex(key)
	mu.Lock()
	defer mu.Unlock()

	f.refill(key, time.Now())

	return f.tokens[key]
}

func (f *TokenBucket) refill(key string, now time.Time) {

	then := f.lastRefills[key]
	elapsed := now.Sub(then)
	if elapsed >= f.refillDuration {
		refills := int(elapsed / f.refillDuration)
		f.tokens[key] = min(f.capacity, f.tokens[key]+refills*f.refillAmount)
		f.lastRefills[key] = now
	}

}
