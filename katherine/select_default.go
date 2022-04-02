package main

import (
	"fmt"
	"time"
)

func main() {
	done := make(chan interface{})
	go func() {
		time.Sleep(time.Second * 3)
		close(done)
	}()

	workerCount := 0
loop:
	for {
		// 如果done堵塞了， 则几乎瞬间执行default 退出select语句
		select {
		case <-done:
			break loop
		default:
		}

		workerCount++
		time.Sleep(time.Second)
	}

	fmt.Printf("Achieved %v cycles of work before signalled to stop.\n", workerCount)
}
