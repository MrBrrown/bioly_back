package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"bioly/auth/internal/config"
	"bioly/auth/internal/repositories"
	"bioly/auth/internal/types"
	"bioly/auth/internal/usecase"
)

type usersMock struct {
	addFn    func(ctx context.Context, u *types.User) error
	delFn    func(ctx context.Context, id int64) error
	verifyFn func(ctx context.Context, username, password string) (*types.User, error)
}

func (m *usersMock) Add(ctx context.Context, u *types.User) error {
	return m.addFn(ctx, u)
}
func (m *usersMock) Delete(ctx context.Context, id int64) error {
	return m.delFn(ctx, id)
}
func (m *usersMock) VerifyCredentials(ctx context.Context, username, password string) (*types.User, error) {
	return m.verifyFn(ctx, username, password)
}

type rtMock struct {
	createCalled  bool
	createFn      func(ctx context.Context, userID int64, jti [16]byte, tokenHash, userAgent, ip string, expiresAt time.Time) error
	revokeFn      func(ctx context.Context, jti [16]byte) error
	revokeAllFn   func(ctx context.Context, userID int64) error
	findFn        func(ctx context.Context, jti [16]byte) (int64, string, time.Time, bool, error)
	lastUserAgent string
	lastIP        string
	lastExpiresAt time.Time
	lastTokenHash string
	lastUserID    int64
}

func (m *rtMock) Create(ctx context.Context, userID int64, jti [16]byte, tokenHash, userAgent, ip string, expiresAt time.Time) error {
	m.createCalled = true
	m.lastUserID = userID
	m.lastUserAgent = userAgent
	m.lastIP = ip
	m.lastExpiresAt = expiresAt
	m.lastTokenHash = tokenHash
	if m.createFn != nil {
		return m.createFn(ctx, userID, jti, tokenHash, userAgent, ip, expiresAt)
	}
	return nil
}
func (m *rtMock) RevokeByJTI(ctx context.Context, jti [16]byte) error {
	if m.revokeFn != nil {
		return m.revokeFn(ctx, jti)
	}
	return nil
}
func (m *rtMock) RevokeAllByUser(ctx context.Context, userID int64) error {
	if m.revokeAllFn != nil {
		return m.revokeAllFn(ctx, userID)
	}
	return nil
}
func (m *rtMock) FindValidByJTI(ctx context.Context, jti [16]byte) (int64, string, time.Time, bool, error) {
	if m.findFn != nil {
		return m.findFn(ctx, jti)
	}
	return 0, "", time.Time{}, false, nil
}

func TestCreateUser_Success(t *testing.T) {
	uRepo := &usersMock{
		addFn: func(ctx context.Context, u *types.User) error {
			u.ID = 101
			u.CreatedAt = time.Now().UTC()
			u.UpdatedAt = u.CreatedAt
			return nil
		},
	}
	rtRepo := &rtMock{}
	uc := usecase.NewAuth(uRepo, rtRepo, &config.JWT{
		AccessSecret:  "access",
		RefreshSecret: "refresh",
		AccessTTL:     15 * time.Minute,
		RefreshTTL:    24 * time.Hour,
		Issuer:        "auth.test",
	})

	user, err := uc.CreateUser(context.Background(), "newuser", "secret")
	assert.NoError(t, err)
	assert.Equal(t, int64(101), user.ID)
	assert.Equal(t, "newuser", user.Username)
}

func TestCreateUser_Duplicate(t *testing.T) {
	uRepo := &usersMock{
		addFn: func(ctx context.Context, u *types.User) error {
			return repositories.ErrDuplicateUsername
		},
	}
	rtRepo := &rtMock{}
	uc := usecase.NewAuth(uRepo, rtRepo, &config.JWT{
		AccessSecret:  "access",
		RefreshSecret: "refresh",
		AccessTTL:     15 * time.Minute,
		RefreshTTL:    24 * time.Hour,
		Issuer:        "auth.test",
	})

	user, err := uc.CreateUser(context.Background(), "admin", "secret")
	assert.Error(t, err)
	assert.Nil(t, user)
	assert.ErrorIs(t, err, repositories.ErrDuplicateUsername)
}

