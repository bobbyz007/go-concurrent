package main

import "fmt"
import "sync"

var a string
var wg sync.WaitGroup

func f()  {
	fmt.Print(a)
	wg.Done()
}

func hello()  {
	a = "hello world"
	go f()
}

func main()  {
	wg.Add(1)
	hello()
	wg.Wait()
}
