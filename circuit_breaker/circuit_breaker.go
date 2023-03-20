package circuitbreaker

import (
	"context"
	"errors"
	"sync"
	"time"
)

var errServiceUnreachable error = errors.New("service unreachable")

// Circuit represents a function interacting with an upstream service;
// it should include an error in its return list
type Circuit func(context.Context) (string, error)

// Breaker wraps a Circuit function to provide a reset mechanism, allowing to retry
// services call applying an exponential backoff
func Breaker(circuit Circuit, threshold uint) Circuit {
	var consecutiveFailures int = 0
	var lastAttemp = time.Now()
	var m sync.RWMutex

	return func(ctx context.Context) (string, error) {
		// read lock
		m.RLock()
		d := consecutiveFailures - int(threshold)

		if d > 0 {
			// exponential backoff
			retryAt := lastAttemp.Add(time.Second * 2 << d)
			if !time.Now().After(retryAt) {
				m.RUnlock()
				return "", errServiceUnreachable
			}
		}
		//release read lock
		m.RUnlock()

		// issue request
		response, err := circuit(ctx)

		// lock on lastAttempt
		m.Lock()
		defer m.Unlock()

		lastAttemp = time.Now()

		// service call failed
		if err != nil {
			consecutiveFailures++
			return response, err
		}

		// service call successfull, reset counter
		consecutiveFailures = 0
		return response, nil
	}
}
