package core

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
)

// BindFunc is a function that extracts a value by key from a source.
type BindFunc func(key string) string

// Binder binds values to struct fields based on tags.
type Binder struct {
	sources map[string]BindFunc
}

// NewBinder creates a new binder.
func NewBinder() *Binder {
	return &Binder{
		sources: make(map[string]BindFunc),
	}
}

// AddSource registers a binding source (e.g., "query", "path", "header").
func (b *Binder) AddSource(tag string, fn BindFunc) {
	b.sources[tag] = fn
}

// Bind populates struct fields from registered sources.
func (b *Binder) Bind(target any) error {
	v := reflect.ValueOf(target)
	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("target must be a pointer to a struct")
	}

	elem := v.Elem()
	typ := elem.Type()

	for i := 0; i < elem.NumField(); i++ {
		field := elem.Field(i)
		fieldType := typ.Field(i)

		if !field.CanSet() {
			continue
		}

		// Try each registered source
		var valStr string
		for tag, fn := range b.sources {
			if key, ok := fieldType.Tag.Lookup(tag); ok {
				valStr = fn(key)
				if valStr != "" {
					break
				}
			}
		}

		// Check for default
		if valStr == "" {
			if def, ok := fieldType.Tag.Lookup("default"); ok {
				valStr = def
			}
		}

		if valStr == "" {
			continue
		}

		if err := setField(field, valStr); err != nil {
			return fmt.Errorf("failed to bind field %s: %w", fieldType.Name, err)
		}
	}

	return nil
}

// BindJSON binds a JSON body to a struct field tagged with body:"json".
func (b *Binder) BindJSON(target any, data []byte) error {
	v := reflect.ValueOf(target)
	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Struct {
		return nil
	}

	elem := v.Elem()
	typ := elem.Type()

	for i := 0; i < elem.NumField(); i++ {
		field := elem.Field(i)
		fieldType := typ.Field(i)

		if tag, ok := fieldType.Tag.Lookup("body"); ok && tag == "json" && field.CanSet() {
			return json.Unmarshal(data, field.Addr().Interface())
		}
	}

	return nil
}

func setField(field reflect.Value, valStr string) error {
	switch field.Kind() {
	case reflect.String:
		field.SetString(valStr)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		val, err := strconv.ParseInt(valStr, 10, 64)
		if err != nil {
			return err
		}
		field.SetInt(val)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		val, err := strconv.ParseUint(valStr, 10, 64)
		if err != nil {
			return err
		}
		field.SetUint(val)
	case reflect.Bool:
		val, err := strconv.ParseBool(valStr)
		if err != nil {
			return err
		}
		field.SetBool(val)
	case reflect.Float32, reflect.Float64:
		val, err := strconv.ParseFloat(valStr, 64)
		if err != nil {
			return err
		}
		field.SetFloat(val)
	default:
		return fmt.Errorf("unsupported type %s", field.Kind())
	}
	return nil
}
