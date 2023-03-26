package future

import (
	"context"
	"testing"
	"time"
)

// mySlowFunction is a wrapper for core functionality, it allows
// asynchronous retrieval of return values through a Future implementation
func mySlowFunction(ctx context.Context) Future {
	resCh := make(chan string)
	errCh := make(chan error)

	go func() {
		select {
		case <-time.After(time.Second * 2):
			resCh <- "OK"
			errCh <- nil
		case <-ctx.Done():
			resCh <- ""
			errCh <- ctx.Err()
		}
	}()

	return &InnerFuture{resCh: resCh, errCh: errCh}
}

func TestFuture(t *testing.T) {
	tests := []struct {
		name     string
		timeout  time.Duration
		expected string
		err      error
	}{
		{"Future returs proper result", 5, "OK", nil},
		{"Context timeout", 1, "", context.DeadlineExceeded},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), tt.timeout*time.Second)
			defer cancel()

			future := mySlowFunction(ctx)

			for i := 0; i < 3; i++ {
				res, err := future.Result()
				if err != tt.err {
					t.Errorf("expected error '%s' - got '%v'", tt.err, err)
				} else if res != tt.expected {
					t.Errorf("expected result '%s' -  got '%s'", tt.expected, res)
				}
			}
		})
	}
}
