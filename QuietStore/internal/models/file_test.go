package models

import (
	"testing"
)

func TestNewFile(t *testing.T) {
	file := NewFile("test.pdf", 1024, "application/pdf", "user123")

	if file.Filename != "test.pdf" {
		t.Errorf("Expected filename 'test.pdf', got '%s'", file.Filename)
	}

	if file.Size != 1024 {
		t.Errorf("Expected size 1024, got %d", file.Size)
	}

	if file.UserID != "user123" {
		t.Errorf("Expected UserID 'user123', got '%s'", file.UserID)
	}
	if file.ID == "" {
		t.Error("Expected ID to be generated")
	}

	if file.UploadedAt.IsZero() {
		t.Error("Expected UploadedAt to be set")
	}
}

func TestFileValidate(t *testing.T) {
	tests := []struct {
		name    string
		file    *File
		wantErr bool
	}{
		{
			name: "valid file",
			file: &File{
				Filename:    "test.pdf",
				Size:        1024,
				ContentType: "application/pdf",
				UserID:      "user123",
			},
			wantErr: false,
		},
		{
			name: "empty filename",
			file: &File{
				Filename:    "",
				Size:        1024,
				ContentType: "application/pdf",
				UserID:      "user123",
			},
			wantErr: true,
		},
		{
			name: "zero size",
			file: &File{
				Filename:    "test.pdf",
				Size:        0,
				ContentType: "application/pdf",
				UserID:      "user123",
			},
			wantErr: true,
		},
		{
			name: "empty content type",
			file: &File{
				Filename:    "test.pdf",
				Size:        1024,
				ContentType: "",
				UserID:      "user123",
			},
			wantErr: true,
		},
		{
			name: "disallowed content type",
			file: &File{
				Filename:    "test.exe",
				Size:        1024,
				ContentType: "application/x-executable",
				UserID:      "user123",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.file.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGenerateStoragePath(t *testing.T) {
	file := &File{
		ID:       "file123",
		Filename: "document.pdf",
		UserID:   "user456",
	}

	file.GenerateStoragePath()

	if file.StoragePath == "" {
		t.Error("Expected storage path to be generated")
	}

	// Check that path contains expected components
	expectedParts := []string{"files", "user456", "file123-document.pdf"}
	for _, part := range expectedParts {
		if !contains(file.StoragePath, part) {
			t.Errorf("Expected storage path to contain '%s', got '%s'", part, file.StoragePath)
		}
	}
}

func TestSanitizeFilename(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"normal.pdf", "normal.pdf"},
		{"file with spaces.pdf", "file_with_spaces.pdf"},
		{"../../../etc/passwd", "passwd"},
		{"file~with`special|chars.pdf", "filewithspecialchars.pdf"},
		{"file$.txt", "file.txt"},
	}

	for _, tt := range tests {
		result := sanitizeFilename(tt.input)
		if result != tt.expected {
			t.Errorf("sanitizeFilename(%s) = %s, want %s", tt.input, result, tt.expected)
		}
	}
}

func TestIsImage(t *testing.T) {
	tests := []struct {
		contentType string
		expected    bool
	}{
		{"image/jpeg", true},
		{"image/png", true},
		{"image/gif", true},
		{"application/pdf", false},
		{"text/plain", false},
	}

	for _, tt := range tests {
		file := &File{ContentType: tt.contentType}
		if got := file.IsImage(); got != tt.expected {
			t.Errorf("IsImage() for %s = %v, want %v", tt.contentType, got, tt.expected)
		}
	}
}

// Helper function for testing
func contains(s, substr string) bool {
	return len(s) >= len(substr) && s[len(s)-len(substr):] == substr ||
		len(s) >= len(substr) && s[:len(substr)] == substr ||
		len(s) > len(substr) && containsMiddle(s, substr)
}

func containsMiddle(s, substr string) bool {
	for i := 1; i < len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
