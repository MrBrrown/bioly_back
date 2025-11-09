package transport

import (
	"bioly/asynclogger"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

func NewRouter(handler *Handler) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RequestLogger(&ChiAdapter{
		Formatter: &asynclogger.DefaultFormatter{},
	}))
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)
	r.Use(render.SetContentType(render.ContentTypeJSON))

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		asynclogger.Warning("[%s] not found %s %s", middleware.GetReqID(r.Context()), r.Method, r.URL.Path)
		http.Error(w, "not found", http.StatusNotFound)
	})
	r.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		asynclogger.Warning("[%s] method not allowed %s %s", middleware.GetReqID(r.Context()), r.Method, r.URL.Path)
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	})

	handler.RegisterRoutes(r)
	return r
}
