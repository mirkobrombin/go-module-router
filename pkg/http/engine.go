package http

import "context"

type Handler = any
type Middleware = func(Handler) Handler

type Engine interface {
	Handle(method, path string, h Handler)
	Use(mw Middleware)
	Group(prefix string) Engine

	Serve(addr string) error
	Shutdown(ctx context.Context) error

	Unwrap() any
}
