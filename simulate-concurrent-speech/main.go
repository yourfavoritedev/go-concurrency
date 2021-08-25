package main

import (
	"fmt"
	"strings"
	"sync"
)

type person struct {
	name  string
	intro string
}

func speak(
	done chan interface{},
	phrase string,
) <-chan string {
	results := make(chan string)
	words := strings.Split(phrase, " ")
	go func() {
		defer close(results)
		for _, v := range words {
			select {
			case <-done:
				return
			case results <- v:
			}
		}
	}()
	return results
}

func main() {
	monica := person{
		name:  "Monica",
		intro: "Lorem Ipsum is simply dummy text of the printing and typesetting industry. Lorem Ipsum has been the industry's standard dummy text ever since the 1500s, when an unknown printer took a galley of type and scrambled it to make a type specimen book. It has survived not only five centuries, but also the leap into electronic typesetting, remaining essentially unchanged. It was popularised in the 1960s with the release of Letraset sheets containing Lorem Ipsum passages, and more recently with desktop publishing software like Aldus PageMaker including versions of Lorem Ipsum.",
	}
	chandler := person{
		name:  "Chandler",
		intro: "There are many variations of passages of Lorem Ipsum available, but the majority have suffered alteration in some form, by injected humour, or randomised words which don't look even slightly believable. If you are going to use a passage of Lorem Ipsum, you need to be sure there isn't anything embarrassing hidden in the middle of text. All the Lorem Ipsum generators on the Internet tend to repeat predefined chunks as necessary, making this the first true generator on the Internet. It uses a dictionary of over 200 Latin words, combined with a handful of model sentence structures, to generate Lorem Ipsum which looks reasonable. The generated Lorem Ipsum is therefore always free from repetition, injected humour, or non-characteristic words etc.",
	}

	doneMonica := make(chan interface{})
	doneChandler := make(chan interface{})

	var wg sync.WaitGroup

	defer close(doneMonica)
	defer close(doneChandler)

	// simulate concurrent speech
	wg.Add(1)
	go func() {
		defer fmt.Println("-----Done-----")
		for v := range speak(doneMonica, monica.intro) {
			fmt.Printf("%s: %s\n", monica.name, v)
		}
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		defer fmt.Println("-----Done-----")
		for v := range speak(doneChandler, chandler.intro) {
			fmt.Printf("%s: %s\n", chandler.name, v)
		}
		wg.Done()
	}()

	wg.Wait()
}
