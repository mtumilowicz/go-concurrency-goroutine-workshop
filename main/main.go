package main

import (
	"context"
	"fmt"
	"go-concurrency-goroutine-workshop/workshop"
	"time"
)

func main2() {
	customer, product, err := workshop.Recommend(1)
	if err != nil {
		fmt.Println("in main: cancelled with error", err)
		return
	}
	fmt.Printf("result: %s, %s\n", customer, product)
}

func main() {
	ctx, cancelTimeout := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancelTimeout()
	sum, err := workshop.Compute(ctx)
	if err != nil {
		fmt.Println("in main2: cancelled with error", err)
		return
	}
	fmt.Printf("result: %d", sum)
}
