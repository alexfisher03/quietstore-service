package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"

	v1 "github.com/alexfisher03/quietstore-service/QuietStore/api/v1"
	"github.com/alexfisher03/quietstore-service/QuietStore/internal/handlers"
)

func main() {
	app := fiber.New(fiber.Config{
		ServerHeader: "QuietStore/1.0",
		ErrorHandler: handlers.CustomErrorHandler,
		BodyLimit:    1024 * 1024 * 4,
	})

	app.Use(logger.New())
	app.Use(recover.New())

	v1.RegisterRoutes(app)

	log.Println("\033[1;36m")
	log.Println("▄▖▖▖▄▖▄▖▄▖▄▖▄▖▄▖▄▖▄▖")
	log.Println("▌▌▌▌▐ ▙▖▐ ▚ ▐ ▌▌▙▘▙▖")
	log.Println("█▌▙▌▟▖▙▖▐ ▄▌▐ ▙▌▌▌▙")
	log.Println(" ▘                  ")
	log.Println("▄▖▄▖▄▖▖▖▄▖▄▖▄▖▄▖")
	log.Println("▚ ▙▖▙▘▌▌▐ ▌ ▙▖▚")
	log.Println("▄▌▙▖▌▌▚▘▟▖▙▖▙▖▄▌")
	log.Println("\033[0m")
	log.Println("\033[1;32m==========================================\033[0m")
	log.Println("\033[1;33m  Server \033[1;35mSTARTING\033[0m \033[1;33mon \033[1;36m:8080\033[0m")
	log.Println("\033[1;32m==========================================\033[0m")
	log.Fatal(app.Listen(":8080"))
}
