> Note that httpz is still unstable.

httpz is a lightweight library built on top of net/http version 1.22. It takes inspiration from Echo's centralized error handling and chi's adherence to the standard library. The problem it aims to solve is that while net/http version 1.22 enhances routing, its functionality is not as user-friendly as other frameworks like Echo and chi.

It functions more like a set of helper functions for net/http rather than a full-fledged web framework. Thanks to net/http handling most of the heavy lifting, httpz has minimal code.

It has the following features:

1. Centralized error handling


2. Convenient route grouping, where you can set middleware for each group or for individual routes.

3. Complete compatibility with the standard library.

### Quick Start

#### Installation

To install httpz, Go 1.22 or higher is required.

```sh
go get github.com/aeilang/httpz
```

#### Hello World

```go
import (
	"log/slog"
	"net/http"
	"os"

	"github.com/aeilang/httpz"
	"github.com/aeilang/httpz/mws"
)

func main() {
	// Create a new mux
	mux := httpz.NewServeMux()

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
  
  // just like net/http's ServerMux
	http.ListenAndServe(":8080", mux)
}
```

> you don't use slog, you can easily create a custom logging middleware.

The complete example can be found in the [example](https://github.com/aeilang/httpz/blob/main//example/hello/main.go) directory

#### grouping:

```go
api := mux.Group("/api/")

// register GET /well route for api group.
// GET /api/well
api.Get("/well", func(w http.ResponseWriter, r *http.Request) error {
	rw := httpz.Unwrap(w)
	return rw.JSON(http.StatusOK, httpz.Map{
		"data": "well well httpz",
	})
})
```

#### Centralized error handling

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
```


You can customize the error handling function:

```go
// Create a new mux
mux := httpz.NewServeMux()

mux.ErrHandler = func(err error, w http.ResponseWriter) {
  // ...
	http.Error(w, err.Error(), http.StatusInternalServerError)
}
```

The default error handling function is as follows:

```go
func DefaultErrHandler(err error, w http.ResponseWriter) {
	rw := Unwrap(w)

	if rw.isCommited || err == nil {
		return
	}

	switch he := err.(type) {
	case *HTTPError:
		rw.JSON(he.StatusCode, Map{"msg": he.Msg})
	default:
		rw.JSON(http.StatusInternalServerError, Map{"msg": he.Error()})
	}
}
```

#### Feel free to contribute your code.

- test

- example

- middleware

- other
