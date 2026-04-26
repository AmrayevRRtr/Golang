package main

import (
	"Practice9/loan-repayment/handler"
	"Practice9/loan-repayment/middleware"
	"Practice9/loan-repayment/store"
	"io"
	"net/http"
	"sync"
	"time"
)

func main() {

	st := store.NewStore()

	mux := http.NewServeMux()
	mux.Handle("/pay", http.HandlerFunc(handler.PayHandler))

	handlerWithMW := middleware.IdempotencyMiddleware(st, mux)

	server := &http.Server{
		Addr:    ":8080",
		Handler: handlerWithMW,
	}

	go func() {
		println("Server running on :8080")
		server.ListenAndServe()
	}()

	time.Sleep(500 * time.Millisecond)

	key := "abc-123-idempotency-key"

	var wg sync.WaitGroup

	for i := 0; i < 8; i++ {
		wg.Add(1)

		go func(i int) {
			defer wg.Done()

			req, _ := http.NewRequest("POST", "http://localhost:8080/pay", nil)
			req.Header.Set("Idempotency-Key", key)

			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				println("err:", err.Error())
				return
			}
			defer resp.Body.Close()

			body, _ := io.ReadAll(resp.Body)

			println("Request", i, "status:", resp.StatusCode, "body:", string(body))
		}(i)
	}

	wg.Wait()
}
