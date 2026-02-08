package main

import (
	"Practice2/internal/handlers"
	"Practice2/internal/middleware"
	"Practice2/internal/storage"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {

	store := storage.NewStore()
	handler := handlers.NewHandler(store)

	mux := http.NewServeMux()
	mux.Handle("/tasks", http.HandlerFunc(handler.Tasks))

	var h http.Handler = mux
	h = middleware.Logging(h)
	h = middleware.APIKey(h)

	server := &http.Server{
		Addr:    ":8080",
		Handler: h,
	}

	go func() {
		log.Println("Server started at :8080")

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server stopped")
}
