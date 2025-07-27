package handlers

import (
	"github.com/alexfisher03/quietstore-service/QuietStore/internal/models"
	"github.com/alexfisher03/quietstore-service/QuietStore/internal/service"
	"github.com/gofiber/fiber/v2"
)

type FileHandler struct {
	storage service.StorageService
}

func NewFileHandler(storage service.StorageService) *FileHandler {
	return &FileHandler{
		storage: storage,
	}
}

func (h *FileHandler) UploadFileHandler(c *fiber.Ctx) error {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "No file uploaded"})
	}

	file, err := fileHeader.Open()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to open file"})
	}
	defer file.Close()

	newFile := models.NewFile(fileHeader.Filename, fileHeader.Size, fileHeader.Header.Get("Content-Type"), "user123")
	if err := newFile.Validate(); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	newFile.GenerateStoragePath()

	if err := h.storage.Store(c.Context(), newFile, file); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "upload file failed: " + err.Error()})
	}
	return c.Status(200).JSON(fiber.Map{
		"message": "File uploaded successfully",
		"file": fiber.Map{
			"id":           newFile.ID,
			"filename":     newFile.Filename,
			"size":         newFile.Size,
			"content_type": newFile.ContentType,
			"uploaded_at":  newFile.UploadedAt,
		},
	})

}

func (h *FileHandler) GetFileByIdHandler(c *fiber.Ctx) error {
	fileID := c.Params("id")
	if fileID == "" {
		return c.Status(400).JSON(fiber.Map{"error": "File ID is required"})
	}

	// mock file creation, fetch from DB later
	file := models.NewFile("aperature01.pdf", 6000, "application/pdf", "user123")
	reader, err := h.storage.Retrieve(c.Context(), file)
	if err != nil {
		if err == service.ErrFileNotFound {
			return c.Status(404).JSON(fiber.Map{"error": "File not found"})
		}
		return c.Status(500).JSON(fiber.Map{"error": "Failed to retrieve file"})
	}
	defer reader.Close()
	if file == nil {
		return c.Status(404).JSON(fiber.Map{"error": "File not found"})
	}
	return c.Status(200).JSON(file.ToResponse())
}

func (h *FileHandler) GetAllFilesForUserHandler(c *fiber.Ctx) error {
	files := []fiber.Map{
		{
			"id":          "file1",
			"filename":    "proj.cpp",
			"size":        2048,
			"uploaded_at": "2023-10-01T12:00:00Z",
		},
		{
			"id":          "file2",
			"filename":    "report.pdf",
			"size":        10240,
			"uploaded_at": "2023-10-02T14:30:00Z",
		},
	}
	return c.JSON(fiber.Map{
		"files": files,
		"count": len(files),
	})
}
