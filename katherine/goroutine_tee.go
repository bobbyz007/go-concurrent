package main

import (
	"fmt"
	"go-concurrent/katherine/util"
)

func main() {
	done := make(chan interface{})
	defer close(done)

	out1, out2 := util.Tee(done, util.Take(done, util.Repeat(done, 1, 2), 4))
	for val1 := range out1 {
		fmt.Printf("out1: %v, out2: %v\n", val1, <-out2)
	}
}
