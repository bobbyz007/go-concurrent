package main

import (
	"fmt"
	"go-concurrent/katherine/util"
	"math/rand"
)

// 一些有用的generator
func main() {
	done := make(chan interface{})
	defer close(done)

	for num := range util.Take(done, util.Repeat(done, 1, 2, 3), 20) {
		fmt.Printf("%v ", num)
	}

	rand := func() interface{} {
		return rand.Int()
	}
	fmt.Println()
	for num := range util.Take(done, util.RepeatFn(done, rand), 10) {
		fmt.Printf("%v\n", num)
	}
}
