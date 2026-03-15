package main

import (
	"Practice5/internal/database"
	"Practice5/internal/handler"
	"Practice5/internal/repository"
	"log"
	"net/http"

	"github.com/joho/godotenv"
)

// TIP <p>To run your code, right-click the code and select <b>Run</b>.</p> <p>Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.</p>
func main() {
	godotenv.Load()

	db := database.NewMySQL()

	repo := repository.NewRepository(db)

	userHandler := handler.NewUserHandler(repo)

	mux := http.NewServeMux()

	userHandler.RegisterRoutes(mux)

	log.Println("Server running on :8080")

	http.ListenAndServe(":8080", mux)
}
