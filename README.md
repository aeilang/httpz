![Logo](./_img/httpz.png "网站Logo")

> httpz v1.0.0 is release, its API is stable.

**[简体中文](https://github.com/aeilang/httpz/blob/main/README_CN.md)**

httpz is a lightweight library built on top of net/http version 1.22. It takes inspiration from Echo's centralized error handling and chi's adherence to the standard library. The problem it aims to solve is that while net/http version 1.22 enhances routing, its functionality is not as user-friendly as other frameworks like Echo and chi.

It functions more like a set of helper functions for net/http rather than a full-fledged web framework. Thanks to net/http handling most of the heavy lifting, httpz has minimal code.

It has the following features:

1. Centralized error handling


2. Convenient route grouping, where you can set middleware for each group or for individual routes.

3. Complete compatibility with the standard library.

# Quick Start

## 1.Installation

To install httpz, Go 1.22 or higher is required.

```sh
go get github.com/aeilang/httpz
```

## 2.Hello World

```go
import (
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
		
		// or you can write it by yourself.
		// hw.Header().Set("Content-Type", "text/plain; charset=UTF-8")
		// hw.WriteHeader(http.StatusOK)
		// hw.Write([]byte("hello httpz"))
		// return nil
	})
  
  // just like net/http's ServerMux
	http.ListenAndServe(":8080", mux)
}
```

> the middleware package is copied from chi/middleware. 

The complete example can be found in the [example](https://github.com/aeilang/httpz/blob/main//example/hello/main.go) directory

## 3.grouping:

```go
// group return a new *ServeMux base on path "/api/"
api := mux.Group("/api/")

// register GET /well route for api group.
// GET /api/well
api.Get("/well", func(w http.ResponseWriter, r *http.Request) error {	
	rw := httpz.NewHelperRW(w)
	return rw.JSON(http.StatusOK, httpz.Map{
		"data": "well well httpz",
	})
})
```

## 4.Centralized error handling

```go
// The parent mux of v2 is api,
// allowing you to group routes infinitely.
v2 := api.Group("/v2/")

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
```


You can customize the error handling function:

```go
// Create a new mux
mux := httpz.NewServeMux()

mux.ErrHandler = func(err error, w http.ResponseWriter) {
  // for example:
	http.Error(w, err.Error(), http.StatusInternalServerError)
}
```

The default error handling function is as follows:

```go
// default centrailzed error handling function.
// only the *HTTPError will triger sending error response.
func DefaultErrHandlerFunc(err error, w http.ResponseWriter) {
	if he, ok := err.(*HTTPError); ok {
		rw := NewHelperRW(w)
		rw.JSON(he.StatusCode, Map{"msg": he.Msg})
	} else {
		slog.Error(err.Error())
	}
}
```

## 5. Binding

You can also bind path parameters, query parameters, form parameters, and req.Body just like in Echo.

```go
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
```

# Feel free to contribute your code.

- test

- example

- middleware

- other
