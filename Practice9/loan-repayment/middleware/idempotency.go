package middleware

import (
	"Practice9/loan-repayment/store"
	"bytes"
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

func IdempotencyMiddleware(st *store.Store, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		key := r.Header.Get("Idempotency-Key")

		if key == "" {
			http.Error(w, "Idempotency-Key required", http.StatusBadRequest)
			return
		}

		if existing, ok := st.Get(key); ok {

			if existing.Status == store.Processing {
				http.Error(w, "Request already in progress", http.StatusConflict)
				return
			}

			if existing.Status == store.Completed {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(existing.StatusCode)
				w.Write(existing.Body)
				return
			}
		}

		ok := st.Start(key)
		if !ok {
			http.Error(w, "Request already in progress", http.StatusConflict)
			return
		}

		println("Processing started for key:", key)

		rec := &responseWriter{
			ResponseWriter: w,
			body:           &bytes.Buffer{},
		}

		next.ServeHTTP(rec, r)

		st.Save(key, rec.statusCode, rec.body.Bytes())
	})
}
