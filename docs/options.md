# Router Options

```go
type Options struct {
    SessionDuration time.Duration 
    Logger          logger.Logger 
    OnError         func(error)  
}
```

that translates to:

| Field             | Type               | Description                                   |
|-------------------|--------------------|-----------------------------------------------|
| `SessionDuration` | `time.Duration`    | Passed to middleware factories.               |
| `Logger`          | `logger.Logger`    | Default: `logger.Nop`.                        |
| `OnError`         | `func(error)`      | Called when a route/handler fails to resolve. |

* **SessionDuration** purely propagates into middleware factories.
  If you have no auth layer, ignore it.
* **Logger** affects every internal debug / error message.
* **OnError** lets you trap reflection errors instead of logging them only.


