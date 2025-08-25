package repo

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type RefreshPGX struct{ pool *pgxpool.Pool }

func NewRefreshPGX(pool *pgxpool.Pool) *RefreshPGX { return &RefreshPGX{pool: pool} }

func (r *RefreshPGX) Insert(ctx context.Context, userID, tokenHash string, expiresAt time.Time) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO refresh_tokens (id, user_id, token_hash, issued_at, expires_at)
		VALUES (gen_random_uuid()::text, $1, $2, NOW(), $3)
	`, userID, tokenHash, expiresAt)
	return err
}

func (r *RefreshPGX) FindValid(ctx context.Context, userID, tokenHash string, now time.Time) (bool, error) {
	var exists bool
	err := r.pool.QueryRow(ctx, `
		SELECT EXISTS (
			SELECT 1 FROM refresh_tokens
			WHERE user_id=$1
			  AND token_hash=$2
			  AND revoked_at IS NULL
			  AND expires_at > $3
		)
	`, userID, tokenHash, now).Scan(&exists)
	return exists, err
}

func (r *RefreshPGX) Revoke(ctx context.Context, userID, tokenHash string) error {
	_, err := r.pool.Exec(ctx, `
		UPDATE refresh_tokens
		   SET revoked_at = NOW()
		 WHERE user_id=$1 AND token_hash=$2 AND revoked_at IS NULL
	`, userID, tokenHash)
	return err
}

func (r *RefreshPGX) Purge(ctx context.Context, expiresBefore time.Time, revokedBefore time.Time) (int64, error) {
	expiresBefore = expiresBefore.UTC()
	revokedBefore = revokedBefore.UTC()

	tag, err := r.pool.Exec(ctx, `
		DELETE FROM refresh_tokens
		WHERE
			expires_at < $1
			OR (revoked_at IS NOT NULL AND revoked_at < $2)
	`, expiresBefore, revokedBefore)
	if err != nil {
		return 0, err
	}
	return tag.RowsAffected(), nil
}
