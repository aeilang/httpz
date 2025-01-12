package main

import (
	"fmt"
	"net/http"

	"github.com/aeilang/httpz"
	"github.com/aeilang/httpz/middleware"
)

func main() {
	// Create a new mux
	mux := httpz.NewServeMux()

	// add logger middleware, it 's copy from chi/middleware
	mux.Use(middleware.Logger)

	// register a GET /hello route
	// GET /hello
	mux.Get("/hello", func(w http.ResponseWriter, r *http.Request) error {
		// rw is a helper responsewriter to send response
		rw := httpz.NewHelperRW(w)
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
		rw := httpz.NewHelperRW(w)
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

func API(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("before api")
		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

func V2(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
		fmt.Println("after api")
	}

	return http.HandlerFunc(fn)
}
