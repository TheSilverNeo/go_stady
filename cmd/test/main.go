package main

import (
	"context"
	"fmt"
)

func main() {
	ctx := context.Background()
	withCancel, cancel := context.WithCancel(ctx)
	cancel()
	cancel()
	_ = withCancel
	fmt.Println("hello world")
}
