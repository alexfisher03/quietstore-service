package handlers

import (
	"fmt"
	"strconv"

	"github.com/alexfisher03/quietstore-service/QuietStore/internal/models"
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

// UploadFileHandler godoc
//
//	@Summary		Upload a file
//	@Description	Uploads a file for the authenticated user
//	@Tags			files
//	@Security		BearerAuth
//	@Accept			multipart/form-data
//	@Param			file	formData	file	true	"file"
//	@Produce		json
//	@Success		200			{object}	models.FileMeta
//	@Failure		400,401,500	{object}	map[string]string
//	@Router			/files [post]
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

// GetUserFilesHandler godoc
//
//	@Summary		List my files
//	@Description	Returns all files belonging to the authenticated user
//	@Tags			files
//	@Security		BearerAuth
//	@Produce		json
//	@Success		200	{array}		models.FileMeta
//	@Failure		401	{object}	map[string]string
//	@Router			/files [get]
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

// GetFileByIdHandler godoc
//
//	@Summary		Get file by ID
//	@Description	Retrieves metadata for a specific file
//	@Tags			files
//	@Security		BearerAuth
//	@Produce		json
//	@Param			id	path		string	true	"File ID"
//	@Success		200	{object}	models.FileMeta
//	@Failure		401	{object}	map[string]string
//	@Failure		404	{object}	map[string]string
//	@Router			/files/{id} [get]
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

// DeleteUserFileByIDHandler godoc
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

func (h *FileHandler) SearchFilesHandler(c *fiber.Ctx) error {
	userID, err := resolveUserID(c)
	if err != nil {
		return err
	}

	q := c.Query("q", "")
	ctype := c.Query("type", "")
	minSize, _ := strconv.ParseInt(c.Query("min_size", "0"), 10, 64)
	maxSize, _ := strconv.ParseInt(c.Query("max_size", "0"), 10, 64)
	limit, _ := strconv.Atoi(c.Query("limit", "50"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))

	files, err := h.storage.SearchFiles(c.Context(), userID, q, ctype, minSize, maxSize, limit, offset)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "search failed: "+err.Error())
	}

	if files == nil {
		files = []*models.File{}
	}

	return c.JSON(files)
}

type RenameFileRequest struct {
	NewName string `json:"new_name"`
}

func (h *FileHandler) RenameFileHandler(c *fiber.Ctx) error {
	userID, err := resolveUserID(c)
	if err != nil {
		return err
	}

	fileID := c.Params("fileID")
	var req RenameFileRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}
	if req.NewName == "" {
		return fiber.NewError(fiber.StatusBadRequest, "new_name is required")
	}

	if err := h.storage.RenameFile(c.Context(), userID, fileID, req.NewName); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "rename failed: "+err.Error())
	}

	return c.JSON(fiber.Map{"File name": req.NewName})
}
