package service

import (
	"context"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/alexfisher03/quietstore-service/QuietStore/internal/models"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinIOConfig struct {
	Endpoint        string
	AccessKeyID     string
	SecretAccessKey string
	BucketName      string
	UseSSL          bool
	Region          string
}

type MinIOStorageService struct {
	client     *minio.Client
	bucketName string
	region     string
}

func NewMinIOStorageService(config MinIOConfig) (*MinIOStorageService, error) {
	client, err := minio.New(config.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(config.AccessKeyID, config.SecretAccessKey, ""),
		Secure: config.UseSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("Failed during MinIO client initialization: %w", err)
	}

	service := &MinIOStorageService{
		client:     client,
		bucketName: config.BucketName,
		region:     config.Region,
	}

	ctx := context.Background()
	if err := service.ensureBucket(ctx); err != nil {
		return nil, fmt.Errorf("Failed to ensure buncket: %w", err)
	}

	if err := service.validateStorageConfig(); err != nil {
		return nil, fmt.Errorf("Failed to validate storage config: %w", err)
	}
	return service, nil
}

func (s *MinIOStorageService) Store(ctx context.Context, file *models.File, content io.Reader) error {
	if err := file.Validate(); err != nil {
		return fmt.Errorf("File validation failed you fucker: %w", err)
	}

	file.GenerateStoragePath()

	opts := minio.PutObjectOptions{
		ContentType: file.ContentType,
		UserMetadata: map[string]string{
			"uploaded-by":   file.UserID,
			"original-name": file.Filename,
		},
	}

	info, err := s.client.PutObject(ctx, s.bucketName, file.StoragePath, content, file.Size, opts)
	if err != nil {
		if minio.ToErrorResponse(err).Code == "NoSuchBucket" {
			if createErr := s.ensureBucket(ctx); createErr != nil {
				return fmt.Errorf("Failed to create bucket: %w", createErr)
			}
			info, err = s.client.PutObject(ctx, s.bucketName, file.StoragePath, content, file.Size, opts)
			if err != nil {
				return fmt.Errorf("upload retry failed : %w", err)
			}
		} else {
			return fmt.Errorf("Failed to upload file: %w", err)
		}
	}

	if info.Size != file.Size {
		return fmt.Errorf("uploaded file size mismatch: expected %d, got %d", file.Size, info.Size)
	}
	return nil
}

func (s *MinIOStorageService) Retrieve(ctx context.Context, file *models.File) (io.ReadCloser, error) {
	if file.StoragePath == "" {
		return nil, ErrFileNotFound
	}

	object, err := s.client.GetObject(ctx, s.bucketName, file.StoragePath, minio.GetObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("Failed to retrieve file: %w", err)
	}

	_, err = object.Stat()
	if err != nil {
		object.Close()
		errResp := minio.ToErrorResponse(err)
		if errResp.Code == "NoSuchKey" {
			return nil, ErrFileNotFound
		}
		return nil, fmt.Errorf("Failed to stat object: %w", err)
	}

	return object, nil
}

func (s *MinIOStorageService) Delete(ctx context.Context, file *models.File) error {
	if file.StoragePath == "" {
		return nil
	}

	err := s.client.RemoveObject(ctx, s.bucketName, file.StoragePath, minio.RemoveObjectOptions{})
	if err != nil {
		if minio.ToErrorResponse(err).Code == "NoSuchKey" {
			return nil
		}
		return fmt.Errorf("Failed to delete file: %w", err)
	}
	return nil
}

func (s *MinIOStorageService) GetPresignedURL(ctx context.Context, file *models.File, expiry time.Duration) (string, error) {
	if file.StoragePath == "" {
		return "", errors.New("file storage path is empty during presigned URL generation")
	}

	url, err := s.client.PresignedGetObject(ctx, s.bucketName, file.StoragePath, expiry, nil)
	if err != nil {
		return "", fmt.Errorf("Failed to generate presigned URL: %w", err)
	}
	return url.String(), nil
}

func (s *MinIOStorageService) ensureBucket(ctx context.Context) error {
	exists, err := s.client.BucketExists(ctx, s.bucketName)
	if err != nil {
		return fmt.Errorf("Failed to check if bucket exists: %w", err)
	}

	if !exists {
		err = s.client.MakeBucket(ctx, s.bucketName, minio.MakeBucketOptions{Region: s.region})
		if err != nil {
			return fmt.Errorf("Failed to create bucket: %w", err)
		}
	}

	return nil
}

func (c *MinIOStorageService) validateStorageConfig() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := c.client.ListBuckets(ctx)
	if err != nil {
		return fmt.Errorf("MinIO connection failed: %w", err)
	}
	return nil
}
