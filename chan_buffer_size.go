package main

import (
	"fmt"
	"time"
)

func main() {
	// 不带缓冲的channel写完就阻塞，这种情况只有其他协程中有对应的读才能解除阻塞。
	// 而带缓冲的channel要直到写满+1才阻塞。
	//with_panic()
	without_panic()
}

func without_panic() {
	ch := make(chan int, 1)
	for {
		select {
		case ch <- 0:
		case ch <- 1:
		}
		i := <-ch
		fmt.Println("Value received:", i) // 随机输出0和1
		time.Sleep(time.Second)
	}
}

func with_panic() {
	ch := make(chan int)
	for {
		select {
		case ch <- 0:
		case ch <- 1:
		}
		i := <-ch
		fmt.Println("Value received:", i) // 报错：fatal error: all goroutines are asleep - deadlock!
	}
}
