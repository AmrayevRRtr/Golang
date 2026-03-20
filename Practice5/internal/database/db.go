package database

import (
	"Practice5/internal/models"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

//func NewMySQL() *sqlx.DB {
//
//	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
//		os.Getenv("DB_USER"),
//		os.Getenv("DB_PASSWORD"),
//		os.Getenv("DB_HOST"),
//		os.Getenv("DB_PORT"),
//		os.Getenv("DB_NAME"),
//	)
//
//	db, err := sqlx.Open("mysql", dsn)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	err = db.Ping()
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	log.Println("Connected to MySQL")
//	return db
//}

var DB *gorm.DB

func InitDB() {
	var err error
	err = godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Connected Succesfully to Database")
	db.Logger = logger.Default.LogMode(logger.Info)
	log.Println("Running Migrations")
	db.AutoMigrate(&models.User{})
	DB = db
}
