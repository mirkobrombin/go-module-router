# Dependency Injection

All transports share the same DI mechanism via `core.Container`.

## Providing Dependencies

```go
// Via transport
t := http.New()
t.Provide("DB", db)
t.Provide("UserService", userService)

// Or via router facade
r := router.New()
r.Provide("DB", db)  // Registers in all transports
```

## Injecting into Handlers

Dependencies are injected by matching field names:

```go
type GetUser struct {
    Meta core.Pattern `method:"GET" path:"/users/{id}"`

    // Injected by field name "DB"
    DB *sql.DB

    // Injected by field name "UserService"
    UserService UserService
}
```

The container matches `Provide("DB", ...)` to fields named `DB`.

## Type Safety

Injection is type-checked at runtime. If the provided value cannot be assigned to the field type, it will be skipped.

## Direct Container Usage

For advanced use cases:

```go
import "github.com/mirkobrombin/go-module-router/v2/pkg/core"

container := core.NewContainer()
container.Provide("Config", config)

// Inject into any struct
container.Inject(&myHandler)

// Retrieve a dependency
if dep, ok := container.Get("Config"); ok {
    cfg := dep.(*Config)
}
```
