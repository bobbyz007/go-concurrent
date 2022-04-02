package main

import (
	"fmt"
	"time"
)

func main() {
	// channels中有任一channel关闭或写入，则or函数返回的合并channel 关闭
	var or func(channels ...<-chan interface{}) <-chan interface{}
	or = func(channels ...<-chan interface{}) <-chan interface{} {
		switch len(channels) {
		case 0:
			return nil
		case 1:
			return channels[0]
		}

		orDone := make(chan interface{})
		go func() {
			defer close(orDone)
			switch len(channels) {
			case 2:
				select {
				case <-channels[0]:
				case <-channels[1]:
				}
			default:
				select {
				case <-channels[0]:
				case <-channels[1]:
				case <-channels[2]:
				// 递归调用or，如果or返回的channel已关闭，则case <- 就不会堵塞。
				case <-or(append(channels[3:], orDone)...):
				}
			}
		}()

		// 返回控制channel
		return orDone
	}

	sig := func(after time.Duration) <-chan interface{} {
		c := make(chan interface{})
		go func() {
			defer close(c)
			time.Sleep(after)
		}()
		return c
	}

	start := time.Now()

	<-or(sig(time.Hour),
		sig(time.Minute),
		sig(time.Second*3),
		sig(time.Hour*2),
		sig(time.Minute*3),
	)

	fmt.Printf("done after %v", time.Since(start))

}
