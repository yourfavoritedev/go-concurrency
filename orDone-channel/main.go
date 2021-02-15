package main

import "fmt"

func main() {
	// This pattern is useful for working with channels who's behavior is independent of the goroutines you're working with.
	// You don't know if the fact that your goroutine was canceled means that the channel you're reading from will have been canceled.
	orDone := func(done, c <-chan interface{}) <-chan interface{} {
		valStream := make(chan interface{})
		go func() {
			defer close(valStream)
			for {
				select {
				case <-done:
					return
				// read from independent channel until we no longer get a value
				case v, ok := <-c:
					// exit the goroutine if ok is falsey (meaning there is nothing (nil) to read from the channel)
					if ok == false {
						return
					}
					select {
					case valStream <- v:
					case <-done:
					}
				}
			}
		}()
		return valStream
	}

	// signal to terminate orDone
	done := make(chan interface{})
	defer close(done)
	// populate channel
	myChan := make(chan interface{})
	integers := []int{1, 2, 3, 4}
	go func() {
		defer close(myChan)
		for _, v := range integers {
			myChan <- v
		}
	}()

	for val := range orDone(done, myChan) {
		fmt.Println(val)
	}
}
