# Logging

The router never imports a concrete logging library.  
Instead it depends on this tiny interface:

```go
type Logger interface {
    Debug(msg string, kv ...any)
    Info(msg string, kv ...any)
    Warn(msg string, kv ...any)
    Error(msg string, kv ...any)
    Fatal(msg string, kv ...any)
}
```

this allows you to plug in any logger you like.

### Included implementations

* `logger.Zap` – wraps a `*zap.Logger`.
* `logger.Nop` – silent no-op logger (default).

### Plugging your own logger

```go
type slogAdapter struct{ h *slog.Logger }

func (s slogAdapter) Debug(m string, kv ...any) { s.h.Debug(m, kv...) }

router.New(registry.Global(), nil, eng, router.Options{
    Logger: slogAdapter{h: slog.New(...)} ,
})
```

