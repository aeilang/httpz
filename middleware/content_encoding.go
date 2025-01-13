// Copyright (c) 2015-present Peter Kieltyka (https://github.com/pkieltyka), Google Inc.
// Declaration: The middleware package is copied from chi/v5/middleware. Source:
// https://github.com/go-chi/chi/tree/master/middleware

package middleware

import (
	"net/http"
	"strings"
)

// AllowContentEncoding enforces a whitelist of request Content-Encoding otherwise responds
// with a 415 Unsupported Media Type status.
func AllowContentEncoding(contentEncoding ...string) func(next http.Handler) http.Handler {
	allowedEncodings := make(map[string]struct{}, len(contentEncoding))
	for _, encoding := range contentEncoding {
		allowedEncodings[strings.TrimSpace(strings.ToLower(encoding))] = struct{}{}
	}
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			requestEncodings := r.Header["Content-Encoding"]
			// skip check for empty content body or no Content-Encoding
			if r.ContentLength == 0 {
				next.ServeHTTP(w, r)
				return
			}
			// All encodings in the request must be allowed
			for _, encoding := range requestEncodings {
				if _, ok := allowedEncodings[strings.TrimSpace(strings.ToLower(encoding))]; !ok {
					w.WriteHeader(http.StatusUnsupportedMediaType)
					return
				}
			}
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}
