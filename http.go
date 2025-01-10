//go:build go1.22

package httpz

import (
	"fmt"
	"net/http"
	"strings"
	"sync"
)

type HandlerFunc func(w http.ResponseWriter, r *http.Request) error

type MiddlewareFunc func(next HandlerFunc) HandlerFunc

type ServeMux struct {
	http.ServeMux
	ErrHandler ErrHandler
	mws        []MiddlewareFunc
	groups     map[string]*ServeMux
	isMaster   bool
	once       sync.Once
}

func NewServeMux() *ServeMux {
	return &ServeMux{
		ServeMux:   http.ServeMux{},
		groups:     make(map[string]*ServeMux),
		ErrHandler: DefaultErrHandler,
		isMaster:   true,
	}
}

func (sm *ServeMux) HandleFunc(pattern string, h HandlerFunc) {
	sm.ServeMux.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		err := h(w, r)
		sm.ErrHandler(err, w)
	})
}

func (sm *ServeMux) addToMux() {
	for prefix, mux := range sm.groups {
		mux.addToMux()
		pre := strings.TrimSuffix(prefix, "/")
		sm.Handle(prefix, http.StripPrefix(pre, mux))
	}
}

func (sm *ServeMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if sm.isMaster {
		sm.once.Do(func() {
			sm.addToMux()
		})
	}

	w = &ResponseWriter{ResponseWriter: w}

	h := http.Handler(&sm.ServeMux)

	hf := toFunc(h)

	for i := len(sm.mws) - 1; i >= 0; i-- {
		hf = sm.mws[i](hf)
	}

	err := hf(w, r)
	sm.ErrHandler(err, w)
}

func (sm *ServeMux) Group(prefix string) *ServeMux {
	mux := NewServeMux()
	mux.isMaster = false
	if _, existed := sm.groups[prefix]; existed {
		panic(fmt.Sprintf("prefix %s already existed", prefix))
	}

	sm.groups[prefix] = mux
	return mux
}

func (sm *ServeMux) Get(path string, h HandlerFunc, m ...MiddlewareFunc) {
	h = use(h, m...)
	sm.HandleFunc(fmt.Sprintf("%s %s", http.MethodGet, path), h)
}

func (sm *ServeMux) Head(path string, h HandlerFunc, m ...MiddlewareFunc) {
	h = use(h, m...)
	sm.HandleFunc(fmt.Sprintf("%s %s", http.MethodHead, path), h)
}

func (sm *ServeMux) Post(path string, h HandlerFunc, m ...MiddlewareFunc) {
	h = use(h, m...)
	sm.HandleFunc(fmt.Sprintf("%s %s", http.MethodPost, path), h)
}

func (sm *ServeMux) Put(path string, h HandlerFunc, m ...MiddlewareFunc) {
	h = use(h, m...)
	sm.HandleFunc(fmt.Sprintf("%s %s", http.MethodGet, path), h)
}

func (sm *ServeMux) Patch(path string, h HandlerFunc, m ...MiddlewareFunc) {
	h = use(h, m...)
	sm.HandleFunc(fmt.Sprintf("%s %s", http.MethodPatch, path), h)
}

func (sm *ServeMux) Delete(path string, h HandlerFunc, m ...MiddlewareFunc) {
	h = use(h, m...)
	sm.HandleFunc(fmt.Sprintf("%s %s", http.MethodDelete, path), h)
}

func (sm *ServeMux) Connect(path string, h HandlerFunc, m ...MiddlewareFunc) {
	h = use(h, m...)
	sm.HandleFunc(fmt.Sprintf("%s %s", http.MethodConnect, path), h)
}

func (sm *ServeMux) Options(path string, h HandlerFunc, m ...MiddlewareFunc) {
	h = use(h, m...)
	sm.HandleFunc(fmt.Sprintf("%s %s", http.MethodOptions, path), h)
}

func (sm *ServeMux) Trace(path string, h HandlerFunc, m ...MiddlewareFunc) {
	h = use(h, m...)
	sm.HandleFunc(fmt.Sprintf("%s %s", http.MethodTrace, path), h)
}

func use(h HandlerFunc, m ...MiddlewareFunc) HandlerFunc {
	for i := len(m) - 1; i >= 0; i-- {
		h = m[i](h)
	}

	return h
}

func (sm *ServeMux) Use(m ...MiddlewareFunc) {
	sm.mws = append(sm.mws, m...)
}

func toFunc(h http.Handler) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		h.ServeHTTP(w, r)
		return nil
	}
}
