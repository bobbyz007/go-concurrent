package main

import (
	"fmt"
	"go-concurrent/katherine/util"
	"log"
	"os"
	"time"
)

func main() {
	doWorkFn := func(done <-chan interface{}, intList ...int) (util.StartGoroutineFn, <-chan interface{}) {
		intChanStream := make(chan (<-chan interface{}))
		intStream := util.Bridge(done, intChanStream)

		doWork := func(done <-chan interface{}, pulseInterval time.Duration) <-chan interface{} {
			intStream := make(chan interface{})
			heartbeat := make(chan interface{})
			go func() {
				defer close(intStream)

				select {
				case intChanStream <- intStream:
				case <-done:
					return
				}

				pulse := time.Tick(pulseInterval)

				for {
				valueLoop:
					// 可能返回重复数
					//for _, intVal := range intList {
					// 不会返回重复数
					for {
						if len(intList) == 0 {
							return
						}
						intVal := intList[0]
						intList = intList[1:]
						// 模拟故障
						if intVal < 0 {
							log.Printf("negative value: %v\n", intVal)
							return
						}
						for {
							select {
							case <-pulse:
								select {
								case heartbeat <- struct{}{}:
								default:
								}
							case intStream <- intVal:
								continue valueLoop
							case <-done:
								return
							}
						}
					}
				}

			}()

			return heartbeat
		}

		return doWork, intStream
	}

	log.SetFlags(log.Ltime | log.LUTC)
	log.SetOutput(os.Stdout)

	done := make(chan interface{})
	defer close(done)
	time.AfterFunc(time.Millisecond*8, func() {
		close(done)
	})

	doWork, intStream := doWorkFn(done, 1, 2, -1, 3, 4, 5, 6)
	doWorkWithSteward := util.NewSteward(time.Millisecond, doWork)

	// 启动steward：此处模拟情况，碰到负数时 goroutine返回，ward中断心跳，ward内部的intStream延时关闭； steward监听到无心挑，则重启ward
	doWorkWithSteward(done, time.Hour)

	for intVal := range util.Take(done, intStream, 6) {
		fmt.Printf("Received: %v\n", intVal)
	}
}
