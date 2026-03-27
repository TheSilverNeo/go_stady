package main

import (
	"fmt"
	"time"
)

// Пример 1: Базовое поведение nil канала (DEADLOCK!)
func example1_BasicNilBehavior() {
	fmt.Println("\n=== Пример 1: Базовое поведение nil канала ===")
	fmt.Println("Закомментирован! Если раскомментировать - DEADLOCK!")
	fmt.Println("✗ Отправка в nil:   ch <- x     → вечное блокирование")
	fmt.Println("✗ Получение из nil: <-ch        → вечное блокирование")
	fmt.Println("✗ close(nil)                    → паника")
}

// Пример 2: Отключение веток в select - ОСНОВНОЕ ПРИМЕНЕНИЕ!
func example2_DisableInSelect() {
	fmt.Println("\n=== Пример 2: Отключение веток в select ===")

	ch1 := make(chan string)
	ch2 := make(chan string)

	go func() {
		time.Sleep(100 * time.Millisecond)
		ch1 <- "сообщение из ch1"
	}()

	go func() {
		time.Sleep(200 * time.Millisecond)
		ch2 <- "сообщение из ch2"
	}()

	for i := 0; i < 2; i++ {
		select {
		case msg := <-ch1:
			fmt.Printf("Итерация %d: получено из ch1: %s\n", i, msg)
			ch1 = nil // ⚡ ОТКЛЮЧАЕМ ch1! Больше в этот case не зайдём

		case msg := <-ch2:
			fmt.Printf("Итерация %d: получено из ch2: %s\n", i, msg)
		}
	}

	fmt.Println("✓ ch1 отключен после первого получения, используется только ch2")
}

// Пример 3: Мержинг двух каналов с отключением
func example3_MergeChannels() {
	fmt.Println("\n=== Пример 3: Мержинг каналов с отключением ===")

	ch1 := make(chan int)
	ch2 := make(chan int)

	// Отправитель ch1
	go func() {
		ch1 <- 1
		ch1 <- 2
		close(ch1)
	}()

	// Отправитель ch2
	go func() {
		ch2 <- 10
		ch2 <- 20
		ch2 <- 30
		close(ch2)
	}()

	// Получатель: когда один канал закончится, отключаем его (nil)
	for ch1 != nil || ch2 != nil {
		select {
		case v, ok := <-ch1:
			if !ok {
				fmt.Println("  ch1 закрыт, отключаем его")
				ch1 = nil // ⚡ Отключаем закрытый канал
				continue
			}
			fmt.Printf("  Из ch1: %d\n", v)

		case v, ok := <-ch2:
			if !ok {
				fmt.Println("  ch2 закрыт, отключаем его")
				ch2 = nil // ⚡ Отключаем закрытый канал
				continue
			}
			fmt.Printf("  Из ch2: %d\n", v)
		}
	}

	fmt.Println("✓ Оба канала обработаны полностью")
}

// Пример 4: Условное отключение (динамическое управление)
func example4_ConditionalDisable(enableFast bool) {
	fmt.Println(fmt.Sprintf("\n=== Пример 4: Условное отключение (enableFast=%v) ===", enableFast))

	normalCh := make(chan string)
	fastCh := make(chan string)

	// Если fast отключен, переменная fastCh станет nil
	if !enableFast {
		fastCh = nil // ⚡ Динамически отключаем fast ветку
	}

	go func() {
		normalCh <- "обработка на нормальной скорости"
		close(normalCh)
	}()

	if enableFast {
		go func() {
			fastCh <- "ускоренная обработка"
			close(fastCh)
		}()
	}

	processed := false
	for {
		select {
		case msg, ok := <-fastCh:
			if !ok {
				fastCh = nil
				continue
			}
			fmt.Printf("  ⚡ Fast: %s\n", msg)
			processed = true

		case msg, ok := <-normalCh:
			if !ok {
				normalCh = nil
			}
			if normalCh == nil && fastCh == nil {
				fmt.Println("  ✓ Обе ветки завершены")
				return
			}
			fmt.Printf("  ✓ Normal: %s\n", msg)
			processed = true

		default:
			if fastCh == nil && normalCh == nil {
				return
			}
		}

		if !processed {
			break
		}
	}
}

// Пример 5: Практический обработчик событий
func example5_EventProcessor() {
	fmt.Println("\n=== Пример 5: Обработчик событий с несколькими источниками ===")

	userEvents := make(chan string, 2)
	systemEvents := make(chan string, 2)

	// Имитируем события
	go func() {
		userEvents <- "click button"
		userEvents <- "input text"
		close(userEvents)
	}()

	go func() {
		time.Sleep(50 * time.Millisecond)
		systemEvents <- "memory warning"
		systemEvents <- "low battery"
		close(systemEvents)
	}()

	// Обработчик - когда один источник закончится, отключаем его
	eventCount := 0
	for userEvents != nil || systemEvents != nil {
		select {
		case event, ok := <-userEvents:
			if !ok {
				fmt.Println("  [User Events] закрыты, отключаем")
				userEvents = nil
				continue
			}
			fmt.Printf("  [USER EVENT #%d] %s\n", eventCount+1, event)
			eventCount++

		case event, ok := <-systemEvents:
			if !ok {
				fmt.Println("  [System Events] закрыты, отключаем")
				systemEvents = nil
				continue
			}
			fmt.Printf("  [SYSTEM EVENT #%d] %s\n", eventCount+1, event)
			eventCount++
		}
	}

	fmt.Printf("✓ Обработано всего событий: %d\n", eventCount)
}

// Пример 6: Разница между nil и closed каналом
func example6_NilVsClosed() {
	fmt.Println("\n=== Пример 6: nil vs closed канал ===")

	normalCh := make(chan int)
	var nilCh chan int // nil по умолчанию

	close(normalCh)

	fmt.Println("Попытка получить из закрытого канала:")
	v, ok := <-normalCh
	fmt.Printf("  Значение: %d, ok: %v (нулевое значение, ok=false)\n", v, ok)

	fmt.Println("\nПопытка получить из nil канала в select:")
	select {
	case v, ok := <-nilCh:
		fmt.Printf("  Никогда сюда не попадём! v=%d, ok=%v\n", v, ok)
	case <-time.After(100 * time.Millisecond):
		fmt.Println("  ✓ Select пропустил nil канал и пошёл в timeout")
	}
}

func main() {
	fmt.Println("╔══════════════════════════════════════════════════════════╗")
	fmt.Println("║     Nil каналы в Go: применение и особенности            ║")
	fmt.Println("╚══════════════════════════════════════════════════════════╝")

	// Пример 1: Базовое поведение (закомментирован, вызывает deadlock)
	example1_BasicNilBehavior()

	// Пример 2: Отключение веток в select (ОСНОВНОЕ ПРИМЕНЕНИЕ)
	example2_DisableInSelect()

	// Пример 3: Мержинг каналов с отключением
	example3_MergeChannels()

	// Пример 4: Условное отключение
	example4_ConditionalDisable(true)
	example4_ConditionalDisable(false)

	// Пример 5: Практический обработчик событий
	example5_EventProcessor()

	// Пример 6: Разница nil vs closed
	example6_NilVsClosed()

	fmt.Println("\n✓ Все примеры выполнены успешно!")
}
