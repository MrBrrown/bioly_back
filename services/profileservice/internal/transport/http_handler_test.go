package transport

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"bioly/profileservice/internal/types"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/stretchr/testify/assert"
)

type mockProfileService struct {
	getProfileFunc func(ctx context.Context, username string) (types.Profile, error)
}

func (m *mockProfileService) GetProfile(ctx context.Context, username string) (types.Profile, error) {
	if m.getProfileFunc != nil {
		return m.getProfileFunc(ctx, username)
	}
	return types.Profile{}, nil
}

func (m *mockProfileService) GetProfileCached(ctx context.Context, username string) (types.Profile, error) {
	if m.getProfileFunc != nil {
		return m.getProfileFunc(ctx, username)
	}
	return types.Profile{}, nil
}

func newTestRouter(t *testing.T, svc *mockProfileService) chi.Router {
	t.Helper()

	r := chi.NewRouter()
	r.Use(middleware.RequestID)

	NewHandler(svc).RegisterRoutes(r)
	return r
}

func TestHandlerPing(t *testing.T) {
	router := newTestRouter(t, &mockProfileService{})
	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "text/plain; charset=utf-8", rec.Header().Get("Content-Type"))
	assert.Equal(t, "pong", rec.Body.String())
}

func TestHandlerHealth(t *testing.T) {
	router := newTestRouter(t, &mockProfileService{})
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "text/plain; charset=utf-8", rec.Header().Get("Content-Type"))
	assert.Equal(t, "OK", rec.Body.String())
}

func TestHandlerGetProfileSuccess(t *testing.T) {
	page := types.JSONB{"bio": "hello"}
	mockSvc := &mockProfileService{
		getProfileFunc: func(ctx context.Context, username string) (types.Profile, error) {
			assert.Equal(t, "john", username)
			return types.Profile{
				Id:        1,
				UserID:    42,
				Username:  username,
				Page:      page,
				CreatedAt: time.Now(),
			}, nil
		},
	}

	router := newTestRouter(t, mockSvc)
	req := httptest.NewRequest(http.MethodGet, "/john", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))

	var resp types.ProfileResponse
	err := json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "john", resp.Username)
	assert.Equal(t, "hello", resp.Page["bio"])
}

func TestHandlerGetProfileError(t *testing.T) {
	expectedErr := errors.New("profile not found")
	mockSvc := &mockProfileService{
		getProfileFunc: func(ctx context.Context, username string) (types.Profile, error) {
			return types.Profile{}, expectedErr
		},
	}

	router := newTestRouter(t, mockSvc)
	req := httptest.NewRequest(http.MethodGet, "/ghost", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))

	var resp map[string]string
	err := json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, expectedErr.Error(), resp["error"])
}
