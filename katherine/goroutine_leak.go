package main

import (
	"fmt"
	"time"
)

func main() {
	doWork := func(strings <-chan string) <-chan interface{} {
		completed := make(chan interface{})

		go func() {
			defer fmt.Println("doWork exited")
			defer close(completed)

			for s := range strings {
				fmt.Println(s)
			}
		}()

		return completed
	}

	// doWork 中的goroutine永远不会结束，造成goroutine泄露
	doWork(nil)

	time.Sleep(time.Hour * 3)
	fmt.Println("Done.")
}
