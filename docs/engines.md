# HTTP Engines

`Engine` is a **minimal frontend** over any HTTP framework, allowing you
to choose the one that fits your needs best or even write your own.

```go
type Engine interface {
    Handle(method, path string, h http.Handler)
    Use(mw Middleware)
    Group(prefix string) Engine

    Serve(addr string) error
    Shutdown(ctx context.Context) error

    Unwrap() any
}
```

### Built-in adapters

| Package                 | Notes                                                   |
| ----------------------- | ------------------------------------------------------- |
| `httpdrv.NewFastHTTP()` | Uses `github.com/valyala/fasthttp` + `fasthttp/router`. |
| `httpdrv.NewStdHTTP()`  | Thin wrapper around `net/http`.                         |

### Writing your own

Implement the three routing methods and the two lifecycle methods, return your
native instance from `Unwrap()` if callers need direct access.

