package handlers

import (
	"fmt"
	"strconv"

	"github.com/alexfisher03/quietstore-service/QuietStore/internal/service"
	"github.com/gofiber/fiber/v2"
)

type FileHandler struct {
	storage service.StorageService
}

func NewFileHandler(storage service.StorageService) *FileHandler {
	return &FileHandler{storage: storage}
}

func resolveUserID(c *fiber.Ctx) (string, error) {
	if v := c.Locals("userID"); v != nil {
		if s, ok := v.(string); ok && s != "" {
			return s, nil
		}
	}
	return "", fiber.NewError(fiber.StatusUnauthorized, "missing user context")
}

func (h *FileHandler) UploadFileHandler(c *fiber.Ctx) error {
	userID, err := resolveUserID(c)
	if err != nil {
		return err
	}

	fh, err := c.FormFile("file")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "no file uploaded")
	}

	f, err := fh.Open()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "open file failed")
	}
	defer f.Close()

	ct := fh.Header.Get("Content-Type")
	meta, err := h.storage.SaveFile(c.Context(), userID, fh.Filename, ct, fh.Size, f)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "save failed: "+err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(meta)
}

func (h *FileHandler) GetUserFilesHandler(c *fiber.Ctx) error {
	userID, err := resolveUserID(c)
	if err != nil {
		return err
	}

	limit, _ := strconv.Atoi(c.Query("limit", "50"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))

	list, err := h.storage.ListFiles(c.Context(), userID, limit, offset)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "list failed: "+err.Error())
	}

	return c.JSON(list)
}

func (h *FileHandler) GetUserFileByIDHandler(c *fiber.Ctx) error {
	userID, err := resolveUserID(c)
	if err != nil {
		return err
	}

	fileID := c.Params("fileID")
	meta, rc, err := h.storage.OpenFile(c.Context(), userID, fileID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "open failed: "+err.Error())
	}
	defer rc.Close()

	if meta.ContentType != "" {
		c.Set("Content-Type", meta.ContentType)
	}
	c.Set("Content-Disposition", fmt.Sprintf(`attachment; filename=%q`, meta.OriginalName))
	return c.SendStream(rc)
}

func (h *FileHandler) DeleteUserFileByIDHandler(c *fiber.Ctx) error {
	userID, err := resolveUserID(c)
	if err != nil {
		return err
	}

	fileID := c.Params("fileID")
	if err := h.storage.DeleteFile(c.Context(), userID, fileID); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "delete failed: "+err.Error())
	}
	return c.JSON(fiber.Map{"message": "deleted"})
}
