//go:build go1.22

package httpz

import (
	"fmt"
	"net/http"
	"strings"
)

// HandlerFunc defines the function signature for a handler.
// It returns an error, which is used for centralized error handling.
type HandlerFunc func(w http.ResponseWriter, r *http.Request) error

// MiddlewareFunc defines the function signature for middleware.
// It wraps an http.Handler and returns a new http.Handler.
type MiddlewareFunc func(next http.Handler) http.Handler

// RouteMiddlewareFunc defines the function signature for route-specific middleware.
// It wraps a HandlerFunc and returns a new HandlerFunc.
type RouteMiddlewareFunc func(next HandlerFunc) HandlerFunc

// ServeMux embeds http.ServeMux and provides additional features like error handling and middleware.
type ServeMux struct {
	http.ServeMux
	ErrHandlerFunc ErrHandlerFunc   // Function for centralized error handling
	mws            []MiddlewareFunc // List of middleware functions
}

// NewServeMux returns a new instance of ServeMux with default settings.
func NewServeMux() *ServeMux {
	return &ServeMux{
		ServeMux:       http.ServeMux{},
		ErrHandlerFunc: DefaultErrHandlerFunc,
	}
}

// HandleFunc registers a new route with a pattern and a handler function.
// The handler function can return an error for centralized error handling.
func (sm *ServeMux) HandleFunc(pattern string, h HandlerFunc) {
	sm.ServeMux.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		err := h(w, r)

		if err != nil {
			sm.ErrHandlerFunc(err, w)
		}
	})
}

// ServeHTTP processes HTTP requests using the registered handlers and middleware.
func (sm *ServeMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h := http.Handler(&sm.ServeMux)

	for i := len(sm.mws) - 1; i >= 0; i-- {
		h = sm.mws[i](h)
	}

	h.ServeHTTP(w, r)
}

// Group creates a new ServeMux for a specific URL prefix, allowing for route grouping.
func (sm *ServeMux) Group(prefix string) *ServeMux {
	if len(prefix) == 0 {
		panic("len(prefix) must greater than 0")
	}

	if prefix[len(prefix)-1] != '/' {
		panic("the last char ine prefix must b /")
	}

	mux := &ServeMux{
		ServeMux:       http.ServeMux{},
		ErrHandlerFunc: sm.ErrHandlerFunc,
	}

	pre := strings.TrimSuffix(prefix, "/")
	sm.Handle(prefix, http.StripPrefix(pre, mux))

	return mux
}

// Get registers a new GET route with optional route-specific middleware.
func (sm *ServeMux) Get(path string, h HandlerFunc, m ...RouteMiddlewareFunc) {
	h = use(h, m...)
	sm.HandleFunc(fmt.Sprintf("%s %s", http.MethodGet, path), h)
}

// Head registers a new HEAD route with optional route-specific middleware.
func (sm *ServeMux) Head(path string, h HandlerFunc, m ...RouteMiddlewareFunc) {
	h = use(h, m...)
	sm.HandleFunc(fmt.Sprintf("%s %s", http.MethodHead, path), h)
}

// Post registers a new POST route with optional route-specific middleware.
func (sm *ServeMux) Post(path string, h HandlerFunc, m ...RouteMiddlewareFunc) {
	h = use(h, m...)
	sm.HandleFunc(fmt.Sprintf("%s %s", http.MethodPost, path), h)
}

// Put registers a new PUT route with optional route-specific middleware.
func (sm *ServeMux) Put(path string, h HandlerFunc, m ...RouteMiddlewareFunc) {
	h = use(h, m...)
	sm.HandleFunc(fmt.Sprintf("%s %s", http.MethodGet, path), h)
}

// Patch registers a new PATCH route with optional route-specific middleware.
func (sm *ServeMux) Patch(path string, h HandlerFunc, m ...RouteMiddlewareFunc) {
	h = use(h, m...)
	sm.HandleFunc(fmt.Sprintf("%s %s", http.MethodPatch, path), h)
}

// Delete registers a new DELETE route with optional route-specific middleware.
func (sm *ServeMux) Delete(path string, h HandlerFunc, m ...RouteMiddlewareFunc) {
	h = use(h, m...)
	sm.HandleFunc(fmt.Sprintf("%s %s", http.MethodDelete, path), h)
}

// Connect registers a new CONNECT route with optional route-specific middleware.
func (sm *ServeMux) Connect(path string, h HandlerFunc, m ...RouteMiddlewareFunc) {
	h = use(h, m...)
	sm.HandleFunc(fmt.Sprintf("%s %s", http.MethodConnect, path), h)
}

// Options registers a new OPTIONS route with optional route-specific middleware.
func (sm *ServeMux) Options(path string, h HandlerFunc, m ...RouteMiddlewareFunc) {
	h = use(h, m...)
	sm.HandleFunc(fmt.Sprintf("%s %s", http.MethodOptions, path), h)
}

// Trace registers a new TRACE route with optional route-specific middleware.
func (sm *ServeMux) Trace(path string, h HandlerFunc, m ...RouteMiddlewareFunc) {
	h = use(h, m...)
	sm.HandleFunc(fmt.Sprintf("%s %s", http.MethodTrace, path), h)
}

// use applies route-specific middleware to a handler function.
func use(h HandlerFunc, m ...RouteMiddlewareFunc) HandlerFunc {
	for i := len(m) - 1; i >= 0; i-- {
		h = m[i](h)
	}

	return h
}

// Use adds middleware to the ServeMux, which will be applied to all routes.
func (sm *ServeMux) Use(m ...MiddlewareFunc) {
	sm.mws = append(sm.mws, m...)
}

// Adator converts a standard http.HandlerFunc to a HandlerFunc that returns an error.
func Adator(fn func(w http.ResponseWriter, r *http.Request)) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		fn(w, r)
		return nil
	}
}
