package router

import (
	"context"

	"github.com/mirkobrombin/go-module-router/v2/pkg/core"
	"github.com/mirkobrombin/go-module-router/v2/pkg/logger"
	"github.com/mirkobrombin/go-module-router/v2/pkg/transport/action"
	"github.com/mirkobrombin/go-module-router/v2/pkg/transport/http"
	"github.com/mirkobrombin/go-signal/v2/pkg/bus"
)

// Re-export core types for convenience
type Handler = core.Handler
type Pattern = core.Pattern

// Meta is an alias for Pattern (backward compatibility with HTTP examples)
type Meta = core.Pattern

// HTTP creates a new HTTP transport.
func HTTP() *http.Transport {
	return http.New()
}

// Action creates a new Action transport for GUI/CLI apps.
func Action() *action.Transport {
	return action.New()
}

// Router is a convenience wrapper that provides both transports.
type Router struct {
	HTTP   *http.Transport
	Action *action.Transport
	Logger logger.Logger
}

// New creates a new multi-transport router.
func New() *Router {
	return &Router{
		HTTP:   http.New(),
		Action: action.New(),
		Logger: logger.Nop,
	}
}

// SetLogger sets the logger for all transports.
func (r *Router) SetLogger(l logger.Logger) {
	r.Logger = l
	r.HTTP.Logger = l
	r.Action.Logger = l
}

// SetBus sets the event bus for action-based dispatching.
func (r *Router) SetBus(b *bus.Bus) {
	r.Action.Bus = b
}

// Provide registers a dependency in all transports.
func (r *Router) Provide(name string, instance any) {
	r.HTTP.Provide(name, instance)
	r.Action.Provide(name, instance)
}

// Register registers a handler in the appropriate transport based on tags.
func (r *Router) Register(prototype Handler) {
	// For now, register in both if applicable
	// TODO: Detect which transport based on tags
	defer func() { recover() }()
	r.HTTP.Register(prototype)
}

// RegisterAction registers an action handler.
func (r *Router) RegisterAction(prototype Handler) {
	r.Action.Register(prototype)
}

// Listen starts the HTTP transport.
func (r *Router) Listen(addr string) error {
	return r.HTTP.Listen(addr)
}

// Dispatch dispatches an action with an optional payload.
func (r *Router) Dispatch(ctx context.Context, action string, payload ...any) (any, error) {
	return r.Action.Dispatch(ctx, action, payload...)
}

// DispatchKey dispatches an action by keybinding.
func (r *Router) DispatchKey(ctx context.Context, key string) (any, error) {
	return r.Action.DispatchKey(ctx, key)
}
