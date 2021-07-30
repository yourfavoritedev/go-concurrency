package main

import (
	"fmt"
	"net/http"
)

// Result consumes an error and response to help with error-handling
type Result struct {
	Error    error
	Response *http.Response
}

func main() {
	checkStatus := func(done <-chan interface{}, urls []string) <-chan Result {
		results := make(chan Result)
		// this goroutine will be responsible for writing to the results channel
		go func() {
			// close results channel after the for-loop finishes or when the done channel closes
			defer close(results)
			for _, url := range urls {
				var result Result
				resp, err := http.Get(url)
				// We capture all returns from http.Get, write them to the channel, and let the main goroutine handle any errors
				result = Result{Error: err, Response: resp}
				// conventional select statement to end goroutine or write to channel
				select {
				case <-done:
					return
				case results <- result:
					fmt.Println("send result to channel")
				}
			}
		}()
		return results
	}

	// channel to signal the above goroutine to terminate
	done := make(chan interface{})
	defer close(done)

	urls := []string{"https://www.google.com", "https://badhost", "https://www.facebook.com"}
	// read from channel, it will stop reading when the channel is closed
	for result := range checkStatus(done, urls) {
		/** By handling errors in the main goroutine we establish a separation of concerns.
		The sub goroutine can focus on writing to the channel.
		The main goroutine can make decisions about what to do with the errors. */
		if result.Error != nil {
			/** this will panic if we do not continue (go to the next iteration)
			 because the current result wiil have a Response of nil
			 trying to print Response.Status (nil has no fields/values) will result in a panic
			**/
			fmt.Printf("error: %v\n", result.Error)
			continue
		}
		fmt.Printf("Response: %v\n", result.Response.Status)
	}
}
