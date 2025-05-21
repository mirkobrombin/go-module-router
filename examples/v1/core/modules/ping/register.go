//go:build !ping
// +build !ping

package ping

import (
	"github.com/mirkobrombin/go-module-router/v1/registry"
	"github.com/valyala/fasthttp"
)

func init() {
	registry.RegisterService("PingService",
		func(_ map[string]any) any { return pingService{} })

	registry.RegisterHandler("PingHandler",
		func(s map[string]any) any {
			return &PingHandler{svc: s["PingService"].(PingService)}
		})

	registry.RegisterRoutes(func() []registry.Route {
		return []registry.Route{
			{
				Method:      fasthttp.MethodGet,
				Path:        "/api/v1/ping",
				HandlerName: "PingHandler.Pong",
				Meta: map[string]any{
					"summary":     "Ping – returns one or more \"pong\" strings",
					"description": "Optional query-param `times` repeats the reply. Example: `/api/v1/ping?times=3`",
					"parameters": []any{
						map[string]any{
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
				},
			},
		}
	})

}
