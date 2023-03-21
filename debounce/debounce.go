package debounce

import (
	"context"
	"sync"
	"time"
)

// Circuit represents a function interacting with an upstream service;
// it should include an error in its return list
type Circuit func(context.Context) (string, error)

// DebounceFirst is a function-first implementation that wraps a Circuit function
// forcing overlapping calls to wait untill the result is cached and guarantees circuit
// is called exactly once, at the beginning of a cluster
func DebounceFirst(circuit Circuit, d time.Duration) Circuit {
	var threshold time.Time
	var result string
	var err error
	var m sync.Mutex

	return func(ctx context.Context) (string, error) {
		// provides thread safety for the entire function
		m.Lock()
		defer func() {
			// reset the time at which a cluster ends
			threshold = time.Now().Add(d)
			m.Unlock()
		}()

		if time.Now().Before(threshold) {
			return result, err
		}

		result, err := circuit(ctx)
		return result, err
	}
}

// DebounceLast is a function is a functionlast implementation that wraps a Circuit function
// and waits for a pause after a seris of calls before calling the inner function.
// Since this implementation won't provide an immediate response, it is useful only if
// your function doesn't need results ASAP
func DebounceLast(circuit Circuit, d time.Duration) Circuit {
	var threshold time.Time = time.Now()
	var ticker *time.Ticker
	var result string
	var err error
	var once sync.Once
	var m sync.Mutex

	return func(ctx context.Context) (string, error) {
		m.Lock()
		defer m.Unlock()

		threshold = time.Now().Add(d)

		once.Do(func() {
			ticker = time.NewTicker(time.Millisecond * 100)

			go func() {
				defer func() {
					m.Lock()
					ticker.Stop()
					// reset to allow the next cluster of calls to be considered
					once = sync.Once{}
					m.Unlock()
				}()

				for {
					select {
					case <-ticker.C:
						m.Lock()
						// verify enough time has passed since the last call
						if time.Now().After(threshold) {
							result, err = circuit(ctx)
							m.Unlock()
							return
						}
						m.Unlock()
					case <-ctx.Done():
						m.Lock()
						result, err = "", ctx.Err()
						m.Unlock()
						return
					}
				}
			}()
		})
		return result, err
	}
}
