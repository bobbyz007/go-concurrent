package main

import (
	"fmt"
	"sync"
)

func main() {
	hello := func(wg *sync.WaitGroup, id int) {
		defer wg.Done()
		fmt.Printf("Hello from %v!\n", id)
	}

	const numberGreeters = 5
	var wg sync.WaitGroup
	wg.Add(numberGreeters)

	for i := 0; i < numberGreeters; i++ {
		go hello(&wg, i+1)
	}

	wg.Wait()
}
