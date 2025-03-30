package internal

import (
	"sync"
	"testing"
	"time"
)

func TestSlidingWindow_Allow(t *testing.T) {
	option := Options{
		Limit:  3,
		Window: time.Second,
	}
	limiter := NewSlidingWindowLimiter(option)

	key := "test_key"

	// Allow 3 requests within the window
	if !limiter.Allow(key) {
		t.Error("Expected Allow to return true")
	}

	if !limiter.Allow(key) {
		t.Error("Expected Allow to return true")
	}

	if !limiter.Allow(key) {
		t.Error("Expected Allow to return true")
	}

	// The 4th request should be denied
	if limiter.Allow(key) {
		t.Error("Expected Allow to return false")
	}

	// Wait for the window to expire
	time.Sleep(time.Second)

	// Should be able to allow again
	if !limiter.Allow(key) {
		t.Error("Expected Allow to return true after window expiry")
	}

	// Test with multiple keys
	key2 := "test_key2"
	if !limiter.Allow(key2) {
		t.Error("Expected Allow to return true for a new key")
	}
}

func TestSlidingWindow_Rate(t *testing.T) {
	option := Options{
		Limit:  10,
		Window: time.Second * 2,
	}
	limiter := NewSlidingWindowLimiter(option)

	expectedRate := float64(10) / float64(2)
	actualRate := limiter.Rate()

	if actualRate != expectedRate {
		t.Errorf("Expected rate %f, but got %f", expectedRate, actualRate)
	}
}

func TestSlidingWindow_Token(t *testing.T) {
	option := Options{
		Limit:  5,
		Window: time.Second,
	}
	limiter := NewSlidingWindowLimiter(option)

	key := "token_key"

	// Initially, all tokens should be available
	if limiter.Token(key) != 5 {
		t.Errorf("Expected 5 tokens, but got %d", limiter.Token(key))
	}

	limiter.Allow(key)
	limiter.Allow(key)

	// After 2 requests, 3 tokens should be available
	if limiter.Token(key) != 3 {
		t.Errorf("Expected 3 tokens, but got %d", limiter.Token(key))
	}

	// Fill the bucket
	limiter.Allow(key)
	limiter.Allow(key)
	limiter.Allow(key)

	// Verify no tokens remain
	if limiter.Token(key) != 0 {
		t.Errorf("Expected 0 tokens, but got %d", limiter.Token(key))
	}
}

func TestSlidingWindow_Concurrency(t *testing.T) {
	option := Options{
		Limit:  100,
		Window: time.Second,
	}
	limiter := NewSlidingWindowLimiter(option)

	key := "concurrent_key"

	var wg sync.WaitGroup
	numGoroutines := 200

	allowedCount := 0
	var mu sync.Mutex // Protects allowedCount

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if limiter.Allow(key) {
				mu.Lock()
				allowedCount++
				mu.Unlock()
			}
		}()
	}

	wg.Wait()

	if allowedCount > option.Limit {
		t.Errorf("Allowed more requests than the limit. Allowed: %d, Limit: %d", allowedCount, option.Limit)
	}
}

func TestSlidingWindow_RemoveSamples(t *testing.T) {
	option := Options{
		Limit:  5,
		Window: time.Second,
	}
	limiter := NewSlidingWindowLimiter(option)
	key := "remove_key"

	now := time.Now()
	past1 := now.Add(-1 * time.Second)
	past2 := now.Add(-2 * time.Second)
	past3 := now.Add(-3 * time.Second)

	// Manually populate timeLogs with outdated timestamps
	limiter.timeLogs[key] = []time.Time{past3, past2, past1}

	limiter.removeSamples(key, now)

	// Only past1 should remain after removeSamples is called
	if len(limiter.timeLogs[key]) != 1 {
		t.Fatalf("Expected 1 log, got %d", len(limiter.timeLogs[key]))
	}

	// Test with no logs
	key2 := "remove_key2"
	limiter.removeSamples(key2, now)

	// timeLogs[key2] should be initialized with an empty slice of Time{}
	if len(limiter.timeLogs[key2]) != 0 {
		t.Fatalf("Expected timeLogs to be initialized. Got length %d", len(limiter.timeLogs[key2]))
	}
}

func TestSlidingWindow_Wait(t *testing.T) {
	option := Options{
		Limit:  1,
		Window: time.Second,
	}
	limiter := NewSlidingWindowLimiter(option)
	key := "wait_key"

	// First request should be allowed immediately
	if !limiter.Allow(key) {
		t.Error("Expected Allow to return true")
	}

	// next request would need to wait
	start := time.Now()

	limiter.Wait(key)

	duration := time.Since(start)

	if float64(time.Second)-float64(duration) < float64(time.Millisecond*10) {
		t.Errorf("Wait time less than 1 seconds")
	}
}
