package main

import (
	"fmt"
	"time"
)

func main() {
	doWork := func(done chan interface{}, pulseInterval time.Duration, nums ...int) (<-chan interface{}, <-chan int) {
		// 容量为1
		heartbeatStream := make(chan interface{}, 1)
		workStream := make(chan int)

		go func() {
			defer close(heartbeatStream)
			defer close(workStream)

			pulse := time.Tick(pulseInterval)

		numLoop:
			for _, n := range nums {
				if n == 2 {
					time.Sleep(time.Second * 3)
				}
				for {
					select {
					case <-done:
						return
					case workStream <- n:
						continue numLoop
					case <-pulse:
						select {
						case heartbeatStream <- struct{}{}:
						default:
						}
					}
				}
			}
		}()

		return heartbeatStream, workStream
	}

	done := make(chan interface{})
	defer close(done)

	intSlice := []int{0, 1, 2, 3, 4, 100}
	const timeout = 2 * time.Second
	heartbeat, results := doWork(done, timeout/2, intSlice...)

	// 保证goroutine的循坏已启动，否则下面的for循环的case可能超时
	<-heartbeat

	i := 0
	// 如果上述goroutine的某个循环迭代很耗时，此处加入超时判断
	for {
		select {
		case r, ok := <-results:
			if ok == false {
				return
			} else if expected := intSlice[i]; r != expected {
				fmt.Errorf("index %v: expected %v, but received %v\n", i, expected, r)
			} else {
				fmt.Printf("index %v: expected == received (%v)\n", i, r)
			}
			i++
		case <-heartbeat:
			fmt.Println("pulse")
		case <-time.After(timeout):
			fmt.Println("test time out")
		}
	}

	fmt.Println("test finished")
}
