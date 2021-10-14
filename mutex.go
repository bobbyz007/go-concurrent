package main

import "sync"

var (
	mutex sync.Mutex
	wg3 sync.WaitGroup
)

func main() {
	wg3.Add(1)
	go func() {
		// 新的goroutine获取锁
		mutex.Lock()
		wg3.Done()
	}()

	wg3.Wait()

	// 主goroutine释放锁
	mutex.Unlock()
}