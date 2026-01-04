# Logging

All transports support a common logger interface.

## Logger Interface

```go
type Logger interface {
    Debug(msg string, kv ...any)
    Info(msg string, kv ...any)
    Warn(msg string, kv ...any)
    Error(msg string, kv ...any)
    Fatal(msg string, kv ...any)
}
```

## Setting the Logger

### Via Router Facade

```go
r := router.New()
r.SetLogger(logger.NewSlog(slog.Default()))
```

### Via Transport

```go
t := http.New()
t.Logger = logger.NewSlog(slog.Default())
```

## Available Implementations

### Nop Logger (Default)

Silent, no output.

```go
t.Logger = logger.Nop
```

### Slog Logger

```go
import "github.com/mirkobrombin/go-module-router/v2/pkg/logger"

t.Logger = logger.NewSlog(slog.Default())
```

### Zap Logger

```go
zapL, _ := zap.NewDevelopment()
t.Logger = &logger.Zap{L: zapL}
```

## Custom Logger

Implement the `Logger` interface for your own logging solution.
