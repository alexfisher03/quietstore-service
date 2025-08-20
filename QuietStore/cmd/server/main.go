package main

import (
	"fmt"
	"log"

	"github.com/alexfisher03/quietstore-service/QuietStore/internal/config"
	"github.com/alexfisher03/quietstore-service/QuietStore/internal/service"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"

	v1 "github.com/alexfisher03/quietstore-service/QuietStore/api/v1"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Println("No .env file loaded (continuing anyway)")
	}
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	storage := service.NewLocalStorageService(cfg.Storage.BasePath)

	app := fiber.New(fiber.Config{
		ServerHeader: "QuietStore/1.0",
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		BodyLimit:    cfg.Server.BodyLimit,
	})

	app.Use(logger.New())
	app.Use(recover.New())

	v1.RegisterRoutes(app, storage)

	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	log.Printf("[BOOT] Starting server at %s", addr)
	log.Printf("[BOOT] ENVIRONMENT: %s", cfg.App.Environment)
	log.Fatal(app.Listen(addr))
}
