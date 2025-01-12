
> httpz v1.0.0 版本已经发布，它的API已经稳定

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
  
  // just like net/http's ServerMux
	http.ListenAndServe(":8080", mux)
}
```

> middleware 包来自chi/middleware

完整的hello world例子在 [example](https://github.com/aeilang/httpz/blob/main//example/hello/main.go) 目录

#### 分组:

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

// GET /api/v2/well/randomID
v2.Get("/well/{id}", func(w http.ResponseWriter, r *http.Request) error {
	id := r.PathValue("id")
	
	// the default error handler just trigered by *HTTPError
	// another error will just be logged,not sending response.
	return errors.New("nomal error")
})
```


自定义全局错误处理函数：

```go
// Create a new mux
mux := httpz.NewServeMux()

mux.ErrHandler = func(err error, w http.ResponseWriter) {
  // for example:
	http.Error(w, err.Error(), http.StatusInternalServerError)
}
```

默认的全局错误处理函数如下：

```go
// default centrailzed error handling function.
// only the *HTTPError will triger error response.
func DefaultErrHandlerFunc(err error, w http.ResponseWriter) {
	if he, ok := err.(*HTTPError); ok {
		rw := NewHelperRW(w)
		rw.JSON(he.StatusCode, Map{"msg": he.Msg})
	} else {
		slog.Error(err.Error())
	}
}
```

#### 欢迎贡献你的代码

- test

- example

- middleware

- other
