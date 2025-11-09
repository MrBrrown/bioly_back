package repositories

import (
	"context"
	"database/sql"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/alexedwards/argon2id"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"

	"bioly/auth/internal/types"
)

func TestUsers_Add_Success(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
	assert.NoError(t, err)
	defer db.Close()

	xdb := sqlx.NewDb(db, "sqlmock")
	defer xdb.Close()

	repo := NewUsers(xdb)

	u := &types.User{Username: "admin", PasswordHash: "$argon2id$v=19$m=65536,t=3,p=2$SALT$HASH"}

	q := regexp.QuoteMeta(`
		INSERT INTO users (username, password_hash)
		VALUES ($1, $2)
		RETURNING id, created_at, updated_at
	`)
	now := time.Now().UTC()
	rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).
		AddRow(int64(1), now, now)

	mock.ExpectQuery(q).
		WithArgs(u.Username, u.PasswordHash).
		WillReturnRows(rows)

	err = repo.Add(context.Background(), u)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), u.ID)
	assert.WithinDuration(t, now, u.CreatedAt, time.Second)
	assert.WithinDuration(t, now, u.UpdatedAt, time.Second)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUsers_Add_Duplicate(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
	assert.NoError(t, err)
	defer db.Close()

	xdb := sqlx.NewDb(db, "sqlmock")
	defer xdb.Close()

	repo := NewUsers(xdb)

	u := &types.User{Username: "admin", PasswordHash: "hash"}

	q := regexp.QuoteMeta(`
		INSERT INTO users (username, password_hash)
		VALUES ($1, $2)
		RETURNING id, created_at, updated_at
	`)
	mock.ExpectQuery(q).
		WithArgs(u.Username, u.PasswordHash).
		WillReturnError(&pq.Error{Code: "23505"})

	err = repo.Add(context.Background(), u)
	assert.Error(t, err)
	assert.ErrorIs(t, err, ErrDuplicateUsername)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUsers_Delete_Success(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
	assert.NoError(t, err)
	defer db.Close()

	xdb := sqlx.NewDb(db, "sqlmock")
	defer xdb.Close()

	repo := NewUsers(xdb)

	q := regexp.QuoteMeta(`DELETE FROM users WHERE id = $1`)
	mock.ExpectExec(q).
		WithArgs(int64(42)).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = repo.Delete(context.Background(), 42)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUsers_Delete_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
	assert.NoError(t, err)
	defer db.Close()

	xdb := sqlx.NewDb(db, "sqlmock")
	defer xdb.Close()

	repo := NewUsers(xdb)

	q := regexp.QuoteMeta(`DELETE FROM users WHERE id = $1`)
	mock.ExpectExec(q).
		WithArgs(int64(99)).
		WillReturnResult(sqlmock.NewResult(0, 0))

	err = repo.Delete(context.Background(), 99)
	assert.Error(t, err)
	assert.ErrorIs(t, err, ErrNotFound)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUsers_VerifyCredentials_Success(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
	assert.NoError(t, err)
	defer db.Close()

	xdb := sqlx.NewDb(db, "sqlmock")
	defer xdb.Close()

	repo := NewUsers(xdb)

	hash, herr := argon2id.CreateHash("secret", argon2id.DefaultParams)
	assert.NoError(t, herr)

	sel := regexp.QuoteMeta(`
		SELECT id, username, password_hash, last_login_at, created_at, updated_at
		FROM users
		WHERE lower(username) = lower($1)
		LIMIT 1
	`)
	now := time.Now().UTC()
	rows := sqlmock.NewRows([]string{"id", "username", "password_hash", "last_login_at", "created_at", "updated_at"}).
		AddRow(int64(7), "Admin", hash, sql.NullTime{}, now, now)

	mock.ExpectQuery(sel).
		WithArgs("admin").
		WillReturnRows(rows)

	upd := regexp.QuoteMeta(`
		UPDATE users
		SET last_login_at = NOW(), updated_at = NOW()
		WHERE id = $1
	`)
	mock.ExpectExec(upd).
		WithArgs(int64(7)).
		WillReturnResult(sqlmock.NewResult(0, 1))

	u, err := repo.VerifyCredentials(context.Background(), "admin", "secret")
	assert.NoError(t, err)
	assert.Equal(t, int64(7), u.ID)
	assert.NotNil(t, u.LastLoginAt)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUsers_VerifyCredentials_InvalidPassword(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
	assert.NoError(t, err)
	defer db.Close()

	xdb := sqlx.NewDb(db, "sqlmock")
	defer xdb.Close()

	repo := NewUsers(xdb)

	hash, herr := argon2id.CreateHash("correct", argon2id.DefaultParams)
	assert.NoError(t, herr)

	sel := regexp.QuoteMeta(`
		SELECT id, username, password_hash, last_login_at, created_at, updated_at
		FROM users
		WHERE lower(username) = lower($1)
		LIMIT 1
	`)
	now := time.Now().UTC()
	rows := sqlmock.NewRows([]string{"id", "username", "password_hash", "last_login_at", "created_at", "updated_at"}).
		AddRow(int64(3), "user", hash, sql.NullTime{}, now, now)

	mock.ExpectQuery(sel).
		WithArgs("user").
		WillReturnRows(rows)

	_, err = repo.VerifyCredentials(context.Background(), "user", "wrong")
	assert.Error(t, err)
	assert.ErrorIs(t, err, ErrInvalidCredentials)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUsers_VerifyCredentials_NoUser(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
	assert.NoError(t, err)
	defer db.Close()

	xdb := sqlx.NewDb(db, "sqlmock")
	defer xdb.Close()

	repo := NewUsers(xdb)

	sel := regexp.QuoteMeta(`
		SELECT id, username, password_hash, last_login_at, created_at, updated_at
		FROM users
		WHERE lower(username) = lower($1)
		LIMIT 1
	`)
	mock.ExpectQuery(sel).
		WithArgs("nouser").
		WillReturnError(sql.ErrNoRows)

	_, err = repo.VerifyCredentials(context.Background(), "nouser", "x")
	assert.Error(t, err)
	assert.ErrorIs(t, err, ErrInvalidCredentials)
	assert.NoError(t, mock.ExpectationsWereMet())
}
