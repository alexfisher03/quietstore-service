package main

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"

	v1 "github.com/alexfisher03/quietstore-service/QuietStore/api/v1"
	"github.com/alexfisher03/quietstore-service/QuietStore/internal/config"
	"github.com/alexfisher03/quietstore-service/QuietStore/internal/handlers"
	"github.com/alexfisher03/quietstore-service/QuietStore/internal/service"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load configuration: %v", err)
	}

	log.Println("\033[1;32m==========================================\033[0m")
	log.Printf("\033[1;33m  Server \033[1;35mSTARTING\033[0m \033[1;33mon \033[1;36m:8080\033[0m")
	log.Printf("\033[1;36m  Running QuietStore in %s mode\033[0m", cfg.App.Environment)
	log.Println("\033[1;32m==========================================\033[0m")

	app := fiber.New(fiber.Config{
		ServerHeader: "QuietStore/1.0",
		ErrorHandler: handlers.CustomErrorHandler,
		BodyLimit:    cfg.Server.BodyLimit,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	})

	app.Use(logger.New())
	app.Use(recover.New())

	storageService, err := initializeStorage(cfg.Storage)
	if err != nil {
		log.Fatalf("Failed to initialize storage: %v", err)
	}

	v1.RegisterRoutes(app, storageService)

	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)

	log.Println("\033[1;36m")
	log.Println("▄▖▖▖▄▖▄▖▄▖▄▖▄▖▄▖▄▖▄▖")
	log.Println("▌▌▌▌▐ ▙▖▐ ▚ ▐ ▌▌▙▘▙▖")
	log.Println("█▌▙▌▟▖▙▖▐ ▄▌▐ ▙▌▌▌▙")
	log.Println(" ▘                  ")
	log.Println("▄▖▄▖▄▖▖▖▄▖▄▖▄▖▄▖")
	log.Println("▚ ▙▖▙▘▌▌▐ ▌ ▙▖▚")
	log.Println("▄▌▙▖▌▌▚▘▟▖▙▖▙▖▄▌")
	log.Println("======================")
	log.Println("Server running on " + addr)
	log.Println("\033[0m")

	log.Fatal(app.Listen(addr))
}

func initializeStorage(storageConfig config.StorageConfig) (service.StorageService, error) {
	minioConfig := service.MinIOConfig{
		Endpoint:        storageConfig.Endpoint,
		AccessKeyID:     storageConfig.AccessKeyID,
		SecretAccessKey: storageConfig.SecretAccessKey,
		BucketName:      storageConfig.BucketName,
		UseSSL:          storageConfig.UseSSL,
		Region:          storageConfig.Region,
	}

	return service.NewMinIOStorageService(minioConfig)
}