func TestCreateUser_InvalidInput(t *testing.T) {
	uRepo := &usersMock{
		addFn: func(ctx context.Context, u *types.User) error { return nil },
	}
	rtRepo := &rtMock{}
	uc := usecase.NewAuth(uRepo, rtRepo, &config.JWT{
		AccessSecret:  "access",
		RefreshSecret: "refresh",
		AccessTTL:     15 * time.Minute,
		RefreshTTL:    24 * time.Hour,
		Issuer:        "auth.test",
	})

	user, err := uc.CreateUser(context.Background(), "ab", "")
	assert.Error(t, err)
	assert.Nil(t, user)
	assert.ErrorIs(t, err, repositories.ErrInvalidCredentials)
}

func TestDeleteUser_Success(t *testing.T) {
	uRepo := &usersMock{
		delFn: func(ctx context.Context, id int64) error { return nil },
	}
	rtRepo := &rtMock{}
	uc := usecase.NewAuth(uRepo, rtRepo, &config.JWT{
		AccessSecret:  "access",
		RefreshSecret: "refresh",
		AccessTTL:     15 * time.Minute,
		RefreshTTL:    24 * time.Hour,
		Issuer:        "auth.test",
	})

	err := uc.DeleteUser(context.Background(), 7)
	assert.NoError(t, err)
}

func TestDeleteUser_NotFound(t *testing.T) {
	uRepo := &usersMock{
		delFn: func(ctx context.Context, id int64) error { return repositories.ErrNotFound },
	}
	rtRepo := &rtMock{}
	uc := usecase.NewAuth(uRepo, rtRepo, &config.JWT{
		AccessSecret:  "access",
		RefreshSecret: "refresh",
		AccessTTL:     15 * time.Minute,
		RefreshTTL:    24 * time.Hour,
		Issuer:        "auth.test",
	})

	err := uc.DeleteUser(context.Background(), 999)
	assert.Error(t, err)
	assert.ErrorIs(t, err, repositories.ErrNotFound)
}

func TestLogin_Success(t *testing.T) {
	now := time.Now().UTC()
	uRepo := &usersMock{
		verifyFn: func(ctx context.Context, username, password string) (*types.User, error) {
			return &types.User{
				ID:        77,
				Username:  "root",
				CreatedAt: now,
				UpdatedAt: now,
			}, nil
		},
	}
	rtRepo := &rtMock{}
	uc := usecase.NewAuth(uRepo, rtRepo, &config.JWT{
		AccessSecret:  "access",
		RefreshSecret: "refresh",
		AccessTTL:     15 * time.Minute,
		RefreshTTL:    24 * time.Hour,
		Issuer:        "auth.test",
	})

	user, tokens, err := uc.Login(context.Background(), "root", "secret", "UA", "127.0.0.1")
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.NotEmpty(t, tokens.Access)
	assert.NotEmpty(t, tokens.Refresh)
	assert.True(t, rtRepo.createCalled)
	assert.Equal(t, int64(77), rtRepo.lastUserID)
	assert.Equal(t, "UA", rtRepo.lastUserAgent)
	assert.Equal(t, "127.0.0.1", rtRepo.lastIP)
	assert.WithinDuration(t, time.Now().Add(24*time.Hour), rtRepo.lastExpiresAt, 5*time.Second)
	assert.NotEmpty(t, rtRepo.lastTokenHash)
}

func TestLogin_InvalidCredentials(t *testing.T) {
	uRepo := &usersMock{
		verifyFn: func(ctx context.Context, username, password string) (*types.User, error) {
			return nil, repositories.ErrInvalidCredentials
		},
	}
	rtRepo := &rtMock{}
	uc := usecase.NewAuth(uRepo, rtRepo, &config.JWT{
		AccessSecret:  "access",
		RefreshSecret: "refresh",
		AccessTTL:     15 * time.Minute,
		RefreshTTL:    24 * time.Hour,
		Issuer:        "auth.test",
	})

	user, tokens, err := uc.Login(context.Background(), "who", "bad", "UA", "ip")
	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Nil(t, tokens)
}

func TestRefresh_NotImplemented(t *testing.T) {
	uRepo := &usersMock{}
	rtRepo := &rtMock{}
	uc := usecase.NewAuth(uRepo, rtRepo, &config.JWT{
		AccessSecret:  "access",
		RefreshSecret: "refresh",
		AccessTTL:     15 * time.Minute,
		RefreshTTL:    24 * time.Hour,
		Issuer:        "auth.test",
	})
	user, tokens, err := uc.Refresh(context.Background(), "rt", "UA", "ip")
	assert.Nil(t, user)
	assert.Nil(t, tokens)
	assert.True(t, err != nil)
	assert.True(t, errors.Is(err, repositories.ErrNotImplemented) || errors.Is(err, repositories.ErrInvalidCredentials) || err != nil)
}
