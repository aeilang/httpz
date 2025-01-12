package httpz

import (
	"log/slog"
	"net/http"
)

// centralized error handling function type.
type ErrHandlerFunc func(err error, w http.ResponseWriter)

// default centrailzed error handling function.
// only the *HTTPError will triger error response.
func DefaultErrHandlerFunc(err error, w http.ResponseWriter) {
	if he, ok := err.(*HTTPError); ok {
		rw := NewHelperRW(w)
		rw.JSON(he.StatusCode, Map{"msg": he.Msg})
	} else {
		slog.Error(err.Error())
	}
}

// The custom Error type is inspired by Echo.
type HTTPError struct {
	StatusCode int
	Msg        string
	errs       []error
}

func NewHTTPError(statusCode int, msg string, errs ...error) *HTTPError {
	return &HTTPError{
		StatusCode: statusCode,
		Msg:        msg,
		errs:       errs,
	}
}

func (e *HTTPError) Error() string {
	msg := e.Msg
	for _, err := range e.errs {
		msg += err.Error()
	}

	return msg
}
