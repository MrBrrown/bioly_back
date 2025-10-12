package server

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"bioly/gateway/internal/transport"
)

type CORSConfig struct {
	AllowedOrigins   []string      // Access-Control-Allow-Origin
	AllowedMethods   []string      // Access-Control-Allow-Methods
	AllowedHeaders   []string      // Access-Control-Allow-Headers
	AllowCredentials bool          // Access-Control-Allow-Credentials
	MaxAge           time.Duration // Access-Control-Max-Age
}

func NewRouter(h *transport.Handler, requestTimeout time.Duration, corsCfg CORSConfig) http.Handler {
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	handlerWithTimeout := http.TimeoutHandler(mux, requestTimeout, "request timed out")

	return applyCORS(handlerWithTimeout, corsCfg)
}

func applyCORS(next http.Handler, cfg CORSConfig) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if len(cfg.AllowedOrigins) > 0 {
			w.Header().Set("Access-Control-Allow-Origin", strings.Join(cfg.AllowedOrigins, ","))
		}
		if len(cfg.AllowedMethods) > 0 {
			w.Header().Set("Access-Control-Allow-Methods", strings.Join(cfg.AllowedMethods, ","))
		}
		if len(cfg.AllowedHeaders) > 0 {
			w.Header().Set("Access-Control-Allow-Headers", strings.Join(cfg.AllowedHeaders, ","))
		}
		if cfg.AllowCredentials {
			w.Header().Set("Access-Control-Allow-Credentials", "true")
		}
		if cfg.MaxAge > 0 {
			w.Header().Set("Access-Control-Max-Age", strconv.FormatInt(int64(cfg.MaxAge.Seconds()), 10))
		}

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}
