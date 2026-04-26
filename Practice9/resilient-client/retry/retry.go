package retry

import (
	"math"
	"math/rand"
	"net/http"
	"time"
)

func IsRetryable(resp *http.Response, err error) bool {
	if err != nil {
		return true
	}

	if resp == nil {
		return true
	}

	switch resp.StatusCode {
	case 429, 500, 502, 503, 504:
		return true
	case 401, 404:
		return false
	default:
		return false
	}
}

type Config struct {
	BaseDelay  time.Duration
	MaxDelay   time.Duration
	MaxRetries int
}

func CalculateBackoff(attempt int, cfg Config) time.Duration {
	backoff := cfg.BaseDelay * time.Duration(math.Pow(2, float64(attempt)))

	if backoff > cfg.MaxDelay {
		backoff = cfg.MaxDelay
	}
	
	return time.Duration(rand.Int63n(int64(backoff)))
}
