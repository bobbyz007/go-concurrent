package main

import (
	"context"
	"time"
)
import "fmt"

func main() {
	gen := func(ctx context.Context) <-chan int {
		dst := make(chan int)
		n := 1
		go func() {
			for {
				select {
				case <- ctx.Done():
					fmt.Println(ctx.Err())
				case dst <- n:
					n++
				}
			}
		}()
		return dst
	}

	ctx, cancel := context.WithCancel(context.Background())
	for n := range gen(ctx) {
		fmt.Println(n)
		if n == 5 {
			cancel()
			break
		}
	}

	time.Sleep(time.Nanosecond)
}
