package main

import (
	"errors"
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
		hw := httpz.NewHelperRW(w)
		return hw.String(http.StatusOK, "hello httpz")

		// or you can write it by yourself
		// hw.Header().Set("Content-Type", "text/plain; charset=UTF-8")
		//		hw.WriteHeader(http.StatusOK)
		//		hw.Write([]byte("hello httpz"))
		//		return nil
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
		hw := httpz.NewHelperRW(w)
		return hw.JSON(http.StatusOK, httpz.Map{
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

	// Get /api/v2/httperr
	v2.Get("/httperr", func(w http.ResponseWriter, r *http.Request) error {

		// only *HTTPError will trigger the global error handling.
		// normal error just will just log the msg.
		return errors.New("some error")
	})

	type User struct {
		Name string `json:"name"`
	}

	// POST /api/v2/user
	v2.Post("/user", func(w http.ResponseWriter, r *http.Request) error {
		var u User
		if err := httpz.Bind(r, &u); err != nil {
			return err
		}

		w.WriteHeader(http.StatusCreated)
		return nil
	})

	// just like net/http's ServerMux
	http.ListenAndServe(":8080", mux)
}

// API middleware for test
func API(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("before api")
		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

// V2 middleware for test
func V2(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
		fmt.Println("after v2")
	}

	return http.HandlerFunc(fn)
}
