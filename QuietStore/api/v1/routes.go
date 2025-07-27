package v1

import (
	"github.com/alexfisher03/quietstore-service/QuietStore/internal/handlers"
	"github.com/alexfisher03/quietstore-service/QuietStore/internal/service"
	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(app *fiber.App, storageService service.StorageService) {

	fileHandler := handlers.NewFileHandler(storageService)

	v1 := app.Group("/api/v1")
	v1.Get("/health", handlers.HealthCheck)
	v1.Post("/files/upload", fileHandler.UploadFileHandler)
	v1.Get("/files/:id", fileHandler.GetFileByIdHandler)
	v1.Get("/files", fileHandler.GetAllFilesForUserHandler)
}
