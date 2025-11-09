package repositories

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
)

type RefreshTokens interface {
	Create(ctx context.Context, userID int64, jti [16]byte, tokenHash, userAgent, ip string, expiresAt time.Time) error
	RevokeByJTI(ctx context.Context, jti [16]byte) error
	RevokeAllByUser(ctx context.Context, userID int64) error
	FindValidByJTI(ctx context.Context, jti [16]byte) (userID int64, tokenHash string, expiresAt time.Time, revoked bool, err error)
}

type refreshTokensImpl struct {
	db *sqlx.DB
}

func NewRefreshTokens(db *sqlx.DB) *refreshTokensImpl {
	return &refreshTokensImpl{db: db}
}

func (r *refreshTokensImpl) Create(ctx context.Context, userID int64, jti [16]byte, tokenHash, userAgent, ip string, expiresAt time.Time) error {
	return ErrNotImplemented
}

func (r *refreshTokensImpl) RevokeByJTI(ctx context.Context, jti [16]byte) error {
	return ErrNotImplemented
}

func (r *refreshTokensImpl) RevokeAllByUser(ctx context.Context, userID int64) error {
	return ErrNotImplemented
}

func (r *refreshTokensImpl) FindValidByJTI(ctx context.Context, jti [16]byte) (userID int64, tokenHash string, expiresAt time.Time, revoked bool, err error) {
	return 0, "", time.Time{}, false, ErrNotImplemented
}
