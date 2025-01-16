package httpz

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewHTTPError(t *testing.T) {
	err := NewHTTPError(http.StatusBadRequest, "bad request")
	assert.Equal(t, http.StatusBadRequest, err.StatusCode)
	assert.Equal(t, "bad request", err.Msg)
	assert.Nil(t, err.Internal)
}

func TestHTTPError_SetInternal(t *testing.T) {
	internalErr := errors.New("internal error")
	err := NewHTTPError(http.StatusInternalServerError, "server error").SetInternal(internalErr)
	assert.Equal(t, internalErr, err.Internal)
}

func TestHTTPError_Error(t *testing.T) {
	err := NewHTTPError(http.StatusNotFound, "not found")
	assert.Equal(t, "code=404, message=not found", err.Error())

	internalErr := errors.New("internal error")
	err.SetInternal(internalErr)
	assert.Equal(t, "code=404, message=not found, internal=internal error", err.Error())
}

func TestHTTPError_Unwrap(t *testing.T) {
	internalErr := errors.New("internal error")
	err := NewHTTPError(http.StatusInternalServerError, "server error").SetInternal(internalErr)
	assert.Equal(t, internalErr, err.Unwrap())
}

func TestDefaultErrHandlerFunc(t *testing.T) {
	rec := httptest.NewRecorder()
	err := NewHTTPError(http.StatusBadRequest, "bad request")
	DefaultErrHandlerFunc(err, rec)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.JSONEq(t, `{"msg":"bad request"}`, rec.Body.String())
}
