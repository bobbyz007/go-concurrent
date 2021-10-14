package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	d := time.Now().Add(time.Millisecond * 100)
	ctx, cancel := context.WithDeadline(context.Background(), d)
	// 或者如下定义
	// ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond * 100)

	// 等main函数结束后，上下文会被关闭，如果不及时关闭会导致上下文ctx和其父context存在的周期比我们想要的长
	defer cancel()

	select {
	case <- time.After(time.Millisecond * 200):
		fmt.Println("overslept")
	case <- ctx.Done():
		fmt.Println(ctx.Err())
	}
}

