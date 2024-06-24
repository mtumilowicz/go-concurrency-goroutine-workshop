package workshop

import (
	"context"
	"errors"
	"time"
)

func Compute(ctx context.Context) (int, error) {
	resultPromise := make(chan int)
	defer close(resultPromise)
	go func() {
		time.Sleep(2 * time.Second)
		resultPromise <- 10
	}()
	select {
	case result := <-resultPromise:
		return result, nil
	case <-ctx.Done():
		return 0, ctx.Err()
	}
}

func Compute2() (int, error) {
	resultPromise := make(chan int)
	defer close(resultPromise)

	go func() {
		result := timeConsumingComputation()
		resultPromise <- result
	}()

	select {
	case result := <-resultPromise:
		return result, nil
	case <-time.After(1 * time.Second):
		return 0, errors.New("timeout occurred")
	}
}

func timeConsumingComputation() int {
	time.Sleep(2 * time.Second)
	return 10
}
