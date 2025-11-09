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

type LoginResponse struct {
	Access  string  `json:"access"`
	Refresh string  `json:"refresh"`
	User    UserDTO `json:"user"`
}

func (lr *LoginResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

type UserDTO struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
}
