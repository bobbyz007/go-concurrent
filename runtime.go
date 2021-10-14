package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

var wg2 sync.WaitGroup

func main() {
	// runtime.GOMAXPROCS(1)
	runtime.GOMAXPROCS(2)

	wg2.Add(2)

	go func() {
		for i := 0; i < 10; i++ {
			fmt.Println(i)
			time.Sleep(time.Second)
		}
		wg2.Done()
	}()
	
	go func() {
		for i := 11; i < 20; i++ {
			fmt.Println(i)
			time.Sleep(time.Second)
		}
		wg2.Done()
	}()

	wg2.Wait()
}

