package main

import (
	"context"
	"fmt"
	"time"

	"github.com/mirkobrombin/go-module-router/v2/pkg/router"
	"github.com/mirkobrombin/go-revert/v2/pkg/engine"
	"github.com/mirkobrombin/go-revert/v2/pkg/step"
	"github.com/mirkobrombin/go-signal/v2/pkg/bus"
)

// Events
type OrderProcessed struct {
	OrderID string
	Status  string
}

// Action Handler
type CheckoutAction struct {
	Meta    router.Meta `action:"order.checkout" keys:"ctrl+b"`
	OrderID string      `json:"order_id"`
}

func (a *CheckoutAction) Handle(ctx context.Context) (any, error) {
	fmt.Printf("[Action] Starting checkout for order: %s\n", a.OrderID)

	// Create a reversible workflow
	wf := engine.New()
	wf.Add(&step.Basic{
		Name: "Reserve Items",
		OnExecute: func(ctx context.Context) error {
			fmt.Println("  [Workflow] Items reserved.")
			return nil
		},
		OnUndo: func(ctx context.Context) error {
			fmt.Println("  [Workflow] Items released.")
			return nil
		},
	})

	wf.Add(&step.Basic{
		Name: "Process Payment",
		OnExecute: func(ctx context.Context) error {
			fmt.Println("  [Workflow] Payment processed.")
			return nil
		},
		OnUndo: func(ctx context.Context) error {
			fmt.Println("  [Workflow] Payment refunded.")
			return nil
		},
	})

	if err := wf.Run(ctx); err != nil {
		return nil, err
	}

	return map[string]string{"status": "completed", "order_id": a.OrderID}, nil
}

func main() {
	r := router.New()

	// Subscribe to OrderProcessed events using the Default Bus (passing nil as bus)
	// The Action transport automatically emits events on this bus by default.
	bus.Subscribe(nil, func(ctx context.Context, e *CheckoutAction) error {
		fmt.Printf("[Signal] Async notification: Order %s was processed successfully!\n", e.OrderID)
		return nil
	})

	// Register Action
	r.RegisterAction(&CheckoutAction{})

	fmt.Println("--- Dispatching Checkout Action ---")
	// Passing payload as a map
	res, err := r.Dispatch(context.Background(), "order.checkout", map[string]any{
		"order_id": "ORD-2024-ABC",
	})
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Response: %v\n", res)

	// Wait a bit for async signals
	time.Sleep(500 * time.Millisecond)
	fmt.Println("Done.")
}
