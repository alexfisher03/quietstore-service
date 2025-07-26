package handlers

import (
	"github.com/gofiber/fiber/v2"
)

func GetAllFilesForUserHandler(c *fiber.Ctx) error {
	// eventually call a service to get files for user
	return c.JSON(fiber.Map{
		"message": "All files",
		"files": []fiber.Map{
			{"id": "1", "filename": "test.pdf"},
			{"id": "2", "filename": "image.png"},
		},
	})
}

func UploadFileHandler(c *fiber.Ctx) error {
	// eventually handle file upload
	return c.JSON(fiber.Map{
		"message": "File uploaded successfully",
		"id":      "generated-file-id",
	})
}

func GetFileByIdHandler(c *fiber.Ctx) error {
	fileId := c.Params("id")
	// eventually call a service to get file by ID
	return c.JSON(fiber.Map{
		"message":  "File details",
		"id":       fileId,
		"filename": "example.pdf",
	})
}
