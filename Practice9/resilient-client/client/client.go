package client

import (
	"Practice9/resilient-client/retry"
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Client struct {
	httpClient *http.Client
	cfg        retry.Config
}

func NewClient() *Client {
	return &Client{
		httpClient: &http.Client{},
		cfg: retry.Config{
			BaseDelay:  500 * time.Millisecond,
			MaxDelay:   5 * time.Second,
			MaxRetries: 5,
		},
	}
}

func (c *Client) ExecutePayment(ctx context.Context, url string) error {
	var lastErr error

	for attempt := 0; attempt < c.cfg.MaxRetries; attempt++ {

		if ctx.Err() != nil {
			return ctx.Err()
		}

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		if err != nil {
			return err
		}

		resp, err := c.httpClient.Do(req)

		var body []byte
		if resp != nil && resp.Body != nil {
			body, _ = io.ReadAll(resp.Body)
			resp.Body.Close()
		}

		if err == nil && resp.StatusCode == 200 {
			fmt.Printf("Attempt %d: Success! Response: %s\n", attempt+1, string(body))
			return nil
		}

		lastErr = err

		if !retry.IsRetryable(resp, err) {
			if err != nil {
				return err
			}
			return fmt.Errorf("non-retryable status: %d", resp.StatusCode)
		}

		if attempt == c.cfg.MaxRetries-1 {
			break
		}

		delay := retry.CalculateBackoff(attempt, c.cfg)

		fmt.Printf("Attempt %d failed: waiting %v...\n", attempt+1, delay)

		select {
		case <-time.After(delay):
		case <-ctx.Done():
			return ctx.Err()
		}
	}

	return fmt.Errorf("payment failed after retries: %w", lastErr)
}
