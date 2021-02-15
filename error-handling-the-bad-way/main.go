package main

import (
	"fmt"
	"net/http"
)

func main() {

	checkStatus := func(done <-chan interface{}, urls []string) <-chan *http.Response {
		responses := make(chan *http.Response)
		go func() {
			defer close(responses)
			for _, url := range urls {
				resp, err := http.Get(url)
				// This is not desirable. The goroutine should not have to know how to handle an error if one occurs
				if err != nil {
					// The error will be printed but the value will not be written to the channel which is ambiguous behaviour
					fmt.Println(err)
					continue
				}

				select {
				case <-done:
					return
				case responses <- resp:
				}
			}
		}()
		return responses
	}

	done := make(chan interface{})
	defer close(done)

	urls := []string{"https://www.google.com", "https://badhost", "https://www.facebook.com"}
	for response := range checkStatus(done, urls) {
		fmt.Printf("Response: %v\n", response.Status)
	}

}
