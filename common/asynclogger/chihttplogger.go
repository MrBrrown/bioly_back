package asynclogger

import (
	"fmt"
	"net/http"
	"time"
)

type Entry interface {
	Write(status int, elapsed time.Duration, method, path, reqID string)
	Panic(reqID string, v any, stack []byte)
}

type Formatter interface {
	NewEntry(r *http.Request) Entry
}

type defaultEntry struct {
	method string
	path   string
}

type DefaultFormatter struct{}

func (f *DefaultFormatter) NewEntry(r *http.Request) Entry {
	return &defaultEntry{
		method: r.Method,
		path:   r.URL.Path,
	}
}

func (e *defaultEntry) Write(status int, elapsed time.Duration, method, path, reqID string) {
	dur := fmt.Sprintf("%.3fms", float64(elapsed.Microseconds())/1000.0)

	if reqID != "" {
		line := fmt.Sprintf("[%s] %s %s -> %d in %s", reqID, e.method, e.path, status, dur)
		if status >= 400 {
			Error("%s", line)
		} else {
			Info("%s", line)
		}
		return
	}

	line := fmt.Sprintf("%s %s -> %d in %s", e.method, e.path, status, dur)
	if status >= 400 {
		Error("%s", line)
	} else {
		Info("%s", line)
	}
}

func (e *defaultEntry) Panic(reqID string, v any, stack []byte) {
	if reqID != "" {
		FatalError("PANIC [%s]: %v\nSTACK:\n%s", reqID, v, stack)
		return
	}
	FatalError("PANIC: %v\nSTACK:\n%s", v, stack)
}
