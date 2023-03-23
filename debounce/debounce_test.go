package debounce

import (
	"context"
	"testing"
	"time"
)

var callCounter int

func testCircuit(ctx context.Context) (string, error) {
	callCounter++
	return "OK", nil
}

func TestDebounceFirst(t *testing.T) {
	callCounter = 0
	expectedCalls := 1
	wait := time.Second * 1

	ctx := context.Background()
	circuit := DebounceFirst(testCircuit, wait)

	for i := 0; i < 100; i++ {
		_, err := circuit(ctx)
		if err != nil {
			t.Errorf("Expected no error - got %s", err)
		}
	}

	if callCounter != expectedCalls {
		t.Errorf("Expected %d to circuit() - got %d", expectedCalls, callCounter)
	}
}

func TestDebounceLast(t *testing.T) {
	callCounter = 0
	expectedCalls := 1
	wait := time.Second * 1

	ctx := context.Background()
	circuit := DebounceLast(testCircuit, wait)

	for i := 0; i < 100; i++ {
		_, err := circuit(ctx)
		if err != nil {
			t.Errorf("Expected no error - got %s", err)
		}
	}

	time.Sleep(wait * 2)
	if callCounter != expectedCalls {
		t.Errorf("Expected %d to circuit() - got %d", expectedCalls, callCounter)
	}
}
