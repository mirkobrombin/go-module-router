package http

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"

	"github.com/mirkobrombin/go-module-router/v2/pkg/core"
	"github.com/mirkobrombin/go-module-router/v2/pkg/logger"
)

// Transport handles HTTP-based routing.
type Transport struct {
	mux        *http.ServeMux
	container  *core.Container
	Logger     logger.Logger
	middleware []func(http.Handler) http.Handler
	prefix     string
}

// New creates a new HTTP transport.
func New() *Transport {
	return &Transport{
		mux:       http.NewServeMux(),
		container: core.NewContainer(),
		Logger:    logger.Nop,
	}
}

// Provide registers a dependency.
func (t *Transport) Provide(name string, instance any) {
	t.container.Provide(name, instance)
}

// Use adds middleware.
func (t *Transport) Use(mw func(http.Handler) http.Handler) {
	t.middleware = append(t.middleware, mw)
}

// Group creates a sub-transport with a prefix.
func (t *Transport) Group(prefix string) *Transport {
	return &Transport{
		mux:        t.mux,
		container:  t.container,
		Logger:     t.Logger,
		middleware: append([]func(http.Handler) http.Handler(nil), t.middleware...),
		prefix:     t.prefix + prefix,
	}
}

// Register adds an HTTP endpoint.
// Reads `method:"GET"` and `path:"/users/{id}"` tags from Pattern field.
func (t *Transport) Register(prototype core.Handler) {
	val := reflect.ValueOf(prototype)
	if val.Kind() != reflect.Ptr || val.Elem().Kind() != reflect.Struct {
		panic("Transport.Register: prototype must be a pointer to a struct")
	}

	elemType := val.Elem().Type()

	var method, path string
	for i := 0; i < elemType.NumField(); i++ {
		field := elemType.Field(i)
		if field.Type == reflect.TypeOf(core.Pattern{}) {
			method = field.Tag.Get("method")
			path = field.Tag.Get("path")
			break
		}
	}

	if method == "" || path == "" {
		panic(fmt.Sprintf("Transport.Register: struct %s missing Pattern with method/path tags", elemType.Name()))
	}

	fullPath := t.prefix + path
	pattern := fmt.Sprintf("%s %s", method, fullPath)

	t.Logger.Info("Registering route", "route", pattern, "handler", elemType.Name())

	var finalHandler http.Handler = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		// Create new instance
		newVal := reflect.New(elemType).Elem()
		newVal.Set(val.Elem())

		instance := newVal.Addr().Interface()

		// Inject dependencies
		t.container.Inject(instance)

		// Bind request data
		binder := core.NewBinder()
		binder.AddSource("path", func(key string) string { return req.PathValue(key) })
		binder.AddSource("query", func(key string) string { return req.URL.Query().Get(key) })
		binder.AddSource("header", func(key string) string { return req.Header.Get(key) })

		if err := binder.Bind(instance); err != nil {
			http.Error(w, fmt.Sprintf("Bad Request: %v", err), http.StatusBadRequest)
			return
		}

		// Bind JSON body if present
		if req.Body != nil {
			body, _ := io.ReadAll(req.Body)
			if len(body) > 0 {
				binder.BindJSON(instance, body)
			}
		}

		// Execute
		handler := instance.(core.Handler)
		resp, err := handler.Handle(req.Context())
		if err != nil {
			t.Logger.Error("Handler failed", "error", err)

			code := http.StatusInternalServerError
			var resp any = map[string]string{"error": err.Error()}

			// Check for optional interfaces
			type StatusCoder interface {
				StatusCode() int
			}
			type Payloader interface {
				Payload() any
			}

			if sc, ok := err.(StatusCoder); ok {
				code = sc.StatusCode()
			}
			if pl, ok := err.(Payloader); ok {
				resp = pl.Payload()
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(code)
			json.NewEncoder(w).Encode(resp)
			return
		}

		// Write response
		if resp != nil {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(resp)
		} else {
			w.WriteHeader(http.StatusNoContent)
		}
	})

	// Apply middleware
	for i := len(t.middleware) - 1; i >= 0; i-- {
		finalHandler = t.middleware[i](finalHandler)
	}

	t.mux.Handle(pattern, finalHandler)
}

// Listen starts the HTTP server.
func (t *Transport) Listen(addr string) error {
	t.Logger.Info("HTTP transport listening", "addr", addr)
	return http.ListenAndServe(addr, t.mux)
}

// Shutdown gracefully shuts down.
func (t *Transport) Shutdown(ctx context.Context) error {
	return nil
}

// ServeHTTP implements http.Handler.
func (t *Transport) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	t.mux.ServeHTTP(w, req)
}

// Mux returns the underlying ServeMux.
func (t *Transport) Mux() *http.ServeMux {
	return t.mux
}
