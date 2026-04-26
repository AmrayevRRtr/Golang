package main

import (
	"context"
	"encoding/json"
	"fmt"
	_ "io"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	ProcessingStatus = "PROCESSING"
	LockTTL          = 5 * time.Minute
	CacheTTL         = 24 * time.Hour
)

type CachedResponse struct {
	StatusCode int               `json:"status_code"`
	Body       []byte            `json:"body"`
	Headers    map[string]string `json:"headers"`
}

type RedisStore struct {
	client *redis.Client
	ctx    context.Context
}

func NewRedisStore(addr string) *RedisStore {
	rdb := redis.NewClient(&redis.Options{Addr: addr})
	return &RedisStore{
		client: rdb,
		ctx:    context.Background(),
	}
}

func (s *RedisStore) StartProcessing(key string) (isNew bool, cached *CachedResponse) {
	ok, err := s.client.SetNX(s.ctx, key, ProcessingStatus, LockTTL).Result()
	if err != nil {
		return false, nil
	}

	if ok {
		return true, nil
	}

	val, err := s.client.Get(s.ctx, key).Result()
	if err != nil || val == ProcessingStatus {
		return false, nil
	}

	var result CachedResponse
	if err := json.Unmarshal([]byte(val), &result); err == nil {
		return false, &result
	}

	return false, nil
}

func (s *RedisStore) Finish(key string, resp CachedResponse) {
	data, _ := json.Marshal(resp)
	s.client.Set(s.ctx, key, data, CacheTTL)
}

func IdempotencyMiddleware(store *RedisStore, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := r.Header.Get("X-Idempotency-Key")
		if key == "" {
			http.Error(w, "X-Idempotency-Key header is required", http.StatusBadRequest)
			return
		}

		isNew, cached := store.StartProcessing(key)

		if !isNew {
			if cached != nil {
				w.Header().Set("X-Cache", "HIT")
				for k, v := range cached.Headers {
					w.Header().Set(k, v)
				}
				w.WriteHeader(cached.StatusCode)
				w.Write(cached.Body)
				return
			}
			http.Error(w, "Duplicate request in progress", http.StatusConflict)
			return
		}

		recorder := httptest.NewRecorder()
		next.ServeHTTP(recorder, r)

		respToCache := CachedResponse{
			StatusCode: recorder.Code,
			Body:       recorder.Body.Bytes(),
			Headers:    map[string]string{"Content-Type": recorder.Header().Get("Content-Type")},
		}

		store.Finish(key, respToCache)

		for k, v := range respToCache.Headers {
			w.Header().Set(k, v)
		}
		w.WriteHeader(recorder.Code)
		w.Write(recorder.Body.Bytes())
	})
}

func main() {
	store := NewRedisStore("localhost:6379")

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Processing complex business logic...")
		time.Sleep(5 * time.Second)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"success","transaction_id":"tx_999"}`))
	})

	fmt.Println("Server running on :8080...")
	http.ListenAndServe(":8080", IdempotencyMiddleware(store, handler))
}
