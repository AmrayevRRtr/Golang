package main

import (
	"Practice5/internal/controllers"
	"Practice5/internal/database"
	"Practice5/internal/models"
	"fmt"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

// TIP <p>To run your code, right-click the code and select <b>Run</b>.</p> <p>Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.</p>
//func main() {
//	godotenv.Load()
//
//	db := database.NewMySQL()
//
//	repo := repository.NewRepository(db)
//
//	userHandler := handler.NewUserHandler(repo)
//
//	mux := http.NewServeMux()
//
//	userHandler.RegisterRoutes(mux)
//
//	log.Println("Server running on :8080")
//
//	http.ListenAndServe(":8080", mux)
//}

func main() {
	var err error
	err = godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	app := fiber.New()
	database.InitDB()
	app.Get("/users/seed", func(c *fiber.Ctx) error {
		var user models.User
		if err := database.DB.Exec("delete from users where 1").Error; err != nil {
			return c.SendStatus(500)
		}
		for i := 1; i <= 20; i++ {
			user.Name = fmt.Sprintf("Name %d", i)
			user.Email = fmt.Sprintf("email%d@gmail.com", i)
			user.Gender = fmt.Sprintf("male")
			user.BirthDate = time.Now().AddDate(-20-i, 0, 0)

			user.CreatedAt = time.Now().Add(-time.Duration(21-i) * time.Hour)

			database.DB.Create(&user)
		}
		return c.SendStatus(fiber.StatusOK)
	})
	app.Get("/users", controllers.GetPaginatedUsers)
	log.Fatal(app.Listen(":3000"))
}
