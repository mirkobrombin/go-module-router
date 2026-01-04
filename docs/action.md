# Action Transport

The Action transport handles event-driven dispatching for GUI apps, TUI editors, CLI tools, and similar applications.

## Creating the Transport

```go
import "github.com/mirkobrombin/go-module-router/v2/pkg/transport/action"

t := action.New()
```

Or via the router facade:

```go
import "github.com/mirkobrombin/go-module-router/v2/pkg/router"

r := router.New()
// r.Action is the Action transport
```

## Defining Actions

Use `core.Pattern` with `action` and `keys` tags:

```go
type SaveAction struct {
    Meta     core.Pattern `action:"file.save" keys:"ctrl+s"`
    Document *Document
}

func (a *SaveAction) Handle(ctx context.Context) (any, error) {
    return a.Document.Save()
}
```

## Registering Actions

```go
t.Provide("Document", doc)
t.Register(&SaveAction{})
t.Register(&OpenAction{})
t.Register(&FindAction{})
```

## Dispatching

By action name:

```go
result, err := t.Dispatch(ctx, "file.save")
```

By keybinding:

```go
result, err := t.DispatchKey(ctx, "ctrl+s")
```

## Querying Registered Actions

```go
// Get all action names
actions := t.Actions()  // ["file.save", "file.open", ...]

// Get all keybindings
bindings := t.KeyBindings()  // {"ctrl+s": "file.save", ...}
```

## Event Bus Integration

The Action transport integrates with `go-signal` to provide an asynchronous notification system. 

By default, the transport uses the **`bus.Default()`** instance. Every dispatched action automatically emits its instance as an event on this bus after execution. 

```go
import "github.com/mirkobrombin/go-signal/v2/pkg/bus"

// Subscribe to any instance of SaveAction using the default bus
bus.Subscribe(nil, func(ctx context.Context, e *SaveAction) error {
    fmt.Println("File saved successfully!")
    return nil
})
```

If you need a custom bus instance, you can provide it during initialization:

```go
customBus := bus.New(bus.WithStrategy(bus.BestEffort))
t := action.New(action.WithBus(customBus))
```

For more advanced architecture patterns using this combo, see [Integration with other libraries](ecosystem.md).

## Use Cases

- **Text Editors:** Save, Open, Find, Replace commands
- **Media Players:** Play, Pause, Skip commands
- **Games:** Action bindings
- **CLI Tools:** Command dispatching
- **Desktop Apps:** Menu actions
