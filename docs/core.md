# Core Concepts

Go Module Router uses a transport-agnostic core that can be used with different dispatching mechanisms.

## Handler Interface

All handlers implement the same interface:

```go
type Handler interface {
    Handle(ctx context.Context) (any, error)
}
```

## Pattern Type

The `Pattern` struct is used in tags to declare routing information. Different transports read different tags:

```go
type MyHandler struct {
    // HTTP transport reads 'method' and 'path'
    Meta core.Pattern `method:"GET" path:"/users/{id}"`

    // Action transport reads 'action' and 'keys'
    Meta core.Pattern `action:"file.save" keys:"ctrl+s"`
}
```

## Container (DI)

The `Container` manages dependency injection:

```go
container := core.NewContainer()
container.Provide("DB", db)
container.Provide("Logger", logger)

// Inject into a struct
container.Inject(&myHandler)
```

## Binder

The `Binder` maps external data to struct fields:

```go
binder := core.NewBinder()
binder.AddSource("query", func(key string) string {
    return request.URL.Query().Get(key)
})
binder.Bind(&handler)
```

## Architecture

```
+----------------+
|    Handler     |  <- Your struct implementing Handle()
+----------------+
        |
        v
+----------------+
|     Core       |  <- DI Container + Binder
+----------------+
        |
        v
+----------------+
|   Transport    |  <- HTTP, Action, CLI, etc.
+----------------+
```
