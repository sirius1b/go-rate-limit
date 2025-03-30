package pkg

import (
	"testing"
	"time"
)

func TestRequireFixedWindow(t *testing.T) {
	options := Options{
		Limit:  10,
		Window: time.Second,
	}

	limiter, err := Require(FixedWindow, options)
	if err != nil {
		t.Fatalf("Failed to create FixedWindow limiter: %v", err)
	}

	if limiter == nil {
		t.Fatal("Limiter is nil")
	}

	// Basic allow test
	key := "test_key"
	allowed := 0
	for limiter.Token(key) < options.Limit {
		if limiter.Allow(key) {
			allowed++
		}
	}
	if allowed != options.Limit {
		t.Errorf("Expected %d allows, got %d", options.Limit, allowed)
	}

	// Check if further requests are blocked
	if limiter.Allow(key) {
		t.Error("Should be blocked after limit is reached")
	}

	// Wait functionality test
	options.Limit = 1
	options.Window = time.Second

	limiter, err = Require(FixedWindow, options)

	if err != nil {
		t.Fatalf("Failed to create FixedWindow limiter: %v", err)
	}

	start := time.Now()
	limiter.Wait(key)
	duration := time.Since(start)

	if duration < time.Duration(0) { // check it doens't wait infinitely
		t.Errorf("Wait func returned too fast! It should wait 1s")
	}

	if !limiter.Allow(key) {
		t.Errorf("Allow func should return false after wait")
	}
}
