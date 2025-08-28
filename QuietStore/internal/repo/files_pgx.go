package repo

import (
	"context"
	"errors"
	"fmt"

	"github.com/alexfisher03/quietstore-service/QuietStore/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type FilesPGX struct{ pool *pgxpool.Pool }

func NewFilesPGX(pool *pgxpool.Pool) *FilesPGX { return &FilesPGX{pool: pool} }

func (r *FilesPGX) Create(ctx context.Context, f *models.File) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO files (id, owner_user_id, object_key, original_name, size_bytes, content_type, sha256, created_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`,
		f.ID, f.OwnerUserID, f.ObjectKey, f.OriginalName, f.SizeBytes, f.ContentType, f.SHA256, f.CreatedAt)
	return err
}

func (r *FilesPGX) ByID(ctx context.Context, id string) (*models.File, error) {
	row := r.pool.QueryRow(ctx, `
		SELECT id, owner_user_id, object_key, original_name, size_bytes, content_type, sha256, created_at, deleted_at
		FROM files WHERE id=$1`, id)
	var f models.File
	if err := row.Scan(&f.ID, &f.OwnerUserID, &f.ObjectKey, &f.OriginalName, &f.SizeBytes, &f.ContentType, &f.SHA256, &f.CreatedAt, &f.DeletedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &f, nil
}

func (r *FilesPGX) ListByOwner(ctx context.Context, ownerID string, limit, offset int) ([]*models.File, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, owner_user_id, object_key, original_name, size_bytes, content_type, sha256, created_at, deleted_at
		FROM files
		WHERE owner_user_id=$1 AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3`, ownerID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []*models.File
	for rows.Next() {
		var f models.File
		if err := rows.Scan(&f.ID, &f.OwnerUserID, &f.ObjectKey, &f.OriginalName, &f.SizeBytes, &f.ContentType, &f.SHA256, &f.CreatedAt, &f.DeletedAt); err != nil {
			return nil, err
		}
		out = append(out, &f)
	}
	return out, rows.Err()
}

func (r *FilesPGX) Delete(ctx context.Context, id string, ownerID string) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM files WHERE id=$1 AND owner_user_id=$2`, id, ownerID)
	return err
}

func (r *FilesPGX) ListByFilters(
	ctx context.Context,
	userID, q, contentType string,
	minSize, maxSize int64,
	limit, offset int,
) ([]*models.File, error) {

	rows, err := r.pool.Query(ctx, `
        SELECT id, owner_user_id, object_key, original_name, size_bytes,
               content_type, sha256, created_at, deleted_at
        FROM files
        WHERE owner_user_id = $1
		  AND deleted_at IS NULL
          AND ($2 = '' OR original_name ILIKE '%' || $2 || '%')
          AND ($3 = '' OR content_type = $3)
          AND ($4 <= 0 OR size_bytes >= $4)
          AND ($5 <= 0 OR size_bytes <= $5)
        ORDER BY created_at DESC
        LIMIT $6 OFFSET $7;
    `, userID, q, contentType, minSize, maxSize, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}

	out := make([]*models.File, 0)
	for rows.Next() {
		var f models.File
		if err := rows.Scan(
			&f.ID, &f.OwnerUserID, &f.ObjectKey, &f.OriginalName,
			&f.SizeBytes, &f.ContentType, &f.SHA256, &f.CreatedAt, &f.DeletedAt,
		); err != nil {
			return nil, fmt.Errorf("scan failed: %w", err)
		}
		out = append(out, &f)
	}
	rows.Close()
	return out, rows.Err()
}

func (r *FilesPGX) UpdateOriginalName(ctx context.Context, fileID, userID, newName string) error {
	_, err := r.pool.Exec(ctx, `
		UPDATE files SET original_name = $1
		WHERE id = $2 AND owner_user_id = $3 AND deleted_at IS NULL`,
		newName, fileID, userID)
	return err
}
