package timeout

import "context"

// SlowFunction represents a function call which may, or may not, complete
// in a reasonable amount of time
type SlowFunction func(string) (string, error)

type WithContext func(context.Context, string) (string, error)

// Timeout wraps a SlowFunction to provide it a context, allowing
// to run it in a separate goroutine for a maximum set amount of time
func Timeout(f SlowFunction) WithContext {
	return func(ctx context.Context, arg string) (string, error) {
		chres := make(chan string)
		cherr := make(chan error)

		// run the slow function in its own goroutine
		go func() {
			res, err := f(arg)
			chres <- res
			cherr <- err
		}()

		select {
		case res := <-chres:
			return res, <-cherr
		// running for too long
		case <-ctx.Done():
			return "", ctx.Err()
		}
	}
}
