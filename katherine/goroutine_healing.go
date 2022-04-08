package main

import (
	"go-concurrent/katherine/util"
	"log"
	"os"
	"time"
)

// 基于heartbeat模式，代表被监控的goroutine
type startGoroutineFn func(done <-chan interface{}, pulseInterval time.Duration) (heartbeat <-chan interface{})

func main() {
	// 返回类型是 startGoroutineFn，表示steward本身也是可以被监控的
	// timeout: ward的超时时间。 如果在timeout时间内 获取不到ward的心跳信息，则会重启ward
	newSteward := func(timeout time.Duration, startGoroutine startGoroutineFn) startGoroutineFn {
		return func(done <-chan interface{}, pulseInterval time.Duration) <-chan interface{} {
			heartbeat := make(chan interface{})
			go func() {
				defer close(heartbeat)

				var wardDone chan interface{}
				var wardHeartbeat <-chan interface{}
				//  启动被监控者
				startWard := func() {
					wardDone = make(chan interface{})
					// We want the ward goroutine to halt if either the steward is halted (done channel),
					//or the steward wants to halt the ward goroutine (warDone channel)
					wardHeartbeat = startGoroutine(util.OrDone(wardDone, done), timeout/2)
				}
				startWard()
				pulse := time.Tick(pulseInterval)

			monitorLoop:
				for {
					timeoutSignal := time.After(timeout)
					for {
						select {
						// steward send out pulses
						case <-pulse:
							select {
							case heartbeat <- struct{}{}:
								log.Println("steward: send out heartbeat")
							default:
							}
						// receive from the ward pulse, then continue monitoring
						case <-wardHeartbeat:
							continue monitorLoop
						case <-timeoutSignal:
							log.Println("steward: ward unhealthy; restarting")
							close(wardDone)
							startWard()
							continue monitorLoop
						case <-done:
							return
						}
					}
				}
			}()

			return heartbeat
		}
	}

	log.SetOutput(os.Stdout)
	log.SetFlags(log.Ltime | log.LUTC)

	doWork := func(done <-chan interface{}, _ time.Duration) <-chan interface{} {
		log.Println("ward: Hello, i'm irresponsible")
		go func() {
			<-done
			log.Println("ward: i'm halting")
		}()
		return nil
	}

	doWorkWithSteward := newSteward(4*time.Second, doWork)

	done := make(chan interface{})
	time.AfterFunc(9*time.Second, func() {
		log.Println("main: halting steward and ward")
		close(done)
	})

	// 读取heartbeat
	for range doWorkWithSteward(done, 2*time.Second) {
	}
	log.Println("Done")
}
