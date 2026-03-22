package main

import "fmt"

func addValue(ch chan<- int, val []int) {
	for _, v := range val {
		ch <- v
	}
	close(ch) // канал завершается для остановки цикла
}

func loopByChanel() {
	ch := make(chan int)
	slice := make([]int, 0, 10)

	for i := 0; i < 10; i++ {
		slice = append(slice, i)
	}
	go addValue(ch, slice)

	for v := range ch { // Цикл работает до закрытия канала
		fmt.Println(v)
	}

	/** Альтернатива
	for {
	    value, ok := <-ch
	    if !ok {
	        break
	    }
	    fmt.Println(value)
	}
	*/
}

func main() {
	loopByChanel()
}
