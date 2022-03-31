package main

import (
	"sync"
	"testing"
)

// go test -bench=. -cpu=1 ./katherine/benchmark_test.go
func BenchmarkContextSwitch(b *testing.B) {
	var wg sync.WaitGroup
	begin := make(chan struct{})
	c := make(chan struct{})

	// 发送的数据
	var token struct{}

	sender := func() {
		defer wg.Done()
		<-begin

		for i := 0; i < b.N; i++ {
			// 发送数据
			c <- token
		}
	}

	receiver := func() {
		defer wg.Done()
		<-begin

		for i := 0; i < b.N; i++ {
			// 读取数据
			<-c
		}
	}

	wg.Add(2)

	go sender()
	go receiver()

	b.StartTimer()
	// 触发sender和receiver开始执行
	close(begin)
	wg.Wait()
}
