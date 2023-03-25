package retry

import (
	"context"
	"log"
	"time"
)

// Effector is a function interacting with a service
type Effector func(ctx context.Context) (string, error)

// Retry wraps an Effector function to provide retry logic.

// Accepts an int describing the maximum number of retry attempts and
// a time.Duration describing the interval between each retry attempt
func Retry(effector Effector, retries int, delay time.Duration) Effector {
	return func(ctx context.Context) (string, error) {
		for r := 0; ; r++ {
			response, err := effector(ctx)
			if err == nil || r >= retries {
				return response, err
			}

			log.Printf("Attempt %d failed; retrying in %v", r+1, delay)

			select {
			case <-time.After(delay):
			case <-ctx.Done():
				return "", ctx.Err()
			}
		}
	}
}
