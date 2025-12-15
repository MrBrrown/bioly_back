package usecases

import (
	"context"
	"errors"
	"testing"
	"time"

	"bioly/profileservice/internal/types"

	"github.com/stretchr/testify/assert"
)

type mockProfileRepo struct {
	getUserIDFunc  func(ctx context.Context, username string) (int64, error)
	getProfileFunc func(ctx context.Context, id int64) (types.Profile, error)
}

func (m *mockProfileRepo) GetUserId(ctx context.Context, username string) (int64, error) {
	return m.getUserIDFunc(ctx, username)
}

func (m *mockProfileRepo) GetProfile(ctx context.Context, id int64) (types.Profile, error) {
	return m.getProfileFunc(ctx, id)
}

func TestProfileServiceGetProfileSuccess(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()

	repo := &mockProfileRepo{
		getUserIDFunc: func(ctx context.Context, username string) (int64, error) {
			assert.Equal("admin", username)
			return 42, nil
		},
		getProfileFunc: func(ctx context.Context, id int64) (types.Profile, error) {
			assert.Equal(int64(42), id)
			return types.Profile{
				Id:        10,
				UserID:    id,
				Username:  "admin",
				Page:      types.JSONB{"bio": "hello"},
				CreatedAt: time.Now(),
			}, nil
		},
	}

	service := NewProfile(repo, nil)
	profile, err := service.GetProfile(ctx, "admin")

	assert.NoError(err)
	assert.Equal(int64(10), profile.Id)
	assert.Equal(int64(42), profile.UserID)
	assert.Equal("admin", profile.Username)
	assert.Equal("hello", profile.Page["bio"])
}

func TestProfileServiceGetProfileUserIDError(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()

	expectedErr := errors.New("user not found")
	repo := &mockProfileRepo{
		getUserIDFunc: func(ctx context.Context, username string) (int64, error) {
			return 0, expectedErr
		},
		getProfileFunc: func(ctx context.Context, id int64) (types.Profile, error) {
			t.Fatalf("GetProfile should not be called on user id error")
			return types.Profile{}, nil
		},
	}

	service := NewProfile(repo, nil)
	_, err := service.GetProfile(ctx, "ghost")

	assert.ErrorIs(err, expectedErr)
}

func TestProfileServiceGetProfileLookupError(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()

	expectedErr := errors.New("profile missing")
	repo := &mockProfileRepo{
		getUserIDFunc: func(ctx context.Context, username string) (int64, error) {
			return 11, nil
		},
		getProfileFunc: func(ctx context.Context, id int64) (types.Profile, error) {
			assert.Equal(int64(11), id)
			return types.Profile{}, expectedErr
		},
	}

	service := NewProfile(repo, nil)
	_, err := service.GetProfile(ctx, "no-profile")

	assert.ErrorIs(err, expectedErr)
}
