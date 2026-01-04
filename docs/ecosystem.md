# Integration with other libraries

One of the most powerful features of `go-module-router` is its seamless integration with other of my libraries. By combining the router with `go-revert` and `go-signal`, you can build complex, reliable, and event-driven architectures with very little code.

## The Pattern: Saga + Event-Driven Actions

When building complex operations (like a checkout process, or a multi-step data migration), you often face two challenges:
1. **Partial Failures**: If the third step fails, you need to undo the first two.
2. **Side Effects**: Other parts of the system (notifications, analytics) need to know when an operation completes, but you don't want to bloat your main handler with that logic.

### 1. Robust Logic with `go-revert`
Inside your Action or HTTP handler, use `go-revert` to define a reversible workflow. If any step fails, the workflow automatically compensates (rolls back) the successful ones.

```go
func (a *OrderAction) Handle(ctx context.Context) (any, error) {
    wf := engine.New()
    wf.Add(&step.Basic{
        Name: "Charge Card",
        OnExecute: func(ctx context.Context) error { ... },
        OnUndo:    func(ctx context.Context) error { ... }, // Refund
    })
    
    return wf.Run(ctx) // If this fails, Card is refunded automatically
}
```

### 2. Decoupled Side Effects with `go-signal`
The real benefit here is that **you don't manually call `bus.Emit`**. When an action is dispatched via the router, the `ActionTransport` automatically emits the action instance as an event on the bus *after* your `Handle` method returns.

This creates a clean separation of concerns:
- **The Handler**: Only cares about the business logic (e.g., saving an order).
- **The Router**: Orchestrates execution and notifies the system.
- **The Listeners**: Handle side effects (email, logs, analytics) without knowing about each other.

```go
import "github.com/mirkobrombin/go-signal/v2/pkg/bus"

b := bus.New()
r.SetBus(b)

// Every time order.checkout is dispatched, this listener will trigger.
bus.Subscribe(b, func(ctx context.Context, e *OrderAction) error {
    fmt.Printf("Side Effect: Email sent for Order %s!\n", e.OrderID)
    return nil
})

// The router will Handle() it, then automatically EmitAsync() it on the bus.
r.Dispatch(ctx, "order.checkout", payload)
```

## Summary of the Combo

| Library | Responsibility | Role in Integration |
|---------|----------------|----------------------|
| **`go-module-router`** | Orchestration | Receives the intent, maps the payload, and executes the handler. |
| **`go-revert`** | Reliability | Ensures that complex operations are atomic (all-or-nothing) via compensations. |
| **`go-signal`** | Decoupling | Notifies the rest of the system asynchronusly that the intent was fulfilled. |

For a complete working example, see `examples/integration/main.go`.
