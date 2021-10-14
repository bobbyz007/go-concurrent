package main

import (
	"fmt"
	"time"
)

func main() {
	quit := make(chan struct{})

	go func() {
		for {
			select {
			case <- quit:
				fmt.Println("Sub goroutine is over")
				return
			default:
				time.Sleep(time.Second)
				fmt.Println("Sub goroutine do something")
			}
		}
	}()

	time.Sleep(3 * time.Second)
	fmt.Println("main goroutine start stop sub goroutine")

	close(quit)

	time.Sleep(time.Second * 10)
	fmt.Println("main goroutine is over")
}


