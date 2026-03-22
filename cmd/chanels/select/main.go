package main

import (
	"fmt"
	"sync"
	"time"
)

func nilChanel() {
	ch1 := make(chan struct{})
	ch2 := make(chan struct{})

	go func() {
		ch2 <- struct{}{}
	}()

	go func() {
		ch1 <- struct{}{}
	}()

	for i := 0; i < 4; i++ {
		select {
		case <-ch1:
			fmt.Println("case 1")
			ch2 = nil // Присваиваем каналу nil. В case 2 больше не попасть
		case <-ch2:
			fmt.Println("case 2")
		case <-time.After(2 * time.Second):
			fmt.Println("After 2 seconds")
			return
		}
	}
}

func test(wg *sync.WaitGroup, ch chan int) {
	defer wg.Done()
	ch <- 1
}

func main() {
	ch := make(chan int)
	wg := &sync.WaitGroup{}

	wg.Add(1)
	go test(wg, ch)

	fmt.Println(<-ch)
	wg.Wait()

	nilChanel()
}
