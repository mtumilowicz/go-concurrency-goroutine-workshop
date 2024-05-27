package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"
)

func makeRequest(ctx context.Context, url string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	return http.DefaultClient.Do(req)
}

type Customer struct {
	ID   int
	Name string
}

type Product struct {
	ID   int
	Name string
}

// String method for Customer struct
func (c Customer) String() string {
	return fmt.Sprintf("Customer{Name: %s}", c.Name)
}

// String method for Product struct
func (p Product) String() string {
	return fmt.Sprintf("Product{Name: %s}", p.Name)
}

func GetCustomer() Customer {
	return Customer{ID: 1, Name: "Customer1"}
}

func GetProduct() Product {
	return Product{ID: 1, Name: "Product1"}
}

func main() {
	ctx, cancelTimeout := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancelTimeout()
	ctx, cancelFunc := context.WithCancelCause(ctx)
	defer cancelFunc(nil)
	ch1 := make(chan Customer)
	ch2 := make(chan Product)
	go func() {
		time.Sleep(2 * time.Second)
		resp, err := makeRequest(ctx, "http://httpbin.org/status/200")
		if err != nil {
			cancelFunc(fmt.Errorf("in status goroutine: %w", err))
			return
		}
		if resp.StatusCode == http.StatusInternalServerError {
			cancelFunc(errors.New("bad status"))
			return
		}
		ch1 <- GetCustomer()
	}()
	go func() {
		resp, err := makeRequest(ctx, "http://httpbin.org/delay/2")
		if err != nil {
			fmt.Println("in delay goroutine:", err)
			cancelFunc(fmt.Errorf("in delay goroutine: %w", err))
			return
		}
		if resp.StatusCode == http.StatusInternalServerError {
			cancelFunc(errors.New("bad status"))
			return
		}
		ch2 <- GetProduct()
	}()
	var response1 Customer
	var response2 Product
	var err error
	for range 2 {
		select {
		case s := <-ch1:
			response1 = s
		case s := <-ch2:
			response2 = s
		case <-ctx.Done():
			err = context.Cause(ctx)
		}
	}
	if err != nil {
		fmt.Println("in main: cancelled with error", err)
		return
	}
	fmt.Printf("result: %s, %s\n", response1, response2)
}
