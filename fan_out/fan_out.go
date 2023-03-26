package fanout

// Split accepts a single receive-only channel and splits its input
// between n channels
func Split(source <-chan int, n int) []<-chan int {
	dests := make([]<-chan int, n)

	for i := 0; i < n; i++ {
		ch := make(chan int)
		dests[i] = ch
		// each channel runs in a separate goroutine
		// competing for reads
		go func() {
			defer close(ch)

			for val := range source {
				ch <- val
			}
		}()
	}

	return dests
}
