package v1

import (
	"github.com/alexfisher03/quietstore-service/QuietStore/internal/handlers"
	"github.com/alexfisher03/quietstore-service/QuietStore/internal/service"
	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(app *fiber.App, storage service.StorageService) {
	v1 := app.Group("/api/v1")
	v1.Get("/health", handlers.HealthCheck)

	fileHandlers := handlers.NewFileHandler(storage)
	v1.Post("/files/upload", fileHandlers.UploadFileHandler)
	v1.Get("/:userID/files", fileHandlers.GetUserFilesHandler)
	v1.Get("/:userID/files/:fileID", fileHandlers.GetUserFileByIDHandler)
	v1.Delete("/:userID/files/:fileID", fileHandlers.DeleteUserFileByIDHandler)
}
