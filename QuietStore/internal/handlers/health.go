package handlers

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

// HealthCheck godoc
//
//	@Summary		Health
//	@Description	Liveness probe
//	@Tags			health
//	@Produce		json
//	@Success		200	{object}	map[string]interface{}
//	@Router			/health [get]
func HealthCheck(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"status":    "ok",
		"service":   "QuietStore",
		"timestamp": time.Now().Format(time.RFC3339),
	})
}
