package httpz

import (
	"net/http"
)

// centralized error handling function type.
type ErrHandler func(err error, w http.ResponseWriter)

// default centrailzed error handling function
func DefaultErrHandler(err error, w http.ResponseWriter) {
	rw := Unwrap(w)

	if rw.isCommited || err == nil {
		return
	}

	switch he := err.(type) {
	case *HTTPError:
		rw.JSON(he.StatusCode, Map{"msg": he.Msg})
	default:
		rw.JSON(http.StatusInternalServerError, Map{"msg": he.Error()})
	}
}

// helper function to get underline *ResponseWriter
func Unwrap(w http.ResponseWriter) *ResponseWriter {
	rw, ok := w.(*ResponseWriter)
	if !ok {
		panic("Unwrap must be used in httpz")
	}

	return rw
}

// The custom Error type is inspired by Echo.
type HTTPError struct {
	StatusCode int
	Msg        string
}

func NewHTTPError(statusCode int, msg string) *HTTPError {
	return &HTTPError{
		StatusCode: statusCode,
		Msg:        msg,
	}
}

func (e *HTTPError) Error() string {
	return e.Msg
}
