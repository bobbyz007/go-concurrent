package main

import (
	"fmt"
	"time"
)

var timeout2 <-chan time.Time
var result chan int

func main() {
	timeout2 = time.After(time.Second * 3)
	result = make(chan int)

	go func() {
		fmt.Println("---begin do task---")
		// time.Sleep(time.Millisecond * 3100)

		result <- 1
		fmt.Println("--end do task---")
	}()

	time.Sleep(time.Second * 2)
	select {
	case e := <-result:
		fmt.Printf("get result: %d\n", e)
	case <-timeout2:
		fmt.Println("get result timeout")
	}

	// wait until goroutine finish
	time.Sleep(time.Second)
}
