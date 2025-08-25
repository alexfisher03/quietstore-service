package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/alexfisher03/quietstore-service/QuietStore/internal/models"
	"github.com/alexfisher03/quietstore-service/QuietStore/internal/repo"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type MinIOStorageService struct {
	s3     *s3.Client
	bucket string
	files  repo.Files
}

func NewMinIOStorageService(s3c *s3.Client, bucket string, files repo.Files) *MinIOStorageService {
	return &MinIOStorageService{s3: s3c, bucket: bucket, files: files}
}

func (m *MinIOStorageService) objectKey(userID, fileID string, t time.Time) string {
	return fmt.Sprintf("user/%s/%04d/%02d/%s", userID, t.Year(), int(t.Month()), fileID)
}

func (m *MinIOStorageService) SaveFile(
	ctx context.Context,
	userID, originalName, contentType string,
	size int64,
	r io.Reader,
) (*models.File, error) {

	now := time.Now()
	id := models.GenerateFileID()
	key := m.objectKey(userID, id, now)

	tmp, err := os.CreateTemp("", "qs-upload-*")
	if err != nil {
		return nil, err
	}
	defer func() {
		tmp.Close()
		_ = os.Remove(tmp.Name())
	}()

	n, err := io.Copy(tmp, r)
	if err != nil {
		return nil, fmt.Errorf("buffer upload: %w", err)
	}

	if _, err := tmp.Seek(0, io.SeekStart); err != nil {
		return nil, err
	}
	h := sha256.New()
	if _, err := io.Copy(h, tmp); err != nil {
		return nil, err
	}
	sum := hex.EncodeToString(h.Sum(nil))

	if _, err := tmp.Seek(0, io.SeekStart); err != nil {
		return nil, err
	}

	_, err = m.s3.PutObject(ctx, &s3.PutObjectInput{
		Bucket:        aws.String(m.bucket),
		Key:           aws.String(key),
		Body:          tmp,
		ContentType:   aws.String(contentType),
		ContentLength: aws.Int64(n),
	})
	if err != nil {
		return nil, err
	}

	f := &models.File{
		ID:           id,
		OwnerUserID:  userID,
		ObjectKey:    key,
		OriginalName: originalName,
		SizeBytes:    n,
		ContentType:  contentType,
		SHA256:       sum,
		CreatedAt:    now,
	}
	if err := m.files.Create(ctx, f); err != nil {
		_, _ = m.s3.DeleteObject(ctx, &s3.DeleteObjectInput{
			Bucket: aws.String(m.bucket),
			Key:    aws.String(key),
		})
		return nil, err
	}
	return f, nil
}

func (m *MinIOStorageService) OpenFile(ctx context.Context, userID, fileID string) (*models.File, io.ReadCloser, error) {
	meta, err := m.files.ByID(ctx, fileID)
	if err != nil || meta == nil {
		return nil, nil, err
	}
	if meta.OwnerUserID != userID {
		return nil, nil, fmt.Errorf("not found")
	}

	obj, err := m.s3.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(m.bucket),
		Key:    aws.String(meta.ObjectKey),
	})
	if err != nil {
		return nil, nil, err
	}
	return meta, obj.Body, nil
}

func (m *MinIOStorageService) ListFiles(ctx context.Context, userID string, limit, offset int) ([]*models.File, error) {
	return m.files.ListByOwner(ctx, userID, limit, offset)
}

func (m *MinIOStorageService) DeleteFile(ctx context.Context, userID, fileID string) error {
	meta, err := m.files.ByID(ctx, fileID)
	if err != nil || meta == nil {
		return err
	}
	if meta.OwnerUserID != userID {
		return fmt.Errorf("not found")
	}

	if _, err := m.s3.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(m.bucket), Key: aws.String(meta.ObjectKey),
	}); err != nil {
		return err
	}
	return m.files.Delete(ctx, fileID, userID)
}
