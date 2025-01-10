package httpz

import "net/http"

type ErrHandler func(err error, w http.ResponseWriter)

func DefaultErrHandler(err error, w http.ResponseWriter) {
	if r, ok := w.(*ResponseWriter); ok && r.isCommited {
		return
	}

	rw := Unwrap(w)

	switch he := err.(type) {
	case *HTTPError:
		rw.JSON(he.StatusCode, Map{"msg": he.Msg})
	default:
		rw.JSON(http.StatusInternalServerError, Map{"msg": he.Error()})
	}
}

func Unwrap(w http.ResponseWriter) *ResponseWriter {
	rw, ok := w.(*ResponseWriter)
	if !ok {
		panic("Unwrap must be used in httpz")
	}

	return rw
}

// 自定义Error类型 参考了echo
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
