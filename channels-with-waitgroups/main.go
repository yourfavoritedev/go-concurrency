package main

import (
	"fmt"
	"strings"
	"sync"
)

func speak(phrase string, wg *sync.WaitGroup, c chan string) {
	go func() {
		defer wg.Done()
		words := strings.Split(phrase, " ")
		for _, w := range words {
			/**
			writing to an unbuffered channel (which is full by default) without any action to read from the channel
			will block the goroutine from completing, creating a deadlock
			*/
			c <- w
		}
	}()
}

func main() {
	var wg sync.WaitGroup
	speakChannel := make(chan string)
	philipWords := "Hello I'm Phillip"
	benjiWords := "Hi I am Benji"
	var allPhrases []string
	allPhrases = append(allPhrases, philipWords, benjiWords)

	for _, v := range allPhrases {
		wg.Add(1)
		speak(v, &wg, speakChannel)
	}

	// wg.Wait() <-- do not do this without a go-routine, this will block and prevent the range loop from executing, causing a deadlock
	// wg.Wait() is synchronized to pass once the WaitGroup counter hits 0, so all go-routines should be completed (speak will decrement the WaitGroup),
	// but the go-routines are blocked until the range (read) over speakChannel executes...
	// unbuffed channels require a read-action (<-c) to accept the incoming value from a send-action (c<-), we cannot write unless trying to read
	// Wrapping wg.Wait() in a go-routine allows the below range loop to execute and read from the channel
	go func() {
		wg.Wait()
		close(speakChannel)
	}()

	fmt.Println("unblocked")

	// this will unblock the queue because we are reading data from the channel
	for v := range speakChannel {
		fmt.Println(v)
	}
}
