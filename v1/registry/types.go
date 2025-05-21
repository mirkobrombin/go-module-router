package registry

import "time"

type ServiceInit func(repos map[string]any) any
type RepoInit func(deps map[string]any) any
type HandlerInit func(services map[string]any) any
type MiddlewareInit func(services map[string]any, sessionDur time.Duration) any
type ModelProvider func() []any
type RouteProvider func() []Route

type Route struct {
	Method      string
	Path        string
	HandlerName string
	Middleware  []string
	Permissions []string
	Meta        map[string]any
}
