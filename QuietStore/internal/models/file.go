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

func GenerateFileID() string {
	return "File_" + uuid.New().String()
}
