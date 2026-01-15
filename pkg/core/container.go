package core

import (
	"github.com/mirkobrombin/go-foundation/pkg/di"
)

// Container manages dependency injection.
type Container struct {
	*di.Container
}

// NewContainer creates a new DI container.
func NewContainer() *Container {
	return &Container{
		Container: di.New(),
	}
}
