package main

import (
	"fmt"
	"runtime"
	"time"
)

// temporary file, may delete anytime.
func main() {
	fmt.Println(runtime.NumCPU())

	getS := func(done chan interface{}) chan interface{} {
		s := make(chan interface{}, 2)
		go func() {
			defer close(s)

			for {
				select {
				case s <- time.Now():
					fmt.Println("added")
					time.Sleep(time.Second)
				case <-done:
					return
				}
			}
		}()
		return s
	}

	done := make(chan interface{})
	s := getS(done)
	time.AfterFunc(5*time.Second, func() {
		close(done)
	})
	for i := range s {
		fmt.Println("i :", i)
		time.Sleep(time.Minute)
	}
}

func f1() {
	fmt.Println("f1")
}

func f2() {
	fmt.Println("f2")
}
