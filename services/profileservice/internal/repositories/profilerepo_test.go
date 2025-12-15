package repositories

import (
	"context"
	"database/sql"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func newTestProfileRepo(t *testing.T) (*profilImpl, sqlmock.Sqlmock, func()) {
	t.Helper()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to init sqlmock: %v", err)
	}

	sqlxdb := sqlx.NewDb(db, "sqlmock")
	repo := &profilImpl{db: sqlxdb}

	cleanup := func() {
		_ = sqlxdb.Close()
	}

	return repo, mock, cleanup
}

func TestGetUserId(t *testing.T) {
	assert := assert.New(t)
	repo, mock, cleanup := newTestProfileRepo(t)
	defer cleanup()

	query := regexp.QuoteMeta("SELECT id FROM auth.users WHERE username = $1")
	username := "admin"
	expectedID := int64(42)

	rows := sqlmock.NewRows([]string{"id"}).AddRow(expectedID)
	mock.ExpectQuery(query).WithArgs(username).WillReturnRows(rows)

	id, err := repo.GetUserId(context.Background(), username)
	assert.NoError(err)
	assert.Equal(expectedID, id)
	assert.NoError(mock.ExpectationsWereMet())
}

func TestGetUserIdError(t *testing.T) {
	assert := assert.New(t)
	repo, mock, cleanup := newTestProfileRepo(t)
	defer cleanup()

	query := regexp.QuoteMeta("SELECT id FROM auth.users WHERE username = $1")
	username := "unknown"

	mock.ExpectQuery(query).WithArgs(username).WillReturnError(sql.ErrNoRows)

	_, err := repo.GetUserId(context.Background(), username)
	assert.Error(err)
	assert.NoError(mock.ExpectationsWereMet())
}

func TestGetProfile(t *testing.T) {
	assert := assert.New(t)
	repo, mock, cleanup := newTestProfileRepo(t)
	defer cleanup()

	query := regexp.QuoteMeta("SELECT id, username, page, created_at FROM profiles.user_page WHERE user_id = $1")
	userID := int64(77)
	profileID := int64(10)
	username := "admin"
	page := []byte(`{"bio":"hello"}`)
	createdAt := time.Now()

	rows := sqlmock.NewRows([]string{"id", "username", "page", "created_at"}).
		AddRow(profileID, username, page, createdAt)
	mock.ExpectQuery(query).WithArgs(userID).WillReturnRows(rows)

	profile, err := repo.GetProfile(context.Background(), userID)
	assert.NoError(err)
	assert.Equal(profileID, profile.Id)
	assert.Equal(userID, profile.UserID)
	assert.Equal(username, profile.Username)
	assert.Equal("hello", profile.Page["bio"])
	assert.True(profile.CreatedAt.Equal(createdAt))
	assert.NoError(mock.ExpectationsWereMet())
}

func TestGetProfileError(t *testing.T) {
	assert := assert.New(t)
	repo, mock, cleanup := newTestProfileRepo(t)
	defer cleanup()

	query := regexp.QuoteMeta("SELECT id, username, page, created_at FROM profiles.user_page WHERE user_id = $1")
	userID := int64(100)

	mock.ExpectQuery(query).WithArgs(userID).WillReturnError(sql.ErrNoRows)

	_, err := repo.GetProfile(context.Background(), userID)
	assert.Error(err)
	assert.NoError(mock.ExpectationsWereMet())
}
