package main

import (
	"fmt"
	"go-concurrency-goroutine-workshop/workshop"
)

func main() {
	customer, product, err := workshop.Recommend(1)
	if err != nil {
		fmt.Println("in main: cancelled with error", err)
		return
	}
	fmt.Printf("result: %s, %s\n", customer, product)
}
