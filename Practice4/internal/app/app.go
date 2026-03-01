package app

import (
	"Practice4/internal/handler"
	"Practice4/internal/middleware"
	"Practice4/internal/repository"
	"Practice4/internal/repository/mysql"
	"Practice4/internal/repository/mysql/users"
	"Practice4/internal/usecase"
	"Practice4/pkg/modules"
	"context"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	httpSwagger "github.com/swaggo/http-swagger"

	"log"
	"net/http"
	"os"
	"strconv"
	"time"
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
	mysqlRepo := users.NewUserRepository(db)
	repositories := repository.NewRepositories(mysqlRepo)
	uc := usecase.NewUserUsecase(repositories.UserRepository)
	h := handler.NewUserHandler(uc)

	mainMux := http.NewServeMux()
	protectedMux := http.NewServeMux()

	mainMux.Handle("/swagger/", httpSwagger.WrapHandler)
	mainMux.HandleFunc("/health", h.Health)

	h.RegisterRoutes(protectedMux)

	mainMux.Handle("/users", middleware.Auth(protectedMux))
	mainMux.Handle("/users/", middleware.Auth(protectedMux))

	server := &http.Server{
		Addr:    ":" + os.Getenv("SERVER_PORT"),
		Handler: middleware.Logging(mainMux),
	}

	go func() {
		log.Println("Server started at :" + os.Getenv("SERVER_PORT"))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	log.Println("Shutting down gracefully...")

	ctxTimeout, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctxTimeout); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	if err := db.DB.Close(); err != nil {
		log.Println("Error closing database:", err)
	}

	log.Println("Server exited properly")
}
