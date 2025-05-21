# Dependency Injection (Zero Code-Gen)

DI in `go-module-router` is **factory-based**:

```
┌──────────┐                ┌────────────┐
│ Factory  │ ----requires→  │   Deps     │
└──────────┘                └────────────┘

```

*   **Repositories** receive a map of _already built_ repositories.
*   **Services** receive a map of repositories (and may build on each other).
*   **Handlers / Middleware** receive a map of services.

Because everything is ordinary Go code you can:

*   Initialise third-party clients inside a factory.
*   Return an interface instead of a concrete struct.
*   Decide at runtime to _skip_ auto-wiring by returning `nil` and injecting
    your own instance later.

### Circular dependencies

Factories run **in registration order**.  
Avoid circular service graphs—or break the cycle by injecting an already built
instance in `services map[string]any` when you call `router.New`.


