package circuitbreaker

import (
	"context"
	"errors"
	"testing"
	"time"
)

var errFailedService error = errors.New("something went wrong")
var calls int

func malfunctioningCircuit(ctx context.Context) (string, error) {
	calls++
	if calls == 2 || calls == 3 {
		return "", errFailedService
	}
	return "OK", nil
}

func TestBreaker(t *testing.T) {
	tests := []struct {
		name     string
		wait     time.Duration
		expected string
		err      error
	}{
		{"working circuit", 0, "OK", nil},
		{"broken circuit - 2s backoff", 0, "", errFailedService},
		{"call after 1st fail", 1, "", errServiceUnreachable},
		{"another failed service call - 4s backoff", 3, "", errFailedService},
		{"immediate call after 2nd fail", 1, "", errServiceUnreachable},
		{"still in 2nd fail backoff", 3, "", errServiceUnreachable},
		{"service back up after 2nd fail", 5, "OK", nil},
	}

	ctx := context.Background()
	circuitBreaker := Breaker(malfunctioningCircuit, 0)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			time.Sleep(tt.wait * time.Second)
			res, err := circuitBreaker(ctx)
			if err != nil {
				if tt.err == nil {
					t.Errorf("Expected no error - got %v", err)
				}
				if tt.err != err {
					t.Errorf("Expected error %v - got %v", tt.err, err)
				}
			} else {
				if tt.err != nil {
					t.Errorf("Expected error %v - got %v", tt.err, err)
				} else if res != tt.expected {
					t.Errorf("Expected %s - got %s", tt.expected, res)
				}
			}

		})
	}
}
