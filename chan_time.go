package main

import (
	"fmt"
	"time"
)

var timeout <-chan time.Time

func main() {
	// 创建了Time定时器管道
	timeout = time.After(time.Second * 3)
	// 从管道中读取3s后的时间
	fmt.Println(<-timeout)
}
