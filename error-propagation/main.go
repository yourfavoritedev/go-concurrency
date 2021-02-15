package main

import (
	"fmt"
	"log"
	"os"

	intermediate "github.com/yourfavoritedev/learngo/first/concurrency/error-propagation/intermediate"
)

func handleError(key int, err error, message string) {
	// bind log to erro message with an Id
	log.SetPrefix(fmt.Sprintf("[logID: %v]: ", key))
	// log the entire error if someones needs to dig into what happened
	log.Printf("%#v", err)
	fmt.Printf("[%v] %v", key, message)
}

func main() {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.Ltime | log.LUTC)

	err := intermediate.RunJob("1")
	if err != nil {
		msg := "There was an unexpcted issue: please report this as a bug."
		// type assertion to check if the error is a well-crafted error to our standards
		if _, ok := err.(intermediate.IntermediateErr); ok {
			msg = err.Error()
		}
		handleError(1, err, msg)
	}
}
