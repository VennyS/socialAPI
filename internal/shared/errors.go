package shared

import (
	"errors"
	"net/http"
)

type HttpError struct {
	Err        error
	StatusCode int
}

func (h *HttpError) Error() string {
	return h.Err.Error()
}

func NewHttpError(err string, statusCode int) *HttpError {
	return &HttpError{Err: errors.New(err), StatusCode: statusCode}
}

var (
	InternalError      = NewHttpError("internal server error", http.StatusInternalServerError)
	InvalidCredentials = NewHttpError("invalid credentials", http.StatusUnauthorized)
)
