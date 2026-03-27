package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"
)

func randomTymeWork() {
	time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
}

func predictableTimeWork(ctx context.Context, ch chan int) error {
	go func() {
		for i := range 1000 {
			randomTymeWork()

			select {
			case ch <- i:
			case <-ctx.Done():
				return
			default:
			}
		}
		close(ch)
	}()

	for {
		select {
		case v, ok := <-ch:
			if !ok {
				fmt.Println("channel closed")
				return nil
			}
			fmt.Printf("read from ch: %d\n", v)
		case <-ctx.Done():
			fmt.Println("context done")
			return ctx.Err()
		}
	}
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	ch := make(chan int)
	if err := predictableTimeWork(ctx, ch); err != nil {
		fmt.Println(err)
	}
}
