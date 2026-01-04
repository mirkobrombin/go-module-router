package ping

import (
	"context"

	"github.com/mirkobrombin/go-module-router/v2/pkg/core"
)

type pingResponse struct {
	Message []string `json:"message"`
}

// PingEndpoint is the declarative endpoint for /api/v1/ping
type PingEndpoint struct {
	// Routing metadata (declarative)
	Meta core.Pattern `method:"GET" path:"/api/v1/ping"`

	// Query parameters (auto-bound)
	Times int `query:"times" default:"1"`

	// Dependencies (injected by field name)
	PingService PingService
}

// OpenAPIMeta returns metadata for OpenAPI/Swagger generation.
func (e *PingEndpoint) OpenAPIMeta() map[string]any {
	return map[string]any{
		"summary":     "Ping – returns one or more \"pong\" strings",
		"description": "Optional query-param `times` repeats the reply. Example: `/api/v1/ping?times=3`",
		"parameters": []map[string]any{
			{
				"name":     "times",
				"in":       "query",
				"required": false,
				"schema": map[string]any{
					"type":    "integer",
					"minimum": 1,
				},
			},
		},
		"responses": map[int]any{
			200: "PingResponse",
		},
	}
}

// Handle implements router.Handler
func (e *PingEndpoint) Handle(ctx context.Context) (any, error) {
	if e.Times < 1 {
		e.Times = 1
	}

	out := make([]string, e.Times)
	for i := range out {
		out[i] = e.PingService.Pong()
	}

	return pingResponse{Message: out}, nil
}
