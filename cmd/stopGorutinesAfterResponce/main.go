package main

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// Функция с непредсказуемым временем исполнения
func slowFunction(id int) (int, error) {
	// Имитируем разное время выполнения (100–1000 мс)
	delay := time.Duration(rand.Intn(900) + 100)
	if delay%2 == 0 {
		return 0, fmt.Errorf("slow function error%d", id)
	}
	time.Sleep(delay * time.Millisecond)

	// Возвращаем ID горутины как результат
	return id, nil
}

func main() {
	const numGoroutines = 100

	// Создаём контекст с возможностью отмены
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Канал для передачи результата (буферизованный на 1)
	resultCh := make(chan int, 1)

	// WaitGroup для ожидания завершения горутин
	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	retryAttempt := 2
	// Запускаем горутины
	for i := 0; i < numGoroutines; i++ {
		go func(goroutineID int, retryAttempt int) {
			defer wg.Done()

			// Проверяем, не отменён ли контекст
			select {
			case <-ctx.Done():
				fmt.Printf("Горутина %d: отменена\n", goroutineID)
				return
			default:
			}

			// Выполняем функцию
			result, err := slowFunction(goroutineID)
			if err != nil {
				for k := 0; k < retryAttempt; i++ {
					fmt.Printf("Горутина %d: перезапуск: %d\n", goroutineID, k)
					result, err = slowFunction(goroutineID)
				}
				fmt.Printf("Горутина %d: ошибка: %v\n", goroutineID, err)
				return
			}

			// Отправляем результат в канал (только если ещё не отправлен)
			select {
			case resultCh <- result:
				fmt.Printf("Горутина %d: отправила результат %d\n", goroutineID, result)
				// Отменяем контекст — останавливаем остальные горутины
				cancel()
			default:
				// Канал уже заполнен — значит, результат уже отправлен другой горутиной
				fmt.Printf("Горутина %d: результат уже получен, пропускаем отправку\n", goroutineID)
			}
		}(i, retryAttempt)
	}

	// Ждём либо результата, либо отмены контекста
	select {
	case result := <-resultCh:
		fmt.Printf("Получен первый результат: %d\n", result)
	case <-ctx.Done():
		fmt.Println("Операция отменена")
	}

	// Ждём завершения всех горутин
	wg.Wait()
	fmt.Println("Все горутины завершены")
}
