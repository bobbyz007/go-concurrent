package main

import (
	"fmt"
	"go-concurrent/katherine/util"
)

func main() {
	genVals := func() <-chan <-chan interface{} {
		chanStream := make(chan (<-chan interface{}))
		go func() {
			defer close(chanStream)
			for i := 0; i < 10; i++ {
				stream := make(chan interface{}, 1)
				stream <- i
				close(stream)
				chanStream <- stream
			}
		}()
		return chanStream
	}

	for v := range util.Bridge(nil, genVals()) {
		fmt.Printf("%v ", v)
	}

}
