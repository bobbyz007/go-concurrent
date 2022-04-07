package main

import (
	"fmt"
	"time"
)

func main() {
	doWork := func(done chan interface{}, nums ...int) (<-chan interface{}, <-chan int) {
		// 容量为1
		heartbeatStream := make(chan interface{}, 1)
		workStream := make(chan int)
		go func() {
			defer close(heartbeatStream)
			defer close(workStream)

			for _, n := range nums {
				// 模拟耗时操作
				if n == 2 {
					time.Sleep(time.Second * 3)
				}

				// 首先非堵塞发送心跳消息
				select {
				case heartbeatStream <- struct{}{}:
				default:
				}

				select {
				case <-done:
					return
				case workStream <- n:
				}
			}
		}()

		return heartbeatStream, workStream
	}

	done := make(chan interface{})
	defer close(done)

	intSlice := []int{0, 1, 2, 3, 4, 100}
	heartbeat, results := doWork(done, intSlice...)

	// 保证goroutine的循坏已启动，否则下面的for循环的case可能超时
	<-heartbeat

	i := 0
	// 存在的问题：如果上述goroutine的某个循环迭代很耗时，可能会导致此处等待很久。
	for r := range results {
		if expected := intSlice[i]; r != expected {
			fmt.Errorf("index %v: expected %v, but received %v\n", i, expected, r)
		} else {
			fmt.Printf("index %v: expected == received (%v)\n", i, r)
		}
		i++
	}

	fmt.Println("test finished")
}
