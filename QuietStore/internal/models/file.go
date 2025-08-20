package models

import (
	"errors"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
)

type File struct {
	ID          string    `json:"id"`
	Filename    string    `json:"filename"`
	Size        int64     `json:"size"`
	ContentType string    `json:"content_type"`
	StoragePath string    `json:"-"`
	UserID      string    `json:"user_id"`
	UploadedAt  time.Time `json:"uploaded_at"`
}

func generateID() string {
	return "UserFile_" + uuid.New().String()
}

func sanitizeFilename(filename string) string {
	filename = filepath.Base(filename)

	filename = strings.ReplaceAll(filename, " ", "_")

	dangerous := []string{"..", "~", "`", "|", ";", "&", "$", "*"}

	for _, char := range dangerous {
		filename = strings.ReplaceAll(filename, char, "des")
	}

	return filename
}

func NewFile(filename string, size int64, contentType, userID string) *File {
	return &File{
		ID:          generateID(),
		Filename:    sanitizeFilename(filename),
		Size:        size,
		ContentType: contentType,
		UserID:      userID,
		UploadedAt:  time.Now(),
	}
}

func (f *File) Validate() error {
	if f.Filename == "" {
		return errors.New("filename cannot be empty")
	}
	if f.Size <= 0 {
		return errors.New("file size must be greater than 0")
	}
	if f.ContentType == "" {
		return errors.New("content type cannot be empty")
	}
	if f.UserID == "" {
		return errors.New("user ID cannot be empty")
	}
	if !isAllowedContentType(f.ContentType) {
		return errors.New("file type not allowed")
	}
	return nil
}

func (f *File) GenerateStoragePath() {
	now := time.Now()
	f.StoragePath = filepath.Join(
		"files",
		f.UserID,
		now.Format("2006"),
		now.Format("01"),
		f.ID+"-"+f.Filename,
	)
}

func (f *File) GetExtension() string {
	return strings.ToLower(filepath.Ext(f.Filename))
}

func (f *File) IsImage() bool {
	switch f.ContentType {
	case "image/jpeg", "image/jpg", "image/png", "image/gif", "image/webp":
		return true
	}
	return false
}

type FileCreateRequest struct {
	Filename    string `json:"filename" validate:"required"`
	Size        int64  `json:"size" validate:"required,min=1"`
	ContentType string `json:"content_type" validate:"required"`
}

// represents file data returned to client side
type FileResponse struct {
	ID          string    `json:"id"`
	Filename    string    `json:"filename"`
	Size        int64     `json:"size"`
	ContentType string    `json:"content_type"`
	UploadURL   string    `json:"upload_url,omitempty"`
	DownloadURL string    `json:"download_url,omitempty"`
	UploadedAt  time.Time `json:"uploaded_at"`
}

func (f *File) ToResponse() *FileResponse {
	return &FileResponse{
		ID:          f.ID,
		Filename:    f.Filename,
		Size:        f.Size,
		ContentType: f.ContentType,
		UploadedAt:  f.UploadedAt,
	}
}

func isAllowedContentType(contentType string) bool {
	allowed := []string{
		"image/jpeg",
		"image/jpg",
		"image/png",
		"image/gif",
		"image/webp",
		"application/pdf",
		"application/json",
		"application/zip",
		"application/docx",
		"application/xlsx",
		"text/plain",
		"text/html",
		"text/css",
		"text/javascript",
		"text/typescript",
		"text/x-go",
		"text/x-python",
		"text/x-c++src",
	}

	for _, ct := range allowed {
		if ct == contentType {
			return true
		}
	}
	return false
}

const (
	MaxFileSize = 1000 * 1024 * 1024
	MinFileSize = 1
)
