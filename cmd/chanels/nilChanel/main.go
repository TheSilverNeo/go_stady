package main

import "fmt"

func main() {
	ch := make(chan struct{})
	ch = nil

	ch <- struct{}{}  // Вечное ожидание
	fmt.Println(<-ch) // Вечное ожидание
}
