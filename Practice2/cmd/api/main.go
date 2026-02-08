package main

import (
	"Practice2/internal/handlers"
	"Practice2/internal/middleware"
	"Practice2/internal/storage"
	"log"
	"net/http"
)

func main() {

	store := storage.NewStore()
	handler := handlers.NewHandler(store)

	mux := http.NewServeMux()
	mux.Handle("/tasks", http.HandlerFunc(handler.Tasks))

	var h http.Handler = mux

	h = middleware.Logging(h)
	h = middleware.APIKey(h)

	log.Println("Server started at :8080")

	log.Fatal(http.ListenAndServe(":8080", h))

}
