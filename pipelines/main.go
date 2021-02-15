package main

import "fmt"

func main() {
	// conventionally known as a generator, this func takes a batch of data and converts it into a channel. Thus providing the initial input for our pipeline.
	generator := func(done <-chan interface{}, integers []int) <-chan int {
		intStream := make(chan int)
		go func() {
			defer close(intStream)
			for _, v := range integers {
				select {
				case <-done:
					return
				case intStream <- v:
				}
			}
		}()
		return intStream
	}

	// a stage in our pipeline, it accepts a stream of data and sends back a stream of data of the same type
	multiply := func(done <-chan interface{}, intStream <-chan int, multiplier int) <-chan int {
		multipliedStream := make(chan int)
		go func() {
			defer close(multipliedStream)
			for v := range intStream {
				select {
				case <-done:
					return
				case multipliedStream <- v * multiplier:
				}
			}
		}()
		return multipliedStream
	}

	// additional stage, same input and output type
	add := func(done <-chan interface{}, intStream <-chan int, additive int) <-chan int {
		additiveStream := make(chan int)
		go func() {
			defer close(additiveStream)
			for v := range intStream {
				select {
				case <-done:
					return
				case additiveStream <- v + additive:
				}
			}
		}()
		return additiveStream

	}

	done := make(chan interface{})
	defer close(done)

	// start pipeline
	intStream := generator(done, []int{1, 2, 3, 4})
	pipeline := add(done, multiply(done, intStream, 2), 1)

	for v := range pipeline {
		fmt.Println(v)
	}
}
