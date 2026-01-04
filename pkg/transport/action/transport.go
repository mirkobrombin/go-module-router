package action

import (
	"context"
	"fmt"
	"reflect"
	"sync"

	"github.com/mirkobrombin/go-module-router/v2/pkg/core"
	"github.com/mirkobrombin/go-module-router/v2/pkg/logger"
)

// Transport handles action-based routing for GUI/CLI applications.
type Transport struct {
	container *core.Container
	Logger    logger.Logger
	handlers  map[string]core.Handler
	keys      map[string]string // keybinding -> action
	mu        sync.RWMutex
}

// New creates a new action transport.
func New() *Transport {
	return &Transport{
		container: core.NewContainer(),
		Logger:    logger.Nop,
		handlers:  make(map[string]core.Handler),
		keys:      make(map[string]string),
	}
}

// Provide registers a dependency.
func (t *Transport) Provide(name string, instance any) {
	t.container.Provide(name, instance)
}

// Register adds an action handler.
// Reads `action:"name"` and `keys:"ctrl+s"` tags from Pattern field.
func (t *Transport) Register(prototype core.Handler) {
	val := reflect.ValueOf(prototype)
	if val.Kind() != reflect.Ptr || val.Elem().Kind() != reflect.Struct {
		panic("Transport.Register: prototype must be a pointer to a struct")
	}

	elemType := val.Elem().Type()

	var actionName, keyBinding string
	for i := 0; i < elemType.NumField(); i++ {
		field := elemType.Field(i)
		if field.Type == reflect.TypeOf(core.Pattern{}) {
			actionName = field.Tag.Get("action")
			keyBinding = field.Tag.Get("keys")
			break
		}
	}

	if actionName == "" {
		panic(fmt.Sprintf("Transport.Register: struct %s missing Pattern with action tag", elemType.Name()))
	}

	t.mu.Lock()
	defer t.mu.Unlock()

	t.handlers[actionName] = prototype
	if keyBinding != "" {
		t.keys[keyBinding] = actionName
	}

	t.Logger.Info("Registered action", "action", actionName, "keys", keyBinding)
}

// Dispatch executes an action by name.
func (t *Transport) Dispatch(ctx context.Context, action string) (any, error) {
	t.mu.RLock()
	prototype, ok := t.handlers[action]
	t.mu.RUnlock()

	if !ok {
		return nil, fmt.Errorf("action not found: %s", action)
	}

	// Create new instance
	val := reflect.ValueOf(prototype)
	elemType := val.Elem().Type()
	newVal := reflect.New(elemType).Elem()
	newVal.Set(val.Elem())

	// Inject dependencies
	instance := newVal.Addr().Interface()
	t.container.Inject(instance)

	// Execute
	handler := instance.(core.Handler)
	return handler.Handle(ctx)
}

// DispatchKey executes an action by keybinding.
func (t *Transport) DispatchKey(ctx context.Context, key string) (any, error) {
	t.mu.RLock()
	action, ok := t.keys[key]
	t.mu.RUnlock()

	if !ok {
		return nil, fmt.Errorf("no action bound to key: %s", key)
	}

	return t.Dispatch(ctx, action)
}

// Actions returns all registered action names.
func (t *Transport) Actions() []string {
	t.mu.RLock()
	defer t.mu.RUnlock()

	actions := make([]string, 0, len(t.handlers))
	for name := range t.handlers {
		actions = append(actions, name)
	}
	return actions
}

// KeyBindings returns all registered keybindings.
func (t *Transport) KeyBindings() map[string]string {
	t.mu.RLock()
	defer t.mu.RUnlock()

	bindings := make(map[string]string, len(t.keys))
	for k, v := range t.keys {
		bindings[k] = v
	}
	return bindings
}
