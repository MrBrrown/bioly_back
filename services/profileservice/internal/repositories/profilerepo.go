package repositories

import (
	"bioly/profileservice/internal/types"
	"context"

	"github.com/jmoiron/sqlx"
)

type Profile interface {
	GetUserId(ctx context.Context, username string) (int64, error)
	GetProfile(ctx context.Context, id int64) (types.Profile, error)
}

type profilImpl struct {
	db *sqlx.DB
}

func NewProfile(db *sqlx.DB) Profile {
	return &profilImpl{db: db}
}

func (r *profilImpl) GetUserId(ctx context.Context, username string) (int64, error) {
	query := `SELECT id FROM auth.users WHERE LOWER(username) = LOWER($1)`
	var id int64

	err := r.db.GetContext(ctx, &id, query, username)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *profilImpl) GetProfile(ctx context.Context, id int64) (types.Profile, error) {
	query := `SELECT id, page, created_at FROM profiles.user_page WHERE user_id = $1`
	profile := types.Profile{}

	err := r.db.GetContext(ctx, &profile, query, id)
	if err != nil {
		return types.Profile{}, err
	}

	profile.UserID = id
	return profile, nil
}
