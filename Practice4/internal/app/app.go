package app

import (
	"Practice4/internal/handler"
	"Practice4/internal/middleware"
	"Practice4/internal/repository"
	"Practice4/internal/repository/mysql"
	"Practice4/internal/usecase"
	"Practice4/pkg/modules"
	"context"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

func Run() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	timeoutStr := os.Getenv("DB_EXEC_TIMEOUT")
	timeoutInt, err := strconv.Atoi(timeoutStr)
	if err != nil {
		timeoutInt = 5
	}

	cfg := &modules.MySQLConfig{
		Host: os.Getenv("DB_HOST"), Port: os.Getenv("DB_PORT"),
		Username: os.Getenv("DB_USER"), Password: os.Getenv("DB_PASSWORD"),
		DBName: os.Getenv("DB_NAME"), SSLMode: os.Getenv("DB_SSLMODE"),
		ExecTimeout: time.Duration(timeoutInt) * time.Second,
	}

	ctx := context.Background()
	db := mysql.NewMySQLDialect(ctx, cfg)
	mysqlRepo := mysql.NewUserRepository(db)
	repositories := repository.NewRepositories(mysqlRepo)
	uc := usecase.NewUserUsecase(repositories.UserRepository)
	h := handler.NewUserHandler(uc)

	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	log.Println("Server started at :8080")

	log.Fatal(http.ListenAndServe(":8080", middleware.Logging(middleware.Auth(mux))))

}
