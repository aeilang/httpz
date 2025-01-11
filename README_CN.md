
> 请注意，httpz 仍然不稳定。

背景： net/http 1.22 虽然增强了路由功能，但使用体验却不如 Echo 和 chi 等框架。

httpz 是一个基于 net/http 1.22 版本构建的轻量级库，借鉴了 Echo 的集中式错误处理以及 chi 的小和轻量。

httpz更像是 net/http 1.22 的一组helper函数，而非一个完整的 Web 框架。由于核心工作仍由 net/http 执行，httpz 的代码量非常少。

它具有以下特性：

- 集中式错误处理

- 便捷的路由分组，可以为每个分组设置中间件，也可以为单个路由设置中间件。

- 完全兼容标准库。

### 快速开始

#### 安装

要安装 httpz，需要 Go 1.22+

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

> 如果你不使用slog, 可以自定义Logger中间件

完整的hello world例子在 [example](https://github.com/aeilang/httpz/blob/main//example/hello/main.go) 目录

#### 分组:

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

#### 集中式错误处理

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


自定义全局错误处理函数：

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

#### 欢迎贡献你的代码

- test

- example

- middleware

- other
