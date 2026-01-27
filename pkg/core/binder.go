package core

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"

	"github.com/mirkobrombin/go-foundation/pkg/tags"
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
	val := reflect.ValueOf(target)
	if val.Kind() != reflect.Ptr || val.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("target must be a pointer to a struct")
	}

	elem := val.Elem()

	// Track which fields have been set to avoid overwriting with lower priority sources?
	// Or just use the first value found (since sources loop order is random map iteration).
	// To preserve original semantics: "Try each registered source... if valStr != '' break".
	// The original implementation iterated fields, then iterated sources.
	// We want to iterate sources (cached), then fields.
	// This changes order: We iterate Source A -> all fields. Source B -> all fields.
	// If Source A and B both set Field X, the last one wins.
	// Original: For Field X, iterate sources (random order). First sets wins.
	// Since original source order was random, "Last wins" in our new loop is equivalent to "First wins" in a reversed random order.
	// So it is acceptable.

	// Collect valuesmap[fieldIndex]value
	values := make(map[int]string)

	// 1. Process explicit sources
	for tag, fn := range b.sources {
		parser := tags.NewParser(tag)
		fields := parser.ParseStruct(target)

		for _, meta := range fields {
			// Skip if already set? (First Wins strategy adaptation)
			// If we want "First Scanned Source Wins", we check existence.
			// Since source iteration is random, this is fine.
			if _, ok := values[meta.Index]; ok {
				continue
			}

			key := meta.RawTag
			if key == "" {
				continue
			}

			if val := fn(key); val != "" {
				values[meta.Index] = val
			}
		}
	}

	// 2. Process defaults
	// Only if not set by sources
	defParser := tags.NewParser("default")
	defFields := defParser.ParseStruct(target)
	for _, meta := range defFields {
		if _, ok := values[meta.Index]; !ok {
			if meta.RawTag != "" {
				values[meta.Index] = meta.RawTag
			}
		}
	}

	// 3. Apply values
	for idx, val := range values {
		field := elem.Field(idx)
		if !field.CanSet() {
			continue
		}
		if err := setField(field, val); err != nil {
			typeField := elem.Type().Field(idx)
			return fmt.Errorf("failed to bind field %s: %w", typeField.Name, err)
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
