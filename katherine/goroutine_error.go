package main

import (
	"fmt"
	"net/http"
)

type Result struct {
	Error    error
	Response *http.Response
}

// 错误应该被当做一等公民对待，随同结果一起返回
func main() {
	checkStatus := func(done <-chan interface{}, urls ...string) <-chan Result {
		results := make(chan Result)
		go func() {
			defer close(results)

			for _, url := range urls {
				resp, err := http.Get(url)
				result := Result{
					Error:    err,
					Response: resp,
				}
				select {
				case <-done:
					return
				case results <- result:
				}
			}
		}()

		return results
	}

	done := make(chan interface{})
	defer close(done)

	errCount := 0
	urls := []string{"https://www.baidu.com", "https://badhost", "a", "b", "c"}
	for result := range checkStatus(done, urls...) {
		if result.Error != nil {
			fmt.Printf("error: %v\n", result.Error)
			errCount++
			if errCount > 3 {
				fmt.Println("Too many errors, breaking!")
				break
			}
			continue
		}

		fmt.Printf("Response: %v\n", result.Response.Status)
	}
}
