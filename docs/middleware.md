# Middleware

Middleware is supported by the HTTP transport.

## Adding Middleware

```go
import "github.com/mirkobrombin/go-module-router/v2/pkg/transport/http"

t := http.New()

t.Use(func(next stdhttp.Handler) stdhttp.Handler {
    return stdhttp.HandlerFunc(func(w stdhttp.ResponseWriter, r *stdhttp.Request) {
        start := time.Now()
        next.ServeHTTP(w, r)
        log.Printf("%s %s %v", r.Method, r.URL.Path, time.Since(start))
    })
})
```

## Middleware Order

Middleware is applied in LIFO order (last registered wraps outermost).

```go
t.Use(logging)   // outer
t.Use(auth)      // inner
// Request flow: logging -> auth -> handler -> auth -> logging
```

## Group-Scoped Middleware

Middleware added to a group applies only to that group's routes:

```go
api := t.Group("/api")
api.Use(authMiddleware)  // Only applies to /api/* routes
api.Register(&GetUser{})
```
