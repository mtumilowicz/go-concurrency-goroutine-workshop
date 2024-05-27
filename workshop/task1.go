package workshop

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"
)

type Customer struct {
	ID   int
	Name string
}

func (c Customer) String() string {
	return fmt.Sprintf("Customer{ID: %d, Name: %s}", c.ID, c.Name)
}

type Product struct {
	ID   int
	Name string
}

func (p Product) String() string {
	return fmt.Sprintf("Product{ID: %d, Name: %s}", p.ID, p.Name)
}

func makeRequest(ctx context.Context, url string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	return http.DefaultClient.Do(req)
}

func fetchCustomer(ctx context.Context, promise chan<- Customer, cancelFunc context.CancelCauseFunc) {
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
	promise <- Customer{ID: 1, Name: "Customer1"}
}

func fetchProduct(ctx context.Context, promise chan<- Product, cancelFunc context.CancelCauseFunc) {
	resp, err := makeRequest(ctx, "http://httpbin.org/delay/2")
	if err != nil {
		cancelFunc(fmt.Errorf("in delay goroutine: %w", err))
		return
	}
	if resp.StatusCode == http.StatusInternalServerError {
		cancelFunc(errors.New("bad status"))
		return
	}
	promise <- Product{ID: 1, Name: "Product1"}
}

func Recommend(customerId int) (Customer, Product, error) {
	ctx, cancelTimeout := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancelTimeout()
	ctx, cancelFunc := context.WithCancelCause(ctx)
	defer cancelFunc(nil)
	customerPromise := make(chan Customer)
	defer close(customerPromise)
	productPromise := make(chan Product)
	defer close(productPromise)
	go fetchCustomer(ctx, customerPromise, cancelFunc)
	go fetchProduct(ctx, productPromise, cancelFunc)
	var customer Customer
	var product Product
	var err error
	for range 2 {
		select {
		case customer = <-customerPromise:
		case product = <-productPromise:
		case <-ctx.Done():
			err = context.Cause(ctx)
		}
	}
	return customer, product, err
}
