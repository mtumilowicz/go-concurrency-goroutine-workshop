package main

import (
	"fmt"
	"go-concurrency-goroutine-workshop/workshop"
	"runtime"
	"time"
)

func main() {
	concurrencyLevel := runtime.NumCPU()

	// sum 1...20 = 210
	numbers := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20}

	totalSum := workshop.Sum(numbers, concurrencyLevel)
	_, err := workshop.Compute2()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("Total sum: %d\n", totalSum)
}

func main2() {
	ch1 := make(chan string)
	ch2 := make(chan string)

	// Goroutine to send data to ch1
	go func() {
		for {
			ch1 <- "Message from channel 1"
			time.Sleep(100 * time.Second)
		}
	}()

	// Goroutine to send data to ch2
	go func() {
		for {
			ch2 <- "Message from channel 2"
			time.Sleep(200 * time.Second)
		}
	}()

	// Using for-select loop to read from both channels
	for {
		select {
		case msg1 := <-ch1:
			fmt.Println("Received from ch1:", msg1)
		case msg2 := <-ch2:
			fmt.Println("Received from ch2:", msg2)
		}
		fmt.Println("loop")
	}
}
