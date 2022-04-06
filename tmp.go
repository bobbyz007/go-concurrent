package main

import (
	"fmt"
	"runtime"
)

// temporary file, may delete anytime.
func main() {
	fmt.Println(runtime.NumCPU())

}

func f1() {
	fmt.Println("f1")
}

func f2() {
	fmt.Println("f2")
}
