package main

import (
	"fmt"
	"sync"
	"time"
)

func read(wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Println("read start")
	ch := make(chan int, 4)

	go func() {
		fmt.Println("write to channel")
		time.Sleep(1 * time.Second)
		ch <- 1 // Запись в канал. Снятие блокировки
	}()

	<-ch // Блокировка до прихода писателя

	fmt.Println("read end\n")
}

func write(wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Println("write start")

	ch := make(chan int, 4)
	ch <- 1 // успешная запись
	fmt.Println(<-ch)

	fmt.Println("write end")
}

func main() {
	wg := &sync.WaitGroup{}

	wg.Add(1)
	read(wg)
	wg.Wait()

	wg.Add(1)
	write(wg)
	wg.Wait()
}
