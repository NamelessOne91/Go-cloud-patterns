package fanin

import (
	"testing"
)

func TestFunnel(t *testing.T) {
	tests := []struct {
		name      string
		numChans  int
		numInputs int
	}{
		{"1 channel, 10 elements", 1, 10},
		{"1 channel, 100 elements", 1, 100},
		{"3 channels, 10 elements", 3, 10},
		{"10 channels, 10 elements", 10, 10},
		{"10 channels, 100 elements", 10, 100},
		{"100 channels, 100 elements", 100, 100},
		{"100 channels, 1000 elements", 100, 1000},
		{"10 channels, 1000000 elements", 10, 1000000},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			expected := tt.numChans * tt.numInputs
			sources := make([]<-chan int, tt.numChans)

			for i := tt.numChans - 1; i >= 0; i-- {
				source := make(chan int)
				sources[i] = source

				go func() {
					defer close(source)
					for num := tt.numInputs; num > 0; num-- {
						source <- num
					}
				}()
			}

			sink := Funnel(sources...)
			var count int
			resMap := make(map[int]int)

			for n := range sink {
				count++
				resMap[n]++
			}

			if count != expected {
				t.Errorf("Expected %d channels * %d elements = %d elements - got %d", tt.numChans, tt.numInputs, expected, count)
			}
			for i := 1; i <= tt.numInputs; i++ {
				amount := resMap[i]
				if amount != tt.numChans {
					t.Errorf("Expected %d to have been inserted %d times - got %d", i, tt.numChans, amount)
				}
			}
		})
	}
}
