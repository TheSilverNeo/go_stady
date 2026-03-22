package main

import (
	"fmt"
	"sync"
)

func read(ch <-chan int) {
	fmt.Println("read start")

	fmt.Println(<-ch) // Zero value

	fmt.Println("read end")
}

func write(ch chan<- int) {
	fmt.Println("write start")

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer func() {
			if err := recover(); err != nil {
				fmt.Println(err)
			}
		}()

		ch <- 1 // Паника. Запись в закрытый канал
	}()
	wg.Wait()
}

func main() {
	readCh := make(chan int, 3)
	close(readCh)
	read(readCh)

	writeCh := make(chan int, 3)
	close(writeCh)
	write(writeCh)
}
