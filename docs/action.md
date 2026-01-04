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

## Use Cases

- **Text Editors:** Save, Open, Find, Replace commands
- **Media Players:** Play, Pause, Skip commands
- **Games:** Action bindings
- **CLI Tools:** Command dispatching
- **Desktop Apps:** Menu actions
