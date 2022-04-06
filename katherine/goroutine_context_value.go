package main

import (
	"context"
	"fmt"
)

func main() {
	ProcessRequest("jane", "abc123")
}

func ProcessRequest(userId, authToken string) {
	ctx := context.WithValue(context.Background(), "userId", userId)
	ctx = context.WithValue(ctx, "authToken", authToken)
	HandleResponse(ctx)
}

func HandleResponse(ctx context.Context) {
	fmt.Printf("handling response %v (%v)", ctx.Value("userId"), ctx.Value("authToken"))
}
