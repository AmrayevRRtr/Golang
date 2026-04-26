package main

import (
	"context"
	"errors"
	"fmt"
	"math"
	"math/rand"
	"time"
)

func doSomethingUnreliable() error {
	if rand.Intn(10) < 7 {
		fmt.Println("Operation failed, retrying...")
		return errors.New("temporary failure")
	}
	fmt.Println("Operation succeeded!")
	return nil
}

type RetryConfig struct {
	maxRetries int
	baseDelay  time.Duration
	maxDelay   time.Duration
}

func Retry(ctx context.Context, cfg RetryConfig) error {
	var err error

	for attempt := 0; attempt < cfg.maxRetries; attempt++ {
		// проверка отмены перед каждой попыткой
		if ctx.Err() != nil {
			return ctx.Err()
		}

		err = doSomethingUnreliable()
		if err == nil {
			return nil
		}

		// если последняя попытка — выходим с ошибкой
		if attempt == cfg.maxRetries-1 {
			return err
		}

		// экспоненциальный backoff с Full Jitter
		backoff := cfg.baseDelay * time.Duration(math.Pow(2, float64(attempt)))
		if backoff > cfg.maxDelay {
			backoff = cfg.maxDelay
		}
		jitter := time.Duration(rand.Int63n(int64(backoff)))

		fmt.Printf("Attempt %d failed, waiting ~%v (max backoff: %v)...\n", attempt+1, jitter, backoff)
		time.Sleep(jitter)
	}
	return err
}

func main() {
	rand.Seed(time.Now().UnixNano())

	cfg := RetryConfig{
		maxRetries: 5,
		baseDelay:  500 * time.Millisecond,
		maxDelay:   5 * time.Second,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	fmt.Println("Starting retry process...\n")

	err := Retry(ctx, cfg)
	if err != nil {
		fmt.Println("\nFinal result: failed:", err)
	} else {
		fmt.Println("\nFinal result: success")
	}
}
