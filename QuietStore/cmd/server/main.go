package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"

	v1 "github.com/alexfisher03/quietstore-service/QuietStore/api/v1"
)

func main() {
	app := fiber.New(fiber.Config{
		ServerHeader: "QuietStore/1.0",
	})

	app.Use(logger.New())
	app.Use(recover.New())

	v1.RegisterRoutes(app)

	log.Println("QuietStore backend running on http://localhost:8080")
	log.Fatal(app.Listen(":8080"))
}
