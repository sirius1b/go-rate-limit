package tokenBucket

import (
	"sync"
	"testing"
	"time"

	internal "github.com/sirius1b/go-rate-limit/internal"
)

func TestTokenBucket_Allow(t *testing.T) {
	options := internal.Options{
		Capacity:       10,
		RefillDuration: time.Second,
		RefillAmount:   1,
	}
	tb := NewTokenBucketLimiter(options)

	key := "test_key"

	// Initially allow request
	if !tb.Allow(key) {
		t.Errorf("Expected initial request to be allowed")
	}

	// Deplete tokens
	for i := 0; i < 9; i++ {
		if !tb.Allow(key) {
			t.Errorf("Expected request %d to be allowed", i+1)
		}
	}

	// Should be rate limited
	if tb.Allow(key) {
		t.Errorf("Expected request to be rate limited")
	}

	// Wait for refill
	time.Sleep(time.Second)

	// Should allow again after refill
	if !tb.Allow(key) {
		t.Errorf("Expected request to be allowed after refill")
	}
}

func TestTokenBucket_Wait(t *testing.T) {
	options := internal.Options{
		Capacity:       1,
		RefillDuration: time.Millisecond * 100,
		RefillAmount:   1,
	}
	tb := NewTokenBucketLimiter(options)

	key := "test_key"

	// First wait should be immediate
	start := time.Now()
	tb.Wait(key)
	duration := time.Since(start)

	if duration > time.Millisecond*10 {
		t.Errorf("First wait took too long: %v", duration)
	}
	tb.Allow(key) // consuming token
	// Second wait should take refill duration
	start = time.Now()
	tb.Wait(key)
	duration = time.Since(start)

	if duration < time.Millisecond*90 {
		t.Errorf("Second wait took too short: %v", duration)
	}

	if duration > time.Millisecond*150 {
		t.Errorf("Second wait took too long: %v", duration)
	}
}

func TestTokenBucket_Limit(t *testing.T) {
	options := internal.Options{
		Capacity:       10,
		RefillDuration: time.Second,
		RefillAmount:   1,
	}
	tb := NewTokenBucketLimiter(options)

	limit := tb.Rate()
	if limit != 1 {
		t.Errorf("Expected limit to be 1, got %f", limit)
	}

	options = internal.Options{
		Capacity:       10,
		RefillDuration: time.Millisecond * 500,
		RefillAmount:   10,
	}
	tb = NewTokenBucketLimiter(options)
	limit = tb.Rate()

	if limit != 20 {
		t.Errorf("Expected limit to be 20, got %f", limit)
	}
}

func TestTokenBucket_Token(t *testing.T) {
	options := internal.Options{
		Capacity:       10,
		RefillDuration: time.Second,
		RefillAmount:   1,
	}
	tb := NewTokenBucketLimiter(options)

	key := "test_key"

	// Initial tokens should be capacity
	tokens := tb.Token(key)
	if tokens != 10 {
		t.Errorf("Expected initial tokens to be 10, got %d", tokens)
	}

	// Deplete tokens with Allow
	for i := 0; i < 5; i++ {
		tb.Allow(key)
	}

	// Check remaining tokens
	tokens = tb.Token(key)
	if tokens != 5 {
		t.Errorf("Expected tokens to be 5, got %d", tokens)
	}

	// Wait for refill
	time.Sleep(time.Second)

	// Check tokens after refill
	tokens = tb.Token(key)
	if tokens != 6 {
		t.Errorf("Expected tokens to be 6, got %d", tokens)
	}

	// Check token never exceeds capacity
	time.Sleep(time.Second * 5)
	tokens = tb.Token(key)
	if tokens != 10 {
		t.Errorf("Expected tokens to be 10, got %d", tokens)
	}
}

func TestTokenBucket_Concurrency(t *testing.T) {
	options := internal.Options{
		Capacity:       10,
		RefillDuration: time.Millisecond * 100,
		RefillAmount:   1,
	}
	tb := NewTokenBucketLimiter(options)

	key := "test_key"
	numGoroutines := 100

	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < 10; j++ {
				tb.Allow(key)
				time.Sleep(time.Millisecond * 5)
			}
		}()
	}

	wg.Wait()

	// Check remaining tokens (approximately)
	tokens := tb.Token(key)
	if tokens < 0 || tokens > 10 {
		t.Errorf("Expected tokens to be between 0 and 10, got %d", tokens)
	}
}
