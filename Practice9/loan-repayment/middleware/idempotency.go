package middleware

import (
	"Practice9/loan-repayment/store"
	"bytes"
	"context"
	"net/http"
)

type responseWriter struct {
	http.ResponseWriter
	body       *bytes.Buffer
	statusCode int
}

func (w *responseWriter) WriteHeader(code int) {
	w.statusCode = code
	w.ResponseWriter.WriteHeader(code)
}

func (w *responseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func IdempotencyMiddlewareRedis(st *store.RedisStore, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		ctx := context.Background()
		key := r.Header.Get("Idempotency-Key")

		if key == "" {
			http.Error(w, "Idempotency-Key required", http.StatusBadRequest)
			return
		}

		val, err := st.Get(ctx, key)

		if err == nil {
			if val == "processing" {
				http.Error(w, "Request in progress", http.StatusConflict)
				return
			}

			w.WriteHeader(http.StatusOK)
			w.Write([]byte(val))
			return
		}

		ok, err := st.Start(ctx, key)
		if err != nil {
			http.Error(w, "Redis error", 500)
			return
		}

		if !ok {
			http.Error(w, "Request in progress", http.StatusConflict)
			return
		}

		rec := &responseWriter{
			ResponseWriter: w,
			body:           &bytes.Buffer{},
		}

		next.ServeHTTP(rec, r)

		st.Save(ctx, key, rec.body.String())
	})
}
