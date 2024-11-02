package httperror

import (
	"fmt"
	"net/http"
)

type (
	Error interface {
		error
		HTTPStatusCode() int
	}

	httpError struct {
		StatusCode  int    `json:"-"`
		OK          bool   `json:"ok"`
		Description string `json:"description"`
	}
)

func (e *httpError) HTTPStatusCode() int {
	return e.StatusCode
}

func (e *httpError) Error() string {
	return e.Description
}

func New(status int, msg string, args ...any) *httpError {
	if len(args) > 0 {
		msg = fmt.Sprintf(msg, args...)
	}
	return &httpError{
		OK:          false,
		StatusCode:  status,
		Description: msg,
	}
}

func NotFound(msg string, args ...any) Error {
	return New(http.StatusNotFound, msg, args...)
}

func Unauthorized(msg string, args ...any) Error {
	return New(http.StatusUnauthorized, msg, args...)
}

func Forbidden(msg string, args ...any) Error {
	return New(http.StatusForbidden, msg, args...)
}

func BadRequest(code, msg string, args ...any) Error {
	return New(http.StatusBadRequest, msg, args...)
}
