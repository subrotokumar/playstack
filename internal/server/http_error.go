package server

import (
	"net/http"

	"gofr.dev/pkg/gofr/logging"
)

type HttpError interface {
	Error() string
	StatusCode() int
	LogLevel() logging.Level
}

type httpError struct {
	msg    string
	status int
	level  logging.Level
	err    error
}

func (e *httpError) Error() string {
	if e.err != nil {
		return e.msg + ": " + e.err.Error()
	}
	return e.msg
}

func (e *httpError) StatusCode() int {
	return e.status
}

func (e *httpError) LogLevel() logging.Level {
	return e.level
}

func (e *httpError) Unwrap() error {
	return e.err
}

func NewHttpError(status int, msg string) HttpError {
	return &httpError{
		msg:    msg,
		status: status,
		level:  logging.ERROR,
	}
}

func WrapHttpError(status int, msg string, err error) HttpError {
	return &httpError{
		msg:    msg,
		status: status,
		level:  logging.ERROR,
		err:    err,
	}
}

func BadRequest(msg string) HttpError {
	return &httpError{
		msg:    msg,
		status: http.StatusBadRequest,
		level:  logging.WARN,
	}
}

func Unauthorized(msg string) HttpError {
	return &httpError{
		msg:    msg,
		status: http.StatusUnauthorized,
		level:  logging.WARN,
	}
}

func Forbidden(msg string) HttpError {
	return &httpError{
		msg:    msg,
		status: http.StatusForbidden,
		level:  logging.WARN,
	}
}

func Internal(err error) HttpError {
	return &httpError{
		msg:    "internal server error",
		status: http.StatusInternalServerError,
		level:  logging.ERROR,
		err:    err,
	}
}
