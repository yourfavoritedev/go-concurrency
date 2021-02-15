package main

import "fmt"

func getTotal(nums []int32, c chan int32) {
	var total int32
	for _, v := range nums {
		total += v
	}
	c <- total
}

func main() {
	totalChan := make(chan int32)
	group1 := []int32{5, 10, 15, 20}
	group2 := []int32{2, 4, 6, 8}
	var allGroups [][]int32
	allGroups = append(allGroups, group1, group2)

	for _, v := range allGroups {
		go getTotal(v, totalChan)
	}

	for range allGroups {
		v := <-totalChan
		fmt.Println(v)
	}
}
