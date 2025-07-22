package health

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

func HealthCheck(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"status":    "OK",
		"service":   "QuietStore",
		"timestamp": time.Now().Format(time.RFC3339),
	})
}
