package main
import (
	"fmt"
	"sync"
)
var (
	counter int//计数器
	wg4      sync.WaitGroup//信号量
	mutex2   sync.Mutex//互斥锁
)
func main() {
	//1.两个信号量
	wg4.Add(2)

	//2.开启两个线程
	go incCounter()
	go incCounter()

	//3.等待子线程结束
	wg4.Wait()
	fmt.Println(counter)
}
func incCounter() {
	defer wg4.Done()
	//2.1.获取锁
	mutex2.Lock()
	//2.2.计数加1
	counter++
	//2.3.释放独占锁
	mutex2.Unlock()
}