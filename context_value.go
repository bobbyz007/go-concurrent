package main

import (
	"context"
	"fmt"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(2)

	type String string

	f := func(ctx context.Context, k String) {
		defer wg.Done()
		if v := ctx.Value(k); v != nil {
			fmt.Println("Found value: ", v)
			return
		}
		fmt.Println("Key not found: ", k)
	}

	k := String("language")
	ctx := context.WithValue(context.Background(), k, "Go")
	go f(ctx, k)
	go f(ctx, String("color"))

	wg.Wait()
}

