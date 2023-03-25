package retry

import (
	"context"
	"errors"
	"testing"
	"time"
)

var errTransientFailure = errors.New("fail")
var count int

// conforms to the Effector signature
func EmulateTransientError(ctx context.Context) (string, error) {
	count++

	if count <= 3 {
		return "", errTransientFailure
	}
	return "OK", nil
}

func TestRetry(t *testing.T) {
	tests := []struct {
		name     string
		retries  int
		wait     time.Duration
		expected string
		err      error
	}{
		{"Enough retries and time to succeed", 4, 5, "OK", nil},
		{"Not enough retries to succeed", 2, 5, "", errTransientFailure},
		{"Not enough time to succeed", 4, 2, "", context.Canceled},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			count = 0
			eff := Retry(EmulateTransientError, tt.retries, 1*time.Second)
			ctx, cancel := context.WithCancel(context.Background())
			go func() { time.Sleep(tt.wait * time.Second); cancel() }()

			res, err := eff(ctx)
			if err != tt.err {
				t.Errorf("Expected error '%s' - got '%s'", tt.err, err)
			} else if res != tt.expected {
				t.Errorf("Expected '%s' - got '%s'", tt.expected, res)
			}
		})
	}
}
