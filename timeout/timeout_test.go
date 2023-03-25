package timeout

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func mySlowFunc(in string) (string, error) {
	time.Sleep(5 * time.Second)
	return in, nil
}

func TestTimeout(t *testing.T) {
	tests := []struct {
		name     string
		timeout  time.Duration
		expected string
		err      error
	}{
		{"Function responds in time", 7, "OK", nil},
		{"Function takes too long", 3, "", context.DeadlineExceeded},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), tt.timeout*time.Second)
			defer cancel()

			f := Timeout(mySlowFunc)
			res, err := f(ctx, "OK")

			if err != tt.err {
				fmt.Printf("Expected  error '%s' - got '%v'", tt.err, err)
			} else if res != tt.expected {
				fmt.Printf("Expected res to be '%s' - got '%s'", tt.expected, res)
			}
		})
	}
}
