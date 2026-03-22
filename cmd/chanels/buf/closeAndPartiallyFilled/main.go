package main

import (
	"fmt"
	"sync"
)

func read(ch <-chan int) {
	fmt.Println("read start")

	fmt.Println(<-ch) // Успешное чтение

	fmt.Println("read end\n")
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

		ch <- 4 // Паника. Запись в закрытый канал
	}()

	wg.Wait()
	fmt.Println("write end")
}

func main() {
	readCh := make(chan int, 3)
	readCh <- 1
	readCh <- 2
	close(readCh)
	read(readCh)

	writeCh := make(chan int, 3)
	writeCh <- 1
	writeCh <- 2
	close(writeCh)
	write(writeCh)
}
