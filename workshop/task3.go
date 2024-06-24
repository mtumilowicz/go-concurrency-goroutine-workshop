package workshop

import (
	"fmt"
)

func Sum(numbers []int, concurrencyLevel int) int {
	length := len(numbers)
	partSize := length / concurrencyLevel
	resultChan := make(chan int, concurrencyLevel)
	errorChan := make(chan error, concurrencyLevel)

	for i := 0; i < concurrencyLevel; i++ {
		start := i * partSize
		end := start + partSize
		if i == concurrencyLevel-1 {
			end = length
		}

		go sum(numbers[start:end], resultChan)
	}

	totalSum := 0
	for i := 0; i < concurrencyLevel; i++ {
		select {
		case res := <-resultChan:
			totalSum += res
		case err := <-errorChan:
			fmt.Printf("Error: %v\n", err)
		}
	}

	return totalSum
}

func sum(numbers []int, resultChan chan int) {
	sum := 0
	for _, number := range numbers {
		sum += number
	}
	resultChan <- sum
}
