package main

import (
	"fmt"
	"time"

	"github.com/sirius1b/go-rate-limit/pkg"
)

func main() {
	// Configure a Sliding Window Log rate limiter that allows 3 requests per second.
	options := pkg.Options{
		Limit:  2,
		Window: time.Millisecond * 500,
	}

	limiter, err := pkg.Require(pkg.SlidingWindowLog, options)
	if err != nil {
		panic(err)
	}

	// Simulate multiple requests.
	for i := 0; i < 10; i++ {
		key := "user789" // Rate limit based on a unique key (e.g., user ID).
		allowed := limiter.Allow(key)
		fmt.Printf("Request %d: Allowed? %v, Rate: %.2f, Tokens: %d\n", i+1, allowed, limiter.Rate(), limiter.Token(key))

		if !allowed {
			fmt.Println("  Rate limited!")
		}

		time.Sleep(200 * time.Millisecond) // Send requests every 200ms.
	}

	// Wait and test again to see sliding window effect.
	time.Sleep(time.Second * 2)
	fmt.Println("--- Sliding Window Effect ---")

	for i := 0; i < 3; i++ {
		key := "user789" // Rate limit based on a unique key (e.g., user ID).
		allowed := limiter.Allow(key)
		fmt.Printf("Request %d: Allowed? %v, Rate: %.2f, Tokens: %d\n", i+1, allowed, limiter.Rate(), limiter.Token(key))

		if !allowed {
			fmt.Println("  Rate limited!")
		}

		time.Sleep(200 * time.Millisecond) // Send requests every 200ms.
	}
}
