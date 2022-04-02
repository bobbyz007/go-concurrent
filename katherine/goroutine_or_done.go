package main

import (
	"fmt"
	"go-concurrent/katherine/util"
	"time"
)

func main() {
	done := make(chan interface{})
	defer close(done)

	s1 := util.Take(done, util.Repeat(done, 1), 5)
	for v := range util.OrDone(done, s1) {
		fmt.Printf("v: %v\n", v)
		time.Sleep(time.Millisecond * 500)
	}
}
