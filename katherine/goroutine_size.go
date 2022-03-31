package main

import (
	"fmt"
	"runtime"
	"sync"
)

func main() {
	memConsumed := func() uint64 {
		runtime.GC()
		var s runtime.MemStats
		runtime.ReadMemStats(&s)
		// 占用内存字节
		return s.Sys
	}

	var c <-chan interface{}
	var wg sync.WaitGroup
	noop := func() {
		wg.Done()
		<-c // 堵塞
	}

	const numGoroutines = 1e4
	wg.Add(numGoroutines)

	before := memConsumed()
	for i := numGoroutines; i > 0; i-- {
		go noop()
	}
	after := memConsumed()

	// 因为所有goroutine没有结束，计算所占用的内存
	fmt.Printf("%.3f KB", float64(after-before)/numGoroutines/1024)
}
