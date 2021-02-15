package main

import (
	"testing"
)

// TestDoWorkGeneratesAllNumbers will test the DoWork heartbeats logic
func TestDoWorkGeneratesAllNumbers(t *testing.T) {
	done := make(chan interface{})
	defer close(done)

	intSlice := []int{0, 1, 2, 3, 4}
	heartbeat, results := DoWork(done, intSlice)

	<-heartbeat

	i := 0
	for r := range results {
		if expected := intSlice[i]; r != expected {
			t.Errorf("index %v: expected %v, but received %v,", i, expected, r)
		}
		i++
	}
}
