// internal/repo/users_pgx.go
package repo

import (
	"context"
	"errors"

	"github.com/alexfisher03/quietstore-service/QuietStore/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UsersPGX struct{ pool *pgxpool.Pool }

func NewUsersPGX(pool *pgxpool.Pool) *UsersPGX { return &UsersPGX{pool: pool} }

func (r *UsersPGX) Create(ctx context.Context, u *models.User) error {
	_, err := r.pool.Exec(ctx, `
    INSERT INTO users (id, username, email, password_hash, created_at)
    VALUES ($1, $2, $3, $4, $5)`,
		u.ID, u.Username, u.Email, u.Password, u.CreatedAt)
	return err
}

func (r *UsersPGX) ByID(ctx context.Context, id string) (*models.User, error) {
	row := r.pool.QueryRow(ctx, `
    SELECT id, username, email, password_hash, created_at
    FROM users WHERE id=$1`, id)
	var u models.User
	if err := row.Scan(&u.ID, &u.Username, &u.Email, &u.Password, &u.CreatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
}

func (r *UsersPGX) ByUsername(ctx context.Context, username string) (*models.User, error) {
	row := r.pool.QueryRow(ctx, `
    SELECT id, username, email, password_hash, created_at
    FROM users WHERE username=$1`, username)
	var u models.User
	if err := row.Scan(&u.ID, &u.Username, &u.Email, &u.Password, &u.CreatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
}

func (r *UsersPGX) Update(ctx context.Context, u *models.User) error {
	_, err := r.pool.Exec(ctx, `
		UPDATE users SET
			username = $1,
			email = $2,
			password_hash = $3
		WHERE id = $4`,
		u.Username, u.Email, u.Password, u.ID)
	return err
}

func (r *UsersPGX) Delete(ctx context.Context, id string) error {
	_, err := r.pool.Exec(ctx, `
		DELETE FROM users WHERE id = $1`, id)
	return err
}

func (r *UsersPGX) List(ctx context.Context, limit, offset int) ([]*models.User, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, username, email, password_hash, created_at
		FROM users
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2`, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		var u models.User
		if err := rows.Scan(&u.ID, &u.Username, &u.Email, &u.Password, &u.CreatedAt); err != nil {
			return nil, err
		}
		users = append(users, &u)
	}
	return users, nil
}
