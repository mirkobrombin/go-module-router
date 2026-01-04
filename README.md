# Go Module Router

A **declarative**, transport-agnostic module router for Go.

Define handlers as structs with tags. The router handles binding, dependency injection, and dispatching. Works for web backends, GUI apps, CLI tools, and more.

## Features

- **Declarative Handlers:** Define handlers using struct tags.
- **Transport Agnostic:** Use the same pattern for HTTP, GUI actions, CLI commands.
- **Auto-Binding:** Parameters are automatically bound to struct fields.
- **Dependency Injection:** Services are injected by field name.
- **Middleware Support:** Standard middleware for HTTP transport.

## Installation

```bash
go get github.com/mirkobrombin/go-module-router/v2
```

## Transports

### HTTP Transport

For web backends and APIs.

```go
package main

import (
    "context"
    "log/slog"

    "github.com/mirkobrombin/go-module-router/v2/pkg/core"
    "github.com/mirkobrombin/go-module-router/v2/pkg/logger"
    "github.com/mirkobrombin/go-module-router/v2/pkg/router"
)

type GetUser struct {
    Meta core.Pattern `method:"GET" path:"/users/{id}"`
    ID   string       `path:"id"`
    DB   *sql.DB
}

func (e *GetUser) Handle(ctx context.Context) (any, error) {
    return User{ID: e.ID}, nil
}

func main() {
    r := router.New()
    r.SetLogger(logger.NewSlog(slog.Default()))
    r.Provide("DB", db)
    r.Register(&GetUser{})
    r.Listen(":8080")
}
```

### Action Transport

For GUI apps, TUI editors, or any event-driven application.

```go
package main

import (
    "context"

    "github.com/mirkobrombin/go-module-router/v2/pkg/core"
    "github.com/mirkobrombin/go-module-router/v2/pkg/transport/action"
)

type SaveAction struct {
    Meta     core.Pattern `action:"file.save" keys:"ctrl+s"`
    Document *Document
}

func (a *SaveAction) Handle(ctx context.Context) (any, error) {
    return a.Document.Save()
}

func main() {
    t := action.New()
    t.Provide("Document", doc)
    t.Register(&SaveAction{})

    // Dispatch by keybinding
    t.DispatchKey(ctx, "ctrl+s")

    // Or by action name
    t.Dispatch(ctx, "file.save")
}
```

## Documentation

- [Core Concepts](docs/core.md)
- [HTTP Transport](docs/http.md)
- [Action Transport](docs/action.md)
- [Dependency Injection](docs/di.md)
- [Integration with my other libraries](docs/ecosystem.md)
- [Middleware](docs/middleware.md)
- [OpenAPI Generation](docs/openapi.md)
- [Logging](docs/logging.md)

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
