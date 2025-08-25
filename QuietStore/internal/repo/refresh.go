package repo

import (
	"context"
	"time"
)

type RefreshTokens interface {
	Insert(ctx context.Context, userID string, tokenHash string, expiresAt time.Time) error
	FindValid(ctx context.Context, userID string, tokenHash string, now time.Time) (bool, error)
	Revoke(ctx context.Context, userID string, tokenHash string) error
	Purge(ctx context.Context, expiresBefore time.Time, revokedBefore time.Time) (int64, error)
}
