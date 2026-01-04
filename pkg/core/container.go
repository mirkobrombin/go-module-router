package core

import "reflect"

// Container manages dependency injection.
type Container struct {
	providers map[string]any
}

// NewContainer creates a new DI container.
func NewContainer() *Container {
	return &Container{
		providers: make(map[string]any),
	}
}

// Provide registers a dependency by name.
func (c *Container) Provide(name string, instance any) {
	c.providers[name] = instance
}

// Inject injects dependencies into a struct by matching field names.
func (c *Container) Inject(target any) {
	val := reflect.ValueOf(target)
	if val.Kind() != reflect.Ptr || val.Elem().Kind() != reflect.Struct {
		return
	}

	elem := val.Elem()
	elemType := elem.Type()

	for i := 0; i < elem.NumField(); i++ {
		fieldType := elemType.Field(i)
		fieldVal := elem.Field(i)

		if !fieldVal.CanSet() {
			continue
		}

		if dep, ok := c.providers[fieldType.Name]; ok {
			depVal := reflect.ValueOf(dep)
			if depVal.Type().AssignableTo(fieldType.Type) {
				fieldVal.Set(depVal)
			}
		}
	}
}

// Get retrieves a dependency by name.
func (c *Container) Get(name string) (any, bool) {
	v, ok := c.providers[name]
	return v, ok
}
