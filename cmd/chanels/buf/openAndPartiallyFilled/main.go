package main

import (
	"fmt"
	"sync"
)

func read(wg *sync.WaitGroup, ch <-chan int) {
	defer wg.Done()

	fmt.Println("read start")

	fmt.Println(<-ch) // Успешное чтение

	fmt.Println("read end\n")
}

func write(wg *sync.WaitGroup, ch chan<- int) {
	defer wg.Done()

	ch <- 3 // Успешная запись
}

func main() {
	wg := &sync.WaitGroup{}

	// Read
	wg.Add(1)

	readCh := make(chan int, 3)
	readCh <- 1
	readCh <- 2

	read(wg, readCh)
	wg.Wait()

	// Write
	wg.Add(1)
	writeCh := make(chan int, 3)
	writeCh <- 1
	writeCh <- 2

	fmt.Println("write start")
	write(wg, writeCh)

	fmt.Println(<-writeCh)
	fmt.Println("write end")
	wg.Wait()
}
