package main

import (
	"Practice9/resilient-client/client"
	"Practice9/resilient-client/server"
	"context"
	"fmt"
	"time"
)

func main() {
	srv := server.NewMockServer()

	go func() {
		fmt.Println("Mock server running on :8080")
		srv.ListenAndServe()
	}()

	time.Sleep(500 * time.Millisecond)

	c := client.NewClient()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := c.ExecutePayment(ctx, "http://localhost:8080/pay")
	if err != nil {
		fmt.Println("Final error:", err)
	}

	_ = srv.Close()
}
