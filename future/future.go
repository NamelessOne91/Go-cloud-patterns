package future

import (
	"sync"
)

// Future provides a placeholder for values that are still generated
// by an asynchronous process
type Future interface {
	Result() (string, error)
}

// InnerFuture implements the Future interface and provides concurrent
// functionality
type InnerFuture struct {
	once  sync.Once
	wg    sync.WaitGroup
	res   string
	err   error
	resCh <-chan string
	errCh <-chan error
}

// Result attempts to read results and send them to the InnerFuture struct,
// subsequent calls will receive cached results
func (f *InnerFuture) Result() (string, error) {
	f.once.Do(func() {
		f.wg.Add(1)
		defer f.wg.Done()
		f.res = <-f.resCh
		f.err = <-f.errCh
	})
	// using a wait group makes the func thread safe
	f.wg.Wait()

	return f.res, f.err
}
