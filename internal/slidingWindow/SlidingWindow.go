package slidingWindow

import (
	"sync"
	"time"

	internal "github.com/sirius1b/go-rate-limit/internal"
)

type SlidingWindow struct {
	timeLogs map[string][]time.Time
	limit    int
	window   time.Duration

	mus      map[string]*sync.Mutex
	globalMu *sync.Mutex
}

func NewSlidingWindowLimiter(option internal.Options) *SlidingWindow {
	return &SlidingWindow{

		limit:    option.Limit,
		window:   option.Window,
		timeLogs: make(map[string][]time.Time),

		mus:      make(map[string]*sync.Mutex),
		globalMu: &sync.Mutex{},
	}
}

func (f *SlidingWindow) Allow(key string) bool {
	mu := f.getMutex(key)
	mu.Lock()
	defer mu.Unlock()
	now := time.Now()

	f.removeSamples(key, now)

	if len(f.timeLogs[key]) < f.limit {
		f.timeLogs[key] = append(f.timeLogs[key], now)
		return true
	}

	return false
}

func (f *SlidingWindow) Wait(key string) bool {
	mu := f.getMutex(key)
	mu.Lock()
	defer mu.Unlock()

	logs, exists := f.timeLogs[key]
	now := time.Now()

	if exists {
		lastLog := logs[0]
		expiryTime := lastLog.Add(f.window)

		if now.After(expiryTime) {
			return true
		}

		sleepTime := time.Since(expiryTime)
		time.Sleep(sleepTime)
	}

	return true
}

func (f *SlidingWindow) Rate() float64 {
	return float64(f.limit) / float64(f.window.Seconds())
}

func (f *SlidingWindow) Token(key string) int {
	mu := f.getMutex(key)
	mu.Lock()
	defer mu.Unlock()

	return f.limit - len(f.timeLogs[key])
}

func (s *SlidingWindow) getMutex(key string) *sync.Mutex {
	s.globalMu.Lock()
	defer s.globalMu.Unlock()
	if _, exists := s.mus[key]; !exists {
		s.mus[key] = &sync.Mutex{}
	}
	return s.mus[key]
}

func (r *SlidingWindow) removeSamples(key string, now time.Time) {

	_, exists := r.timeLogs[key]
	if !exists {
		r.timeLogs[key] = make([]time.Time, 0)
	} else {
		threshold := now.Add(-r.window)
		removeIndex := -1
		for index := len(r.timeLogs[key]) - 1; index >= 0; index-- {
			if threshold.After(r.timeLogs[key][index]) {
				removeIndex = index
				break
			}
		}
		if removeIndex != -1 {
			r.timeLogs[key] = r.timeLogs[key][removeIndex+1:]
		}
	}
}
