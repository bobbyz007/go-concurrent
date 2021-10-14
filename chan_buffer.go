package main

import "fmt"

func main() {
	ch := make(chan int, 10)
	wg := make(chan int)

	go printer(ch, wg)

	for i := 1; i < 100; i++ {
		ch <- i
	}

	// stop reading from ch
	close(ch)
	fmt.Println("wait sub goroutine over")

	// wait to read from wg
	<- wg
	fmt.Println("main goroutine over")
}

func printer(ch <-chan int, wg chan<- int)  {
	// 接收通道
	for i := range ch {
		fmt.Println(i)
	}

	// 写入元素
	wg <- 1
	close(wg)
}
