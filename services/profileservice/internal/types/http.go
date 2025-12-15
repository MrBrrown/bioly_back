package types

import (
	"net/http"

	"github.com/go-chi/render"
)

type ErrResponse struct {
	HTTPStatusCode int    `json:"-"`
	Error          string `json:"error"`
	AppCode        int64  `json:"app_code,omitempty"`
}

func (e *ErrResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

func ErrInvalidRequest(status int, err error) *ErrResponse {
	return &ErrResponse{
		HTTPStatusCode: status,
		Error:          err.Error(),
	}
}

type ProfileResponse struct {
	Username string `json:"username"`
	Page     JSONB  `json:"page"`
}

func (pr *ProfileResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}
