package models

import (
	"time"

	"github.com/google/uuid"
)

type File struct {
	ID           string     `json:"id"`
	OwnerUserID  string     `json:"owner_user_id"`
	ObjectKey    string     `json:"object_key"`
	OriginalName string     `json:"original_name"`
	SizeBytes    int64      `json:"size_bytes"`
	ContentType  string     `json:"content_type,omitempty"`
	SHA256       string     `json:"sha256,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
	DeletedAt    *time.Time `json:"deleted_at,omitempty"`
}

type FileMeta struct {
	ID          string    `json:"id" example:"file_123"`
	UserID      string    `json:"user_id" example:"User_65b80522-50be-4012-9964-550369cdcff7"`
	Name        string    `json:"name" example:"report.pdf"`
	Size        int64     `json:"size" example:"102400"`
	ContentType string    `json:"content_type" example:"application/pdf"`
	CreatedAt   time.Time `json:"created_at" example:"2025-08-26T05:53:20Z"`
}

func GenerateFileID() string {
	return "File_" + uuid.New().String()
}
