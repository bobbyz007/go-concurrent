package main

import (
	"fmt"
	"go-concurrent/katherine/util"
	"math/rand"
	"runtime"
	"sync"
	"time"
)

func main() {
	toInt := func(done <-chan interface{}, valueStream <-chan interface{}) <-chan int {
		intStream := make(chan int)
		go func() {
			defer close(intStream)
			for v := range valueStream {
				select {
				case <-done:
					return
				case intStream <- v.(int):
				}
			}
		}()
		return intStream
	}

	primeFinder := func(done <-chan interface{}, intStream <-chan int) <-chan interface{} {
		primeStream := make(chan interface{})
		go func() {
			defer close(primeStream)
			for v := range intStream {
				select {
				case <-done:
					return
				default:
					if isPrime(v) {
						primeStream <- v
					}
				}
			}
		}()

		return primeStream
	}

	fanIn := func(done <-chan interface{}, channels ...<-chan interface{}) <-chan interface{} {
		var wg sync.WaitGroup
		multiplexedStream := make(chan interface{})
		multiplex := func(c <-chan interface{}) {
			defer wg.Done()
			for i := range c {
				select {
				case <-done:
					return
				case multiplexedStream <- i:
				}
			}
		}

		wg.Add(len(channels))
		for _, c := range channels {
			go multiplex(c)
		}

		go func() {
			wg.Wait()
			close(multiplexedStream)
		}()

		return multiplexedStream
	}

	rand := func() interface{} {
		return rand.Intn(50_000_000)
	}

	done := make(chan interface{})
	defer close(done)

	start := time.Now()
	randIntStream := toInt(done, util.RepeatFn(done, rand))
	fmt.Println("Primes: ")

	// no fan-out fan-in
	// ---------------------------------------------------------------
	/*for prime := range take(done, primeFinder(done, randIntStream), 10) {
		fmt.Printf("\t%d\n", prime)
	}*/
	// ---------------------------------------------------------------

	// with fan-out fan-in
	// ---------------------------------------------------------------
	numFinders := runtime.NumCPU()
	fmt.Printf("Spining up %d prime finders.\n", numFinders)
	finders := make([]<-chan interface{}, numFinders)
	for i := 0; i < numFinders; i++ {
		finders[i] = primeFinder(done, randIntStream)
	}
	for prime := range util.Take(done, fanIn(done, finders...), 10) {
		fmt.Printf("\t%d\n", prime)
	}
	// ---------------------------------------------------------------

	fmt.Printf("Search took: %v", time.Since(start))
}

func isPrime(i int) bool {
	if i <= 0 {
		return false
	}

	// 模拟耗时stage
	var end = i - 1 //int(math.Sqrt(float64(i)))
	for ; end >= 2; end-- {
		if i%end == 0 {
			return false
		}
	}

	return true
}
