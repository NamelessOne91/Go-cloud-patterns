package fanin

import "sync"

// Funnel multiplexes one, or more, input channels onto a single destination channel
func Funnel(sources ...<-chan int) <-chan int {
	// shared output channel
	dest := make(chan int)

	var wg sync.WaitGroup
	wg.Add(len(sources))

	// multiplex each source onto the shared channel
	for _, ch := range sources {
		go func(c <-chan int) {
			defer wg.Done()

			for n := range c {
				dest <- n
			}
		}(ch)
	}

	go func() {
		wg.Wait()
		close(dest)
	}()

	return dest
}
