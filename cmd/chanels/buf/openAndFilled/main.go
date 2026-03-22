package main

import (
	"fmt"
	"sync"
)

func read(wg *sync.WaitGroup, readCh <-chan int) {
	defer wg.Done()

	fmt.Println("read start")
	fmt.Println(<-readCh) // Успешное чтение
	fmt.Println("read end\n")
}

func write(wg *sync.WaitGroup, writeCh chan int) {
	fmt.Println("write start")

	wg.Add(1)
	go func() {
		defer wg.Done()
		fmt.Println("reader start")
		<-writeCh // Читатель, блокировка снята
		fmt.Println("reader end")
	}()
	writeCh <- 4 // Блокировка до прихода читателя
	wg.Wait()

	fmt.Println("write end")
}

func main() {
	wg := &sync.WaitGroup{}

	readCh := make(chan int, 3)
	readCh <- 1
	readCh <- 2
	readCh <- 3

	wg.Add(1)
	read(wg, readCh)
	wg.Wait()

	write(wg, readCh)
}
