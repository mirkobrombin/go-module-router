package core

import "context"

// Handler is the interface that all endpoints must implement.
type Handler interface {
	Handle(ctx context.Context) (any, error)
}

// Pattern is used in struct tags to declare routing patterns.
// Different transports interpret different tags.
type Pattern struct{}
