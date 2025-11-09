package transport_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/stretchr/testify/assert"

	"bioly/auth/internal/repositories"
	"bioly/auth/internal/transport"
	"bioly/auth/internal/types"
	"bioly/auth/internal/usecase"
)

// --- mock AuthService ---

type authMock struct {
	loginFn      func(username, password, ua, ip string) (*types.User, *usecase.Tokens, error)
	refreshFn    func(refresh, ua, ip string) (*types.User, *usecase.Tokens, error)
	createUserFn func(username, password string) (*types.User, error)
	deleteUserFn func(id int64) error
}

func (m *authMock) Login(_ ctx, username, password, ua, ip string) (*types.User, *usecase.Tokens, error) {
	return m.loginFn(username, password, ua, ip)
}
func (m *authMock) Refresh(_ ctx, refresh, ua, ip string) (*types.User, *usecase.Tokens, error) {
	return m.refreshFn(refresh, ua, ip)
}
func (m *authMock) CreateUser(_ ctx, username, password string) (*types.User, error) {
	return m.createUserFn(username, password)
}
func (m *authMock) DeleteUser(_ ctx, id int64) error {
	return m.deleteUserFn(id)
}

type ctx = context.Context

// --- helpers ---

func makeRouter(h *transport.Handler) http.Handler {
	r := chi.NewRouter()
	r.Use(render.SetContentType(render.ContentTypeJSON))
	h.RegisterRoutes(r)
	return r
}

func doJSON(t *testing.T, router http.Handler, method, path string, body any, headers map[string]string) *httptest.ResponseRecorder {
	t.Helper()
	var buf bytes.Buffer
	if body != nil {
		enc := json.NewEncoder(&buf)
		_ = enc.Encode(body)
	}
	req := httptest.NewRequest(method, path, &buf)
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	if _, ok := headers["Content-Type"]; !ok {
		req.Header.Set("Content-Type", "application/json")
	}
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	return w
}

// --- tests ---

func TestLogin_Success(t *testing.T) {
	m := &authMock{
		loginFn: func(username, password, ua, ip string) (*types.User, *usecase.Tokens, error) {
			assert.Equal(t, "admin", username)
			assert.Equal(t, "secret", password)
			return &types.User{ID: 1, Username: "admin", CreatedAt: time.Now().UTC()}, &usecase.Tokens{
				Access:  "access.jwt",
				Refresh: "refresh.token",
			}, nil
		},
	}
	h := transport.NewHandler(m)
	router := makeRouter(h)

	w := doJSON(t, router, http.MethodPost, "/login", map[string]string{
		"username": "admin",
		"password": "secret",
	}, map[string]string{"User-Agent": "UA", "X-Real-IP": "127.0.0.1"})

	assert.Equal(t, http.StatusOK, w.Code)
	var resp struct {
		Access  string        `json:"access"`
		Refresh string        `json:"refresh"`
		User    types.UserDTO `json:"user"`
	}
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, "access.jwt", resp.Access)
	assert.Equal(t, "refresh.token", resp.Refresh)
	assert.Equal(t, int64(1), resp.User.ID)
	assert.Equal(t, "admin", resp.User.Username)
}

func TestLogin_BadJSON(t *testing.T) {
	m := &authMock{loginFn: func(_, _, _, _ string) (*types.User, *usecase.Tokens, error) { return nil, nil, nil }}
	h := transport.NewHandler(m)
	router := makeRouter(h)

	w := doJSON(t, router, http.MethodPost, "/login", map[string]any{
		"username": "admin",
		// password отсутствует
	}, nil)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestRefresh_Success(t *testing.T) {
	m := &authMock{
		refreshFn: func(refresh, ua, ip string) (*types.User, *usecase.Tokens, error) {
			assert.Equal(t, "rt-123", refresh)
			return &types.User{ID: 2, Username: "u2"}, &usecase.Tokens{Access: "new.access", Refresh: "new.refresh"}, nil
		},
	}
	h := transport.NewHandler(m)
	router := makeRouter(h)

	w := doJSON(t, router, http.MethodPost, "/refresh", map[string]string{
		"refresh": "rt-123",
	}, map[string]string{"User-Agent": "UA", "X-Forwarded-For": "10.0.0.1"})

	assert.Equal(t, http.StatusOK, w.Code)
	var resp struct {
		Access  string        `json:"access"`
		Refresh string        `json:"refresh"`
		User    types.UserDTO `json:"user"`
	}
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, "new.access", resp.Access)
	assert.Equal(t, "new.refresh", resp.Refresh)
	assert.Equal(t, int64(2), resp.User.ID)
}

func TestRefresh_BadJSON(t *testing.T) {
	m := &authMock{refreshFn: func(_, _, _ string) (*types.User, *usecase.Tokens, error) { return nil, nil, nil }}
	h := transport.NewHandler(m)
	router := makeRouter(h)

	w := doJSON(t, router, http.MethodPost, "/refresh", map[string]any{
		// нет поля refresh
	}, nil)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreateUser_Success(t *testing.T) {
	m := &authMock{
		createUserFn: func(username, password string) (*types.User, error) {
			assert.Equal(t, "newbie", username)
			assert.Equal(t, "pass", password)
			return &types.User{ID: 10, Username: "newbie"}, nil
		},
	}
	h := transport.NewHandler(m)
	router := makeRouter(h)

	w := doJSON(t, router, http.MethodPost, "/users", map[string]string{
		"username": "newbie",
		"password": "pass",
	}, nil)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp struct {
		User types.UserDTO `json:"user"`
	}
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, int64(10), resp.User.ID)
	assert.Equal(t, "newbie", resp.User.Username)
}

func TestCreateUser_Conflict(t *testing.T) {
	m := &authMock{
		createUserFn: func(username, password string) (*types.User, error) {
			return nil, repositories.ErrDuplicateUsername
		},
	}
	h := transport.NewHandler(m)
	router := makeRouter(h)

	w := doJSON(t, router, http.MethodPost, "/users", map[string]string{
		"username": "admin",
		"password": "x",
	}, nil)

	assert.Equal(t, http.StatusConflict, w.Code)
}

func TestDeleteUser_Success(t *testing.T) {
	m := &authMock{
		deleteUserFn: func(id int64) error {
			assert.Equal(t, int64(5), id)
			return nil
		},
	}
	h := transport.NewHandler(m)
	router := makeRouter(h)

	w := doJSON(t, router, http.MethodDelete, "/users/"+strconv.FormatInt(5, 10), nil, nil)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDeleteUser_NotFound(t *testing.T) {
	m := &authMock{
		deleteUserFn: func(id int64) error { return repositories.ErrNotFound },
	}
	h := transport.NewHandler(m)
	router := makeRouter(h)

	w := doJSON(t, router, http.MethodDelete, "/users/999", nil, nil)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestDeleteUser_BadID(t *testing.T) {
	m := &authMock{deleteUserFn: func(id int64) error { return nil }}
	h := transport.NewHandler(m)
	router := makeRouter(h)

	w := doJSON(t, router, http.MethodDelete, "/users/abc", nil, nil)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}
