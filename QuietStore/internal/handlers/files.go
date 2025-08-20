package handlers

import (
	"fmt"

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

func FindFileByID(files []*models.File, fileID string) *models.File {
	for _, file := range files {
		if file.ID == fileID {
			return file
		}
	}
	return nil
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
	print(newFile.StoragePath)

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

func (h *FileHandler) GetUserFilesHandler(c *fiber.Ctx) error {
	userID := c.Params("userID")
	if userID == "" {
		return c.Status(400).JSON(fiber.Map{"error": "User ID is required but empty"})
	}

	files, err := h.storage.ListUserFiles(c.Context(), userID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to list user files during storage lookup"})
	}

	var responses []models.FileResponse
	for _, file := range files {
		responses = append(responses, *file.ToResponse())
	}

	return c.JSON(responses)
}

func (h *FileHandler) GetUserFileByIDHandler(c *fiber.Ctx) error {
	userID := c.Params("userID")
	fileID := c.Params("fileID")

	if userID == "" || fileID == "" {
		return c.Status(400).JSON(fiber.Map{"error": "User ID or File ID are missing"})
	}

	files, err := h.storage.ListUserFiles(c.Context(), userID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to list user files"})
	}

	var match *models.File
	match = FindFileByID(files, fileID)
	if match == nil {
		return c.Status(404).JSON(fiber.Map{"error": "File not found"})
	}

	reader, err := h.storage.Retrieve(c.Context(), match)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to retrieve file"})
	}

	c.Set("Content-Type", match.ContentType)
	c.Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, match.Filename))
	return c.SendStream(reader)
}

func (h *FileHandler) DeleteUserFileByIDHandler(c *fiber.Ctx) error {
	userID := c.Params("userID")
	fileID := c.Params("fileID")

	if userID == "" || fileID == "" {
		return c.Status(400).JSON(fiber.Map{"error": "User ID or File ID are missing"})
	}

	files, err := h.storage.ListUserFiles(c.Context(), userID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to list user files"})
	}

	var targetFile *models.File
	if targetFile = FindFileByID(files, fileID); targetFile == nil {
		return c.Status(404).JSON(fiber.Map{"error": "File not found"})
	}

	if err := h.storage.DeleteFile(c.Context(), targetFile); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to delete file"})
	}
	return c.Status(200).JSON(fiber.Map{"message": "File deleted successfully"})
}
