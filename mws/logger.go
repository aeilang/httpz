package mws

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"time"

	"github.com/aeilang/httpz"
)

func Logger(logger *slog.Logger) httpz.MiddlewareFunc {
	return func(next httpz.HandlerFunc) httpz.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) error {
			now := time.Now()

			err := next(w, r)

			rw := httpz.Unwrap(w)

			if rw.StatusCode() >= 200 && rw.StatusCode() < 400 {
				logger.Info(
					r.URL.Path,
					slog.String("method", r.Method),
					slog.String("status", fmt.Sprintf("%d %s", rw.StatusCode(), http.StatusText(rw.StatusCode()))),
					slog.Duration("elapsed", time.Since(now)),
				)
			} else {
				logger.Error(
					r.URL.Path,
					slog.String("method", r.Method),
					slog.String("status", fmt.Sprintf("%d %s", rw.StatusCode(), http.StatusText(rw.StatusCode()))),
					slog.Duration("elapsed", time.Since(now)),
				)
			}

			return err
		}
	}
}

func API(next httpz.HandlerFunc) httpz.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		log.Println("apiapi")
		return next(w, r)
	}
}

func V2(next httpz.HandlerFunc) httpz.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		log.Println("v2v2")
		return next(w, r)
	}
}
