package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	// returns a read-only channel
	newRandStream := func(done <-chan interface{}) <-chan int {
		randStream := make(chan int)
		go func() {
			defer fmt.Println("newRandStream closure existed.")
			defer close(randStream)
			for {
				// select blocks until one of the cases is activated
				// all cases (channel reads and writes) are run simultaneously and the first one that is ready will execute
				select {
				case randStream <- rand.Int():
				// terminates the gorountine when an update is received by the done channel
				case <-done:
					return
				}
			}
		}()
		return randStream
	}

	// we create a channel and pass it to the sub goroutine, we can use it as a signal to finish the goroutine.
	done := make(chan interface{})
	randStream := newRandStream(done)
	fmt.Println("3 random ints:")
	for i := 1; i <= 3; i++ {
		fmt.Printf("%d: %d\n", i, <-randStream)
	}
	// closing the channel will signal the gorountine above, the case <- done is hit and it returns.
	close(done)
	// gives time for final closure messages to print
	time.Sleep(1 * time.Second)
}
