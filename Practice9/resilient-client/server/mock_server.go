package server

import (
	"encoding/json"
	"net/http"
)

func NewMockServer() *http.Server {
	attempt := 0

	mux := http.NewServeMux()

	mux.HandleFunc("/pay", func(w http.ResponseWriter, r *http.Request) {
		attempt++

		if attempt < 4 {
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte("temporary error"))
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"status": "success",
		})
	})

	return &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
}
