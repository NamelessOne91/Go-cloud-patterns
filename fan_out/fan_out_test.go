package fanout

import (
	"sync"
	"testing"
)

func TestFanOut(t *testing.T) {
	tests := []struct {
		name         string
		destinations int
		elements     int
	}{
		{"1 destination, 10 elements", 1, 10},
		{"2 destinations, 10 elements", 2, 10},
		{"5 destinations, 10 elements", 5, 10},
		{"10 destination, 10 elements", 10, 10},
		{"10 destinations, 100 elements", 10, 100},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			source := make(chan int)
			dests := Split(source, tt.destinations)

			go func() {
				for i := 1; i <= tt.elements; i++ {
					source <- i
				}
				close(source)
			}()

			resMap := make(map[int]bool)
			var m sync.RWMutex

			var wg sync.WaitGroup
			wg.Add(len(dests))

			for _, ch := range dests {
				go func(dest <-chan int) {
					defer wg.Done()

					for val := range dest {
						m.Lock()
						resMap[val] = true
						m.Unlock()
					}
				}(ch)
			}
			wg.Wait()

			for i := 1; i <= tt.elements; i++ {
				_, ok := resMap[i]
				if !ok {
					t.Errorf("Expected key %d not found", i)
				}
			}
		})
	}
}
