package transport

import (
	"fmt"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"

	"bioly/asynclogger"
	"bioly/auth/internal/repositories"
	"bioly/auth/internal/types"
	"bioly/auth/internal/usecase"
)

type Handler struct {
	auth usecase.AuthService
}

func NewHandler(a usecase.AuthService) *Handler {
	return &Handler{auth: a}
}

func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Get("/ping", h.ping)
	r.Get("/health", h.health)
	r.Post("/login", h.login)
	r.Post("/refresh", h.refresh)
	r.Post("/users", h.createUser)
	r.Delete("/users/{id}", h.deleteUser)
}

func (h *Handler) ping(w http.ResponseWriter, r *http.Request) {
	render.PlainText(w, r, "pong")
}

func (h *Handler) health(w http.ResponseWriter, r *http.Request) {
	render.PlainText(w, r, "OK")
}

func clientIP(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		parts := strings.Split(xff, ",")
		return strings.TrimSpace(parts[0])
	}
	if xr := r.Header.Get("X-Real-IP"); xr != "" {
		return xr
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (lr *loginRequest) Bind(r *http.Request) error {
	if lr.Username == "" || lr.Password == "" {
		return fmt.Errorf("username and password are required")
	}
	return nil
}

type refreshRequest struct {
	Refresh string `json:"refresh"`
}

func (rr *refreshRequest) Bind(r *http.Request) error {
	if rr.Refresh == "" {
		return fmt.Errorf("refresh token is required")
	}
	return nil
}

type createUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (c *createUserRequest) Bind(r *http.Request) error {
	if c.Username == "" || c.Password == "" {
		return fmt.Errorf("username and password are required")
	}
	return nil
}

type okResponse struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
}

func (o *okResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (h *Handler) login(w http.ResponseWriter, r *http.Request) {
	reqID := middleware.GetReqID(r.Context())
	start := time.Now().UTC()

	var req loginRequest
	if err := render.Bind(r, &req); err != nil {
		asynclogger.Warning("[%s] login bind failed ip=%s ua=%q err=%v", reqID, clientIP(r), r.Header.Get("User-Agent"), err)
		render.Render(w, r, types.ErrInvalidRequest(http.StatusBadRequest, fmt.Errorf("invalid request body")))
		return
	}

	ua := r.Header.Get("User-Agent")
	ip := clientIP(r)
	user, tokens, err := h.auth.Login(r.Context(), req.Username, req.Password, ua, ip)
	if err != nil {
		asynclogger.Warning("[%s] login failed ip=%s ua=%q username=%q dur=%s err=%v", reqID, ip, ua, req.Username, time.Since(start), err)
		render.Render(w, r, types.ErrInvalidRequest(http.StatusUnauthorized, err))
		return
	}

	asynclogger.Info("[%s] login success ip=%s ua=%q user_id=%d username=%q dur=%s", reqID, ip, ua, user.ID, user.Username, time.Since(start))

	resp := &types.LoginResponse{
		Access:  tokens.Access,
		Refresh: tokens.Refresh,
		User:    types.UserDTO{ID: user.ID, Username: user.Username},
	}
	render.Render(w, r, resp)
}

func (h *Handler) refresh(w http.ResponseWriter, r *http.Request) {
	reqID := middleware.GetReqID(r.Context())
	start := time.Now().UTC()

	var req refreshRequest
	if err := render.Bind(r, &req); err != nil {
		asynclogger.Warning("[%s] refresh bind failed ip=%s ua=%q err=%v", reqID, clientIP(r), r.Header.Get("User-Agent"), err)
		render.Render(w, r, types.ErrInvalidRequest(http.StatusBadRequest, fmt.Errorf("invalid request body")))
		return
	}

	ua := r.Header.Get("User-Agent")
	ip := clientIP(r)
	user, tokens, err := h.auth.Refresh(r.Context(), req.Refresh, ua, ip)
	if err != nil {
		asynclogger.Warning("[%s] refresh failed ip=%s ua=%q dur=%s err=%v", reqID, ip, ua, time.Since(start), err)
		render.Render(w, r, types.ErrInvalidRequest(http.StatusNotImplemented, err))
		return
	}

	asynclogger.Info("[%s] refresh success ip=%s ua=%q user_id=%d dur=%s", reqID, ip, ua, user.ID, time.Since(start))

	resp := &types.LoginResponse{
		Access:  tokens.Access,
		Refresh: tokens.Refresh,
		User:    types.UserDTO{ID: user.ID, Username: user.Username},
	}
	render.Render(w, r, resp)
}

func (h *Handler) createUser(w http.ResponseWriter, r *http.Request) {
	reqID := middleware.GetReqID(r.Context())
	start := time.Now().UTC()

	var req createUserRequest
	if err := render.Bind(r, &req); err != nil {
		asynclogger.Warning("[%s] createUser bind failed ip=%s ua=%q err=%v", reqID, clientIP(r), r.Header.Get("User-Agent"), err)
		render.Render(w, r, types.ErrInvalidRequest(http.StatusBadRequest, fmt.Errorf("invalid request body")))
		return
	}

	user, err := h.auth.CreateUser(r.Context(), req.Username, req.Password)
	if err != nil {
		switch err {
		case repositories.ErrDuplicateUsername:
			asynclogger.Warning("[%s] createUser duplicate username=%q dur=%s", reqID, req.Username, time.Since(start))
			render.Render(w, r, types.ErrInvalidRequest(http.StatusConflict, err))
			return
		case repositories.ErrInvalidCredentials:
			asynclogger.Warning("[%s] createUser invalid input username=%q dur=%s", reqID, req.Username, time.Since(start))
			render.Render(w, r, types.ErrInvalidRequest(http.StatusBadRequest, err))
			return
		default:
			asynclogger.Error("[%s] createUser failed username=%q dur=%s err=%v", reqID, req.Username, time.Since(start), err)
			render.Render(w, r, types.ErrInvalidRequest(http.StatusInternalServerError, fmt.Errorf("internal error")))
			return
		}
	}

	asynclogger.Info("[%s] createUser success user_id=%d username=%q dur=%s", reqID, user.ID, user.Username, time.Since(start))
	render.Render(w, r, &types.LoginResponse{
		Access:  "",
		Refresh: "",
		User:    types.UserDTO{ID: user.ID, Username: user.Username},
	})
}

func (h *Handler) deleteUser(w http.ResponseWriter, r *http.Request) {
	reqID := middleware.GetReqID(r.Context())
	start := time.Now().UTC()

	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		asynclogger.Warning("[%s] deleteUser bad id=%q", reqID, idStr)
		render.Render(w, r, types.ErrInvalidRequest(http.StatusBadRequest, fmt.Errorf("invalid id")))
		return
	}

	if err := h.auth.DeleteUser(r.Context(), id); err != nil {
		if err == repositories.ErrNotFound {
			asynclogger.Warning("[%s] deleteUser not found id=%d dur=%s", reqID, id, time.Since(start))
			render.Render(w, r, types.ErrInvalidRequest(http.StatusNotFound, err))
			return
		}
		asynclogger.Error("[%s] deleteUser failed id=%d dur=%s err=%v", reqID, id, time.Since(start), err)
		render.Render(w, r, types.ErrInvalidRequest(http.StatusInternalServerError, fmt.Errorf("internal error")))
		return
	}

	asynclogger.Info("[%s] deleteUser success id=%d dur=%s", reqID, id, time.Since(start))
	render.Render(w, r, &okResponse{Status: "ok", Message: "user deleted"})
}
