package httpz

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestServeMux(t *testing.T) {
	mux := NewServeMux()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
	assert.Equal(t, "404 page not found\n", rec.Body.String())
}

func TestServeMux_HandleFunc(t *testing.T) {
	mux := NewServeMux()
	mux.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) error {
		w.Write([]byte("test"))
		return nil
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "test", rec.Body.String())
}

func TestMuxGet(t *testing.T) {
	mux := NewServeMux()

	mux.Get("/", func(w http.ResponseWriter, r *http.Request) error {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("hello world"))
		return nil
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "hello world", rec.Body.String())
}

func TestMuxGetWithPathParameter(t *testing.T) {
	mux := NewServeMux()
	mux.Get("/{id}", func(w http.ResponseWriter, r *http.Request) error {
		id := r.PathValue("id")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, id)
		return nil
	})

	req := httptest.NewRequest(http.MethodGet, "/1", nil)
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "1", rec.Body.String())
}

func TestHandlerMiddleware(t *testing.T) {
	mux := NewServeMux()
	buf := new(bytes.Buffer)

	mux.Use(func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			buf.WriteString("1")
			next.ServeHTTP(w, r)
			buf.WriteString("-1")
		}
		return http.HandlerFunc(fn)
	})

	mux.Use(func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			buf.WriteString("2")
			next.ServeHTTP(w, r)
			buf.WriteString("-2")
		}
		return http.HandlerFunc(fn)
	})

	mux.Use(func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			buf.WriteString("3")
			next.ServeHTTP(w, r)
			buf.WriteString("-3")
		}
		return http.HandlerFunc(fn)
	})

	mux.Use(func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			buf.WriteString("4")
			next.ServeHTTP(w, r)
			buf.WriteString("-4")
		}
		return http.HandlerFunc(fn)
	}, func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			buf.WriteString("5")
			next.ServeHTTP(w, r)
			buf.WriteString("-5")
		}
		return http.HandlerFunc(fn)
	})

	mux.Get("/", func(w http.ResponseWriter, r *http.Request) error {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
		return nil
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "OK", rec.Body.String())
	assert.Equal(t, "12345-5-4-3-2-1", buf.String())
}

func TestRouteMiddleware(t *testing.T) {
	buf := new(bytes.Buffer)
	mux := NewServeMux()

	mux.Get("/", func(w http.ResponseWriter, r *http.Request) error {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "OK")
		return nil
	}, func(next HandlerFunc) HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) error {
			buf.WriteString("1")
			err := next(w, r)
			buf.WriteString("-1")
			return err
		}
	}, func(next HandlerFunc) HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) error {
			buf.WriteString("2")
			err := next(w, r)
			buf.WriteString("-2")
			return err
		}
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "OK", rec.Body.String())
	assert.Equal(t, "12-2-1", buf.String())
}

func TestErrorHandler(t *testing.T) {
	mux := NewServeMux()

	mux.Get("/", func(w http.ResponseWriter, r *http.Request) error {
		return NewHTTPError(http.StatusBadRequest, "bad")
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.JSONEq(t, `{"msg":"bad"}`, rec.Body.String())
}

func TestGroup(t *testing.T) {
	mux := NewServeMux()

	v1 := mux.Group("/v1/")
	auth := v1.Group("/auth/")

	auth.Get("/user", func(w http.ResponseWriter, r *http.Request) error {
		fmt.Fprint(w, "user")
		return nil
	})

	req := httptest.NewRequest(http.MethodGet, "/v1/auth/user", nil)
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "user", rec.Body.String())
}

func TestServeMux_Adator(t *testing.T) {
	mux := NewServeMux()
	mux.HandleFunc("/adator", Adator(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("adator"))
	}))

	req := httptest.NewRequest(http.MethodGet, "/adator", nil)
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "adator", rec.Body.String())
}
