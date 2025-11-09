package repositories

import (
	"bioly/auth/internal/types"
	"context"
	"errors"
	"time"

	"github.com/alexedwards/argon2id"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

var ErrInvalidCredentials = errors.New("invalid credentials")
var ErrNotImplemented = errors.New("not implemented")
var ErrDuplicateUsername = errors.New("username already exists")
var ErrNotFound = errors.New("not found")

type Users interface {
	Add(ctx context.Context, u *types.User) error
	Delete(ctx context.Context, id int64) error
	VerifyCredentials(ctx context.Context, username, password string) (*types.User, error)
}

type usersImpl struct {
	db *sqlx.DB
}

func NewUsers(db *sqlx.DB) *usersImpl {
	return &usersImpl{db: db}
}

func (r *usersImpl) Add(ctx context.Context, u *types.User) error {
	q := `
		INSERT INTO users (username, password_hash)
		VALUES ($1, $2)
		RETURNING id, created_at, updated_at
	`
	err := r.db.QueryRowxContext(ctx, q, u.Username, u.PasswordHash).Scan(&u.ID, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return ErrDuplicateUsername
		}
		return err
	}
	return nil
}

func (r *usersImpl) Delete(ctx context.Context, id int64) error {
	res, err := r.db.ExecContext(ctx, `DELETE FROM users WHERE id = $1`, id)
	if err != nil {
		return err
	}
	aff, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if aff == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *usersImpl) VerifyCredentials(ctx context.Context, username, password string) (*types.User, error) {
	var u types.User
	err := r.db.GetContext(ctx, &u, `
		SELECT id, username, password_hash, last_login_at, created_at, updated_at
		FROM users
		WHERE lower(username) = lower($1)
		LIMIT 1
	`, username)
	if err != nil {
		return nil, ErrInvalidCredentials
	}
	ok, err := argon2id.ComparePasswordAndHash(password, u.PasswordHash)
	if err != nil || !ok {
		return nil, ErrInvalidCredentials
	}
	_, _ = r.db.ExecContext(ctx, `
		UPDATE users
		SET last_login_at = NOW(), updated_at = NOW()
		WHERE id = $1
	`, u.ID)
	now := time.Now().UTC()
	u.LastLoginAt = &now
	return &u, nil
}
