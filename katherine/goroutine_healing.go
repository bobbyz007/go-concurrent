package main

import (
	"go-concurrent/katherine/util"
	"log"
	"os"
	"time"
)

func main() {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.Ltime | log.LUTC)

	// 模拟无心跳消息，认为该goroutine有问题
	doWork := func(done <-chan interface{}, _ time.Duration) <-chan interface{} {
		log.Println("ward: Hello, i'm irresponsible")
		go func() {
			<-done
			log.Println("ward: i'm halting")
		}()
		return nil
	}

	// 监控的超时时间是4s，如果4s没有监听到ward心跳消息，则会重启ward
	doWorkWithSteward := util.NewSteward(4*time.Second, doWork)

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
