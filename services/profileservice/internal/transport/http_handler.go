package transport

import (
	"bioly/common/asynclogger"
	"bioly/profileservice/internal/types"
	"bioly/profileservice/internal/usecases"
	"fmt"
	"math/rand/v2"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type Handler struct {
	profile usecases.ProfileService

	randomNames []string
}

func NewHandler(p usecases.ProfileService) *Handler {
	handler := &Handler{profile: p}

	handler.randomNames = make([]string, 0, 100_000)
	for i := range 100_000 {
		handler.randomNames = append(handler.randomNames, fmt.Sprintf("user_%d", i))
	}

	return handler
}

func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Get("/ping", h.ping)
	r.Get("/health", h.health)

	r.Get("/{username}", h.getProfile)

	r.Get("/internal/randompage", h.testGetProfile)
	r.Get("/internal/CachedRandomPage", h.testGetProfileCached)
}

func (h *Handler) ping(w http.ResponseWriter, r *http.Request) {
	render.PlainText(w, r, "pong")
}

func (h *Handler) health(w http.ResponseWriter, r *http.Request) {
	render.PlainText(w, r, "OK")
}

func (h *Handler) getProfile(w http.ResponseWriter, r *http.Request) {
	reqID := middleware.GetReqID(r.Context())
	start := time.Now().UTC()

	idStr := chi.URLParam(r, "username")

	profile, err := h.profile.GetProfile(r.Context(), idStr)
	if err != nil {
		asynclogger.Error("[%s] failed to get profile for username %s: %v", reqID, idStr, err)
		render.Render(w, r, types.ErrInvalidRequest(http.StatusNotFound, fmt.Errorf("profile not found")))
		return
	}

	resp := &types.ProfileResponse{
		Username: profile.Username,
		Page:     profile.Page,
	}

	elapsed := time.Since(start)
	asynclogger.Info("[%s] get profile for username %s in %v", reqID, idStr, elapsed)
	render.Render(w, r, resp)
}

func (h *Handler) testGetProfile(w http.ResponseWriter, r *http.Request) {
	reqID := middleware.GetReqID(r.Context())
	start := time.Now().UTC()

	idStr := h.randomNames[rand.IntN(len(h.randomNames))]

	profile, err := h.profile.GetProfile(r.Context(), idStr)
	if err != nil {
		asynclogger.Error("[%s] failed to get profile for username %s: %v", reqID, idStr, err)
		render.Render(w, r, types.ErrInvalidRequest(http.StatusNotFound, fmt.Errorf("profile not found")))
		return
	}

	resp := &types.ProfileResponse{
		Username: profile.Username,
		Page:     profile.Page,
	}

	elapsed := time.Since(start)
	asynclogger.Info("[%s] get profile for username %s in %v", reqID, idStr, elapsed)
	render.Render(w, r, resp)
}

func (h *Handler) testGetProfileCached(w http.ResponseWriter, r *http.Request) {
	reqID := middleware.GetReqID(r.Context())
	start := time.Now().UTC()

	idStr := h.randomNames[rand.IntN(len(h.randomNames))]

	profile, err := h.profile.GetProfileCached(r.Context(), idStr)
	if err != nil {
		asynclogger.Error("[%s] failed to get profile for username %s: %v", reqID, idStr, err)
		render.Render(w, r, types.ErrInvalidRequest(http.StatusNotFound, fmt.Errorf("profile not found")))
		return
	}

	resp := &types.ProfileResponse{
		Username: profile.Username,
		Page:     profile.Page,
	}

	elapsed := time.Since(start)
	asynclogger.Info("[%s] get profile for username %s in %v", reqID, idStr, elapsed)
	render.Render(w, r, resp)
}
