package throttle

import (
	"context"
	"testing"
	"time"
)

func myEffector(ctx context.Context) (string, error) {
	return "OK", nil
}

func TestThrottle(t *testing.T) {
	tests := []struct {
		name     string
		max      uint
		refill   uint
		interval time.Duration
		calls    uint
		expected string
		err      error
	}{
		{"No throttling applied", 3, 1, 1 * time.Second, 2, "OK", nil},
		{"No bucket equals instant throttling", 0, 1, 1 * time.Second, 1, "", errThrottling},
		{"Too many calls", 2, 1, 1 * time.Second, 6, "", errThrottling},
		{"Refill prevents throttling", 3, 1, 1 * time.Second, 5, "OK", nil},
		{"Context deadline", 10, 1, 1 * time.Second, 15, "", context.DeadlineExceeded},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := Throttle(myEffector, tt.max, tt.refill, 2*time.Second)
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

			var res string
			var err error
			for i := tt.calls; i > 0; i-- {
				res, err = e(ctx)
				time.Sleep(tt.interval)
			}

			cancel()
			if err != tt.err {
				t.Errorf("Expected error '%s' - got '%s'", tt.err, err)
			} else if res != tt.expected {
				t.Errorf("Expected result '%s' - got '%s'", tt.expected, res)
			}
		})
	}
}
