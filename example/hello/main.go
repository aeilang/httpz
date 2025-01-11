package main

import (
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/aeilang/httpz"
	"github.com/aeilang/httpz/mws"
)

func main() {
	// Create a new mux
	mux := httpz.NewServeMux()
	
	mux.ErrHandler = func(err error, w http.ResponseWriter) {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	// logger use slog package
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	// register a logger middleware
	mux.Use(mws.Logger(logger))

	// register a GET /hello route
	// GET /hello
	mux.Get("/hello", func(w http.ResponseWriter, r *http.Request) error {
		// Unwrap w to get its underlying implementation.
		// rw includes helper methods for sending responses.
		rw := httpz.Unwrap(w)
		return rw.String(http.StatusOK, "hello httpz")
	})

	// Group routes based on /api/, making sure to include the trailing slash /.
	// This is required by the standard library syntax,
	// which handles all requests starting with /api.
	api := mux.Group("/api/")
	// use API middleware for this api group. just for testing the ability。
	api.Use(API)

	// register GET /well route for api group.
	// GET /api/well
	api.Get("/well", func(w http.ResponseWriter, r *http.Request) error {
		rw := httpz.Unwrap(w)
		return rw.JSON(http.StatusOK, httpz.Map{
			"data": "well well httpz",
		})
	})

	// The parent mux of v2 is api,
	// allowing you to group routes infinitely.
	v2 := api.Group("/v2/")

	// use V2 middleware for this group, just for testing the ability.
	v2.Use(V2)

	// GET /api/v2/hello
	v2.Get("/hello", func(w http.ResponseWriter, r *http.Request) error {
		// centralized error handling in tests
		return httpz.NewHTTPError(http.StatusBadRequest, "bad reqeust")
	})

	// testing path parameters and centrialzed error handling.
	// GET /api/v2/well/randomID
	v2.Get("/well/{id}", func(w http.ResponseWriter, r *http.Request) error {
		id := r.PathValue("id")

		return httpz.NewHTTPError(http.StatusBadRequest, id)
	})

	// just like net/http's ServerMux
	http.ListenAndServe(":8080", mux)
}

// middleware for example
func API(next httpz.HandlerFunc) httpz.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		log.Println("apiapi")
		return next(w, r)
	}
}

// middleware for example
func V2(next httpz.HandlerFunc) httpz.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		log.Println("v2v2")
		return next(w, r)
	}
}
