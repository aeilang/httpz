// Copyright (c) 2015-present Peter Kieltyka (https://github.com/pkieltyka), Google Inc.
// Declaration: The middleware package is copied from chi/v5/middleware. Source:
// https://github.com/go-chi/chi/tree/master/middleware

package middleware

import (
	"context"
	"net/http"
)

// WithValue is a middleware that sets a given key/value in a context chain.
func WithValue(key, val interface{}) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			r = r.WithContext(context.WithValue(r.Context(), key, val))
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}
