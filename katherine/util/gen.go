package util

import (
	"log"
	"time"
)

// 无限重复
func Repeat(done <-chan interface{}, values ...interface{}) <-chan interface{} {
	valueStream := make(chan interface{})
	go func() {
		defer close(valueStream)

		// infinite loop
		for {
			for _, v := range values {
				select {
				case <-done:
					return
				case valueStream <- v:
				}
			}
		}
	}()
	return valueStream
}

func RepeatFn(done <-chan interface{}, fn func() interface{}) <-chan interface{} {
	valueStream := make(chan interface{})
	go func() {
		defer close(valueStream)

		// infinite loop
		for {
			select {
			case <-done:
				return
			case valueStream <- fn():
			}
		}
	}()
	return valueStream
}

// 返回valuedStream前num个元素
func Take(done <-chan interface{}, valueStream <-chan interface{}, num int) chan interface{} {
	takeStream := make(chan interface{})
	go func() {
		defer close(takeStream)

		for i := 0; i < num; i++ {
			select {
			case <-done:
				return
			case takeStream <- <-valueStream:
			}
		}
	}()
	return takeStream
}

func OrDone(done, c <-chan interface{}) <-chan interface{} {
	valStream := make(chan interface{})
	go func() {
		defer close(valStream)
		for {
			select {
			case <-done:
				return
			// 检查c是否关闭
			case v, ok := <-c:
				log.Printf("ordone: %v\n", ok)
				if ok == false {
					return
				}
				select {
				case valStream <- v:
				case <-done: // 保证立即响应done
				}
			}
		}
	}()
	return valStream
}

// in 分发到 out1,out2,...
func Tee(done <-chan interface{}, in <-chan interface{}) (_, _ <-chan interface{}) {
	out1 := make(chan interface{})
	out2 := make(chan interface{})
	go func() {
		defer close(out1)
		defer close(out2)
		for val := range OrDone(done, in) {
			var out1, out2 = out1, out2
			for i := 0; i < 2; i++ {
				select {
				case <-done:
				case out1 <- val:
					out1 = nil
				case out2 <- val:
					out2 = nil
				}
			}
		}
	}()
	return out1, out2
}

// 多个channel 汇总到一个channel
func Bridge(done <-chan interface{}, chanStream <-chan <-chan interface{}) <-chan interface{} {
	valStream := make(chan interface{})
	go func() {
		defer close(valStream)

		for {
			var stream <-chan interface{}
			select {
			case maybeStream, ok := <-chanStream:
				log.Printf("bridge: %v\n", ok)
				if ok == false {
					return
				}
				stream = maybeStream
			case <-done:
				return
			}

			// 如果or 或 stream关闭，则不执行循环，防止valStream保存nil
			for val := range OrDone(done, stream) {
				select {
				case valStream <- val:
				case <-done:
					return
				}
			}
		}
	}()
	return valStream
}

// 基于heartbeat模式，代表被监控的goroutine
type StartGoroutineFn func(done <-chan interface{}, pulseInterval time.Duration) (heartbeat <-chan interface{})

// 返回类型是 startGoroutineFn，表示steward本身也是可以被监控的
// timeout: ward的超时时间。 如果在timeout时间内 获取不到ward的心跳信息，则会重启ward
func NewSteward(timeout time.Duration, startGoroutine StartGoroutineFn) StartGoroutineFn {
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
				wardHeartbeat = startGoroutine(OrDone(wardDone, done), timeout/2)
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
