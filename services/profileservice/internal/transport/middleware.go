package transport

import (
	"bioly/common/asynclogger"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
)

type ChiAdapter struct {
	Formatter asynclogger.Formatter
}

func (a *ChiAdapter) NewLogEntry(r *http.Request) middleware.LogEntry {
	entry := a.Formatter.NewEntry(r)
	return &chiEntry{r: r, entry: entry}
}

type chiEntry struct {
	r     *http.Request
	entry asynclogger.Entry
}

func (e *chiEntry) Write(status, _ int, _ http.Header, elapsed time.Duration, _ any) {
	reqID := middleware.GetReqID(e.r.Context())
	e.entry.Write(status, elapsed, e.r.Method, e.r.URL.Path, reqID)
}

func (e *chiEntry) Panic(v any, stack []byte) {
	reqID := middleware.GetReqID(e.r.Context())
	e.entry.Panic(reqID, v, stack)
}
