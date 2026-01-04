package action

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"sync"

	"github.com/mirkobrombin/go-module-router/v2/pkg/core"
	"github.com/mirkobrombin/go-module-router/v2/pkg/logger"
	"github.com/mirkobrombin/go-signal/v2/pkg/bus"
)

// Transport handles action-based routing for GUI/CLI applications.
type Transport struct {
	container *core.Container
	Logger    logger.Logger
	handlers  map[string]core.Handler
	keys      map[string]string // keybinding -> action
	Bus       *bus.Bus
	mu        sync.RWMutex
}

type Option func(*Transport)

func WithBus(b *bus.Bus) Option {
	return func(t *Transport) { t.Bus = b }
}

// New creates a new action transport.
func New(opts ...Option) *Transport {
	t := &Transport{
		container: core.NewContainer(),
		Logger:    logger.Nop,
		handlers:  make(map[string]core.Handler),
		keys:      make(map[string]string),
		Bus:       bus.Default(),
	}
	for _, opt := range opts {
		opt(t)
	}
	return t
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

// Dispatch executes an action by name with an optional payload.
func (t *Transport) Dispatch(ctx context.Context, action string, payload ...any) (any, error) {
	t.mu.RLock()
	prototype, ok := t.handlers[action]
	busInstance := t.Bus
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

	// Real payload binding
	if len(payload) > 0 && payload[0] != nil {
		if err := t.applyPayload(instance, payload[0]); err != nil {
			return nil, fmt.Errorf("payload binding failed: %w", err)
		}
	}

	// Execute
	handler := instance.(core.Handler)
	res, err := handler.Handle(ctx)

	// If a bus is present, emit the action as an event asynchronously
	if busInstance != nil {
		bus.EmitAsync(ctx, busInstance, instance)
	}

	return res, err
}

// applyPayload maps data from the payload to the target struct.
// It supports map[string]any and structs, using reflection and "json" tags.
func (t *Transport) applyPayload(target any, payload any) error {
	dstVal := reflect.ValueOf(target).Elem()
	dstType := dstVal.Type()

	srcVal := reflect.ValueOf(payload)
	if srcVal.Kind() == reflect.Ptr {
		srcVal = srcVal.Elem()
	}

	switch srcVal.Kind() {
	case reflect.Map:
		for _, key := range srcVal.MapKeys() {
			k := fmt.Sprintf("%v", key.Interface())
			v := srcVal.MapIndex(key)
			t.setFieldByNameOrTag(dstVal, dstType, k, v)
		}
	case reflect.Struct:
		srcType := srcVal.Type()
		for i := 0; i < srcVal.NumField(); i++ {
			field := srcType.Field(i)
			val := srcVal.Field(i)
			t.setFieldByNameOrTag(dstVal, dstType, field.Name, val)
		}
	}

	return nil
}

func (t *Transport) setFieldByNameOrTag(dst reflect.Value, dstType reflect.Type, name string, val reflect.Value) {
	for i := 0; i < dst.NumField(); i++ {
		fieldMeta := dstType.Field(i)
		field := dst.Field(i)

		if !field.CanSet() {
			continue
		}

		// Match by name or json tag
		tag := fieldMeta.Tag.Get("json")
		if strings.EqualFold(fieldMeta.Name, name) || (tag != "" && strings.Split(tag, ",")[0] == name) {
			if field.Type() == val.Type() {
				field.Set(val)
			} else if val.Kind() == reflect.Interface && reflect.TypeOf(val.Interface()) == field.Type() {
				field.Set(reflect.ValueOf(val.Interface()))
			}
			// Future: add type conversion (string to int, etc) if needed
			break
		}
	}
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
