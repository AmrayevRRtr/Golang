package app

import (
	"fmt"
	"log"
	"os"
	v1 "practice7/internal/controller/http/v1"
	"practice7/internal/entity"
	"practice7/internal/pkg/mysql"
	"practice7/internal/usecase"
	"practice7/internal/usecase/repo"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/gorm/logger"
)

func Run() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	db, err := mysql.New(dsn)
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}

	err = db.Conn.AutoMigrate(&entity.User{})
	if err != nil {
		log.Fatalf("Migrations error: %v", err)
	}

	userRepo := repo.NewUserRepo(db)
	userUseCase := usecase.NewUserUseCase(userRepo)

	handler := gin.Default()

	v1.NewUserRoutes(handler.Group("/v1"), userUseCase, logger.Default)

	handler.Run(":8080")
}
