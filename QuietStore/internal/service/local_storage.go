package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/alexfisher03/quietstore-service/QuietStore/internal/models"
)

type LocalStorageService struct {
	BasePath string
}

func NewLocalStorageService(basePath string) *LocalStorageService {
	return &LocalStorageService{BasePath: basePath}
}

func (s *LocalStorageService) Store(ctx context.Context, file *models.File, content io.Reader) error {
	if file == nil || file.StoragePath == "" {
		return ErrInvalidFile
	}

	fullPath := filepath.Join(s.BasePath, file.StoragePath)

	if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	out, err := os.Create(fullPath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer out.Close()

	if _, err := io.Copy(out, content); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	metaPath := fullPath + ".meta.json"
	metaFile, err := os.Create(metaPath)
	if err != nil {
		return fmt.Errorf("failed to create metadata file: %w", err)
	}
	defer metaFile.Close()

	if err := json.NewEncoder(metaFile).Encode(file); err != nil {
		return fmt.Errorf("failed to encode meta data: %w", err)
	}

	return nil
}

func (s *LocalStorageService) Retrieve(ctx context.Context, file *models.File) (io.ReadCloser, error) {
	if file == nil || file.StoragePath == "" {
		return nil, ErrInvalidFile
	}

	fullPath := filepath.Join(s.BasePath, file.StoragePath)

	f, err := os.Open(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, ErrFileNotFound
		}
		return nil, fmt.Errorf("failed to open file: %w", err)
	}

	return f, nil
}

func (s *LocalStorageService) ListUserFiles(ctx context.Context, userID string) ([]*models.File, error) {
	root := filepath.Join(s.BasePath, "files", userID)
	var files []*models.File

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		if info.IsDir() || !strings.HasSuffix(path, ".meta.json") {
			return nil
		}

		metaFile, err := os.Open(path)
		if err != nil {
			log.Printf("[SKIP] Failed to open metadata file %s: %v", path, err)
			return nil
		}
		defer metaFile.Close()

		var f models.File
		f.StoragePath = strings.TrimSuffix(strings.TrimPrefix(path, s.BasePath+string(filepath.Separator)), ".meta.json")
		if err := json.NewDecoder(metaFile).Decode(&f); err != nil {
			log.Printf("[SKIP] Failed to decode meta file: %s", path)
			return nil
		}
		files = append(files, &f)
		return nil
	})

	if err != nil {
		return nil, err
	}

	return files, nil
}

func (s *LocalStorageService) DeleteFile(ctx context.Context, file *models.File) error {
	if file == nil || file.StoragePath == "" {
		return ErrInvalidFile
	}

	fullPath := filepath.Join(s.BasePath, file.StoragePath)
	metaPath := fullPath + ".meta.json"

	if err := os.Remove(fullPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete file: %w", err)
	}

	if err := os.Remove(metaPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete metadata: %w", err)
	}

	return nil
}
