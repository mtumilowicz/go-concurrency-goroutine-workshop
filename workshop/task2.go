package workshop

import (
	"context"
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
