package main

import (
	"fmt"
	"sync"
)

func main() {
	var wg sync.WaitGroup

	for _, salutation := range []string{"hello", "greetings", "good day"} {
		wg.Add(1)

		// 有问题，都指向salutation变量
		/*go func() {
			defer wg.Done()
			fmt.Println(salutation)
		}()*/

		go func(salutation string) {
			defer wg.Done()
			fmt.Println(salutation)
		}(salutation)
	}

	wg.Wait()
}
