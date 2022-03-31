package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	const QueueSize = 10
	var wg sync.WaitGroup
	wg.Add(QueueSize)

	c := sync.NewCond(&sync.Mutex{})
	queue := make([]interface{}, 0, QueueSize)

	removeFromQueue := func(delay time.Duration) {
		defer wg.Done()

		time.Sleep(delay)
		c.L.Lock()

		queue = queue[1:]
		fmt.Println("Removed from queue")
		c.L.Unlock()
		c.Signal()
	}

	for i := 0; i < QueueSize; i++ {
		c.L.Lock()
		// 队列长度满了，不能添加了
		for len(queue) == 2 {
			c.Wait()
		}
		fmt.Println("Adding to queue")
		queue = append(queue, struct{}{})

		go removeFromQueue(time.Second)
		c.L.Unlock()
	}

	// wait for consume all elements
	wg.Wait()
}
