package main

import (
	"fmt"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	done := make(chan interface{})

	// spawn two different go-routines that are blocked by the done channel
	wg.Add(1)
	go func() {
		defer wg.Done()
		select {
		/**
		* the syntax <-done, has a use case other than "reading" from a channel
		* it's a blocker, waiting for a signal that a value can be read or the channel is closed
		* closing the channel will trigger the following case and allow it to execute
		 */
		case <-done:
			fmt.Println("signal from done channel")
			return
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		select {
		/**
		* this case is also anticipating a signal to the done channel
		* either a value can be read from the channel or the channel is closed
		 */
		case <-done:
			fmt.Println("signal from done channel woof")
			return
		}
	}()

	close(done)
	wg.Wait()
}
