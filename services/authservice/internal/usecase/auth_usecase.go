package usecase

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"strconv"
	"strings"
	"time"

	"github.com/alexedwards/argon2id"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"bioly/auth/internal/config"
	"bioly/auth/internal/repositories"
	"bioly/auth/internal/types"
)

type AuthService interface {
	Login(ctx context.Context, username, password, userAgent, ip string) (*types.User, *Tokens, error)
	Refresh(ctx context.Context, refreshToken, userAgent, ip string) (*types.User, *Tokens, error)
	CreateUser(ctx context.Context, username, password string) (*types.User, error)
	DeleteUser(ctx context.Context, id int64) error
}

type Tokens struct {
	Access  string
	Refresh string
}

type authImpl struct {
	users   repositories.Users
	rt      repositories.RefreshTokens
	jwtConf *config.JWT
	nowFn   func() time.Time
}

func NewAuth(users repositories.Users, rt repositories.RefreshTokens, jwtConf *config.JWT) AuthService {
	return &authImpl{
		users:   users,
		rt:      rt,
		jwtConf: jwtConf,
		nowFn:   func() time.Time { return time.Now().UTC() },
	}
}

func (a *authImpl) Login(ctx context.Context, username, password, userAgent, ip string) (*types.User, *Tokens, error) {
	user, err := a.users.VerifyCredentials(ctx, username, password)
	if err != nil {
		return nil, nil, err
	}
	now := a.nowFn()
	access, err := a.signAccess(user, now)
	if err != nil {
		return nil, nil, err
	}
	rtPlain, rtHash, jti, exp, err := a.newRefresh()
	if err != nil {
		return nil, nil, err
	}
	if err := a.rt.Create(ctx, user.ID, jti, rtHash, userAgent, ip, exp); err != nil && err != repositories.ErrNotImplemented {
		return nil, nil, err
	}
	return user, &Tokens{Access: access, Refresh: rtPlain}, nil
}

func (a *authImpl) Refresh(ctx context.Context, refreshToken, userAgent, ip string) (*types.User, *Tokens, error) {
	return nil, nil, repositories.ErrNotImplemented
}

func (a *authImpl) Logout(ctx context.Context, refreshToken string) error {
	return repositories.ErrNotImplemented
}

func (a *authImpl) LogoutAll(ctx context.Context, userID int64) error {
	return repositories.ErrNotImplemented
}

func (a *authImpl) signAccess(u *types.User, now time.Time) (string, error) {
	claims := jwt.MapClaims{
		"iss":  a.jwtConf.Issuer,
		"sub":  strconv.FormatInt(u.ID, 10),
		"name": u.Username,
		"iat":  now.Unix(),
		"exp":  now.Add(a.jwtConf.AccessTTL).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(a.jwtConf.AccessSecret))
}

func (a *authImpl) newRefresh() (plain string, hash string, jti uuid.UUID, exp time.Time, err error) {
	buf := make([]byte, 32)
	_, err = rand.Read(buf)
	if err != nil {
		return
	}
	plain = base64.RawURLEncoding.EncodeToString(buf)
	jti = uuid.New()
	exp = a.nowFn().Add(a.jwtConf.RefreshTTL)
	hash, err = argon2id.CreateHash(plain, argon2id.DefaultParams)
	return
}

func (a *authImpl) CreateUser(ctx context.Context, username, password string) (*types.User, error) {
	u := strings.TrimSpace(username)
	if len(u) < 3 || len(u) > 64 || password == "" {
		return nil, repositories.ErrInvalidCredentials
	}
	hash, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		return nil, err
	}
	user := &types.User{
		Username:     u,
		PasswordHash: hash,
	}
	if err := a.users.Add(ctx, user); err != nil {
		return nil, err
	}
	user.PasswordHash = ""
	return user, nil
}

func (a *authImpl) DeleteUser(ctx context.Context, id int64) error {
	return a.users.Delete(ctx, id)
}
