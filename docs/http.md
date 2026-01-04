# HTTP Transport

The HTTP transport handles web requests and APIs.

## Creating the Transport

```go
import "github.com/mirkobrombin/go-module-router/v2/pkg/transport/http"

t := http.New()
```

Or via the router facade:

```go
import "github.com/mirkobrombin/go-module-router/v2/pkg/router"

r := router.New()
// r.HTTP is the HTTP transport
```

## Defining Endpoints

Use `core.Pattern` with `method` and `path` tags:

```go
type GetUser struct {
    Meta core.Pattern `method:"GET" path:"/users/{id}"`

    // Path parameter
    ID string `path:"id"`

    // Query parameter
    Details bool `query:"details" default:"false"`

    // Header
    Token string `header:"Authorization"`

    // JSON body
    Body CreateUserRequest `body:"json"`

    // Dependencies
    DB *sql.DB
}

func (e *GetUser) Handle(ctx context.Context) (any, error) {
    return User{ID: e.ID}, nil
}
```

## Registering Endpoints

```go
t.Provide("DB", db)
t.Register(&GetUser{})
t.Register(&CreateUser{})
```

## Middleware

```go
t.Use(func(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        log.Println("Request:", r.URL.Path)
        next.ServeHTTP(w, r)
    })
})
```

## Route Groups

```go
api := t.Group("/api/v1")
api.Register(&GetUser{})  // -> GET /api/v1/users/{id}
```

## Starting the Server

```go
t.Listen(":8080")
```

## Accessing the Raw Mux

```go
mux := t.Mux()
mux.HandleFunc("/health", healthHandler)
```
