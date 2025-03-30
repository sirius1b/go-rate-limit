package internal

import (
	"sync"
	"testing"
	"time"
)

func TestFixedWindowLimiter_Allow(t *testing.T) {
	limiter := NewFixedWindowLimiter(Options{
		Limit:  2,
		Window: time.Second,
	})

	key := "test_key"

	// First two requests should be allowed
	if !limiter.Allow(key) {
		t.Error("First request should be allowed")
	}
	if !limiter.Allow(key) {
		t.Error("Second request should be allowed")
	}

	// Third request should be denied
	if limiter.Allow(key) {
		t.Error("Third request should be denied")
	}

	// Wait for the window to pass
	time.Sleep(time.Second)

	// Request after the window should be allowed again
	if !limiter.Allow(key) {
		t.Error("Request after window should be allowed")
	}
}

func TestFixedWindowLimiter_Allow_Concurrency(t *testing.T) {
	limiter := NewFixedWindowLimiter(Options{
		Limit:  10,
		Window: time.Second,
	})

	key := "test_key"
	numRequests := 20
	var allowedCount int
	var wg sync.WaitGroup

	for i := 0; i < numRequests; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if limiter.Allow(key) {
				allowedCount++
			}
		}()
	}

	wg.Wait()

	if limiter.Token(key) > 0 {
		t.Errorf("Allowed more requests than the limit. Allowed: %d, Limit: %d", allowedCount, 10-limiter.Token(key))
	}
}

func TestFixedWindowLimiter_Wait(t *testing.T) {
	limiter := NewFixedWindowLimiter(Options{
		Limit:  1,
		Window: time.Millisecond * 100,
	})

	key := "test_key"

	// First request should be allowed
	if !limiter.Allow(key) {
		t.Error("First request should be allowed")
	}

	startTime := time.Now()
	// Second request should be block
	if !limiter.Wait(key) {
		t.Error("Second request should return true")
	}
	elapsedTime := time.Since(startTime)

	// Check the waiting time to be roughly within the window
	if elapsedTime < time.Millisecond*100 {
		t.Errorf("Wait time should be around the window %v, was: %v", limiter.window, elapsedTime)
	}

	// Allow should now return true because of Wait
	if !limiter.Allow(key) {
		t.Error("Should be allowed, as waited")
	}

}

func TestFixedWindowLimiter_Limit(t *testing.T) {
	limit := 5
	limiter := NewFixedWindowLimiter(Options{
		Limit:  limit,
		Window: time.Millisecond * 500,
	})

	if limiter.Rate() != 10 {
		t.Errorf("Rate should be %d, but was %f", 10, limiter.Rate())
	}
}

func TestFixedWindowLimiter_Token(t *testing.T) {
	limiter := NewFixedWindowLimiter(Options{
		Limit:  3,
		Window: time.Second,
	})

	key := "test_key"

	// Initial token count should be 0
	if limiter.Token(key) != 3 {
		t.Errorf("Initial token count should be 3, but was %d", limiter.Token(key))
	}

	// Allow one request
	limiter.Allow(key)

	// Token count should be 1
	if limiter.Token(key) != 2 {
		t.Errorf("Token count should be 2, but was %d", limiter.Token(key))
	}

	// Allow two more requests
	limiter.Allow(key)
	limiter.Allow(key)

	// Token count should be 3
	if limiter.Token(key) != 0 {
		t.Errorf("Token count should be 0, but was %d", limiter.Token(key))
	}

	// Wait for the window to pass
	time.Sleep(time.Second)

	// Token count should still be 3 because the startTime and limits are only updated inside the Allow function
	if limiter.Token(key) != 0 {
		t.Errorf("Token count should be 0 after waiting since it does not update on token call, but was %d", limiter.Token(key))
	}

	// Call Allow one more time to reset the window and counter
	limiter.Allow(key)

	// Token count should be 1
	if limiter.Token(key) != 2 {
		t.Errorf("Token count should be 2 after allow, but was %d", limiter.Token(key))
	}
}

func TestFixedWindowLimiter_MultipleKeys(t *testing.T) {
	limiter := NewFixedWindowLimiter(Options{
		Limit:  1,
		Window: time.Second,
	})

	key1 := "key1"
	key2 := "key2"

	if !limiter.Allow(key1) {
		t.Error("key1: First request should be allowed")
	}
	if limiter.Allow(key1) {
		t.Error("key1: Second request should be denied")
	}

	if !limiter.Allow(key2) {
		t.Error("key2: First request should be allowed")
	}
	if limiter.Allow(key2) {
		t.Error("key2: Second request should be denied")
	}
}
