package main

import (
	"fmt"
	"time"

	"github.com/sirius1b/go-rate-limit/pkg"
)

func main() {
	// Configure a Token Bucket rate limiter that allows 2 requests per second
	// and can accumulate a burst of up to 5 tokens.
	options := pkg.Options{
		Capacity:       5,
		RefillAmount:   1,
		RefillDuration: time.Second,
	}

	limiter, err := pkg.Require(pkg.TokenBucket, options)
	if err != nil {
		panic(err)
	}

	// Simulate multiple requests.
	for i := 0; i < 10; i++ {
		key := "user456" // Rate limit based on a unique key (e.g., user ID).
		allowed := limiter.Allow(key)
		fmt.Printf("Request %d: Allowed? %v, Rate: %.2f, Tokens: %d\n", i+1, allowed, limiter.Rate(), limiter.Token(key))

		if !allowed {
			fmt.Println("  Rate limited!")
		}

		time.Sleep(200 * time.Millisecond) // Send requests every 500ms
	}

	// Wait a bit to refill tokens.
	time.Sleep(time.Second * 2)
	fmt.Println("--- Refilling Tokens ---")

	for i := 0; i < 3; i++ {
		key := "user456" // Rate limit based on a unique key (e.g., user ID).
		allowed := limiter.Allow(key)
		fmt.Printf("Request %d: Allowed? %v, Rate: %.2f, Tokens: %d\n", i+1, allowed, limiter.Rate(), limiter.Token(key))

		if !allowed {
			fmt.Println("  Rate limited!")
		}

		time.Sleep(200 * time.Millisecond) // Send requests every 500ms
	}
}
