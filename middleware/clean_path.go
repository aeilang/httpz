package middleware

import (
	"net/http"
	"path"
)

// CleanPath middleware will clean out double slash mistakes from a user's request path.
// For example, if a user requests /users//1 or //users////1 will both be treated as: /users/1
func CleanPath(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		r.URL.RawPath = path.Clean(r.URL.RawPath)
		r.URL.Path = path.Clean(r.URL.Path)

		next.ServeHTTP(w, r)
	})
}
