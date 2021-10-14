package main

import "fmt"
import "sort"

var c chan int
var nums []int

func main() {
	// 无缓冲通道
	c = make(chan int)
	nums = []int{5, 3, 7, 2, 9, 1, 6}

	go func() {
		sort.Ints(nums)
		c <- 1
	}()

	// 堵塞，直到通道内有元素
	<- c
	fmt.Println(nums)
}
