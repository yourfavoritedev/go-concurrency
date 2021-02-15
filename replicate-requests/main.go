package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func main() {
	doWork := func(
		done <-chan interface{},
		id int,
		wg *sync.WaitGroup,
		result chan<- int,
	) {
		started := time.Now()
		defer wg.Done()

		// Simulate random load
		simulatedLoadTime := time.Duration(1*rand.Intn(10)) * time.Second

		/** use two separate select blocks because we want to send/receive two different values, the time.After (receive) and the id (send).
		/ if they were in the same select block, then we could only use one value at a time, the other will get lost. */
		select {
		// do not want to return on <-done because we still want to log the time it took
		case <-done:
		// alternative path if <-done takes too long to finish, time.After() can serve as a time-limit
		case <-time.After(simulatedLoadTime):
		}

		select {
		case <-done:
		case result <- id:
		}

		took := time.Since(started)
		// Display how long handlers would have taken
		if took < simulatedLoadTime {
			took = simulatedLoadTime
		}
		fmt.Printf("%v took %v\n", id, took)
	}

	done := make(chan interface{})
	result := make(chan int)

	var wg sync.WaitGroup
	wg.Add(10)

	for i := 0; i < 10; i++ {
		go doWork(done, i, &wg, result)
	}

	firstReturned := <-result
	close(done)
	wg.Wait()

	fmt.Printf("Received an answer from #%v\n", firstReturned)
}
