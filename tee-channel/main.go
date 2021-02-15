package main

import "fmt"

func main() {
	orDone := func(done, c <-chan interface{}) <-chan interface{} {
		valStream := make(chan interface{})
		go func() {
			defer close(valStream)
			for {
				select {
				case <-done:
					return
				case v, ok := <-c:
					if ok == false {
						return
					}
					select {
					case valStream <- v:
						fmt.Println("write from valStream")
					case <-done:
					}
				}
			}
		}()
		return valStream
	}

	tee := func(done <-chan interface{}, in <-chan interface{}) (_, _ <-chan interface{}) {
		out1 := make(chan interface{})
		out2 := make(chan interface{})
		go func() {
			defer close(out1)
			defer close(out2)
			for val := range orDone(done, in) {
				// shadow out1, out2
				var out1, out2 = out1, out2
				// Use one select statement so that writes to out1 and out2 don't block each other
				// To ensure both are written to, we'll perform two iterations of the select statement: one for each outbound channel
				for i := 0; i < 2; i++ {
					fmt.Println("write to teechannel")
					select {
					case <-done:
					case out1 <- val:
						out1 = nil // Once we've written to a channel, we set its shadowed copy to nil (full) so that further writes will block
					case out2 <- val:
						out2 = nil // Same
					}
				}
			}
		}()
		return out1, out2
	}

	done := make(chan interface{})
	defer close(done)

	myChan := make(chan interface{})
	integers := []int{1, 2, 3, 4}
	go func() {
		defer close(myChan)
		for _, v := range integers {
			fmt.Println("Write from myChan")
			myChan <- v
		}
	}()

	out1, out2 := tee(done, myChan)
	for val1 := range out1 {
		fmt.Printf("out1: %v, out2: %v\n", val1, <-out2)
	}
}
