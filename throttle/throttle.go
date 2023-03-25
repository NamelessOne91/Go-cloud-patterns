package throttle

import (
	"context"
	"errors"
	"sync"
	"time"
)

var errThrottling = errors.New("too many calls")

// Effector is a function interacting with a service
type Effector func(context.Context) (string, error)

// Throttle wraps an Effector function to provide rate-limiting logic.
//
// It uses the token bucket strategy: a function call consumes one (or more) token from the bucket,
// which then refills at a fixed rate.
// When there are not enough tokens left, a custom throttling strategy is applied
func Throttle(e Effector, max uint, refill uint, d time.Duration) Effector {
	// token bucket
	var tokens = max
	var once sync.Once

	return func(ctx context.Context) (string, error) {
		if ctx.Err() != nil {
			return "", ctx.Err()
		}

		once.Do(func() {
			ticker := time.NewTicker(d)

			go func() {
				defer ticker.Stop()

				for {
					select {
					case <-ctx.Done():
						return
					// bucket refill
					case <-ticker.C:
						t := tokens + refill
						if t > max {
							t = max
						}
						tokens = t
					}
				}
			}()
		})
		// can also return the results of the last function call
		// or use a queue to retry calls later
		if tokens <= 0 {
			return "", errThrottling
		}
		tokens--

		return e(ctx)
	}
}
