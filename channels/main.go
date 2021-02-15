package main

import "fmt"

func main() {
	data := make([]int, 3)

	loopData := func(c chan<- int) {
		defer close(c)
		for _, v := range data {
			c <- v
		}
	}

	c := make(chan int)
	go loopData(c)

	for num := range c {
		fmt.Println(num)
	}
}
