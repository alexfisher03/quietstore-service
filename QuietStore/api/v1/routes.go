package v1

import (
	"github.com/alexfisher03/quietstore-service/QuietStore/internal/handlers"
	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(app *fiber.App) {
	v1 := app.Group("/api/v1")
	v1.Get("/health", handlers.HealthCheck)
}
