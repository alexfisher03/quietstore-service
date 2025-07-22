package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"

	"github.com/alexfisher03/QuietStore/internal/handlers/health"
)

func main() {
	app := fiber.New(fiber.Config{
		ServerHeader: "QuietStore/1.0",
		ErrorHandler: customErrorHandler,
		BodyLimit:    1024 * 1024 * 4,
	})

	app.Use(logger.New())
	app.Use(recover.New())

	app.Get("/health", health.HealthCheck)
	app.Post("/echo", echoHandler)

	log.Fatal(app.Listen(":8080"))
}

func echoHandler(c *fiber.Ctx) error {
	var body map[string]interface{}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	return c.JSON(fiber.Map{
		"received": body,
		"method":   c.Method(),
		"path":     c.Path(),
	})
}

func customErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError

	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
	}

	return c.Status(code).JSON(fiber.Map{
		"error": err.Error(),
		"code":  code,
	})
}
