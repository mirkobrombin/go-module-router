# Router Metadata

Every `registry.Route` now has an **optional** `Meta map[string]any` field you 
can use for *anything*—the router itself never touches it at runtime.

| Use-case                    | Example keys                                           |
| --------------------------- | ------------------------------------------------------ |
| **OpenAPI / Swagger**       | `summary`, `description`, `parameters`, `responses`, … |
| **Rate-limit declarations** | `rateLimit`: `"100req/min"`                            |
| **AB test labels**          | `bucket`: `"A1"`                                       |
| **GraphQL operation names** | `operation`: `"GetUser"`                               |
| **Custom middleware hints** | `skipAuth`: `true`                                     |

```go
registry.RegisterRoutes(func() []registry.Route{{
    Method:      http.MethodGet,
    Path:        "/api/v1/ping",
    HandlerName: "PingHandler.Pong",

    Meta: map[string]any{
        "summary":     "Ping – returns one or more \"pong\" strings",
        "description": "Optional query-param `times` repeats the reply. Example: `/api/v1/ping?times=3`",
        "parameters": []any{
            map[string]any{
                "name": "times",
                "in":   "query",
                "schema": map[string]any{
                    "type":    "integer",
                    "minimum": 1,
                },
            },
        },
        "responses": map[string]any{
            "200": "PingResponse",
        },
        "rateLimit": "5req/sec",
    },
}})
```

## Consuming the metadata

Since all metadata lives in the registry, **you can walk**
`registry.Global().RouteProviders` at startup to generate documentation,
wire middleware, export metrics, etc.

For a full, working example of a minimal OpenAPI‐3.0 generator (producing JSON 
like the one below), see `examples/v1/core/swagger/generator.go`.

```json
{
  "openapi": "3.0.3",
  "info": {
    "title": "Example API",
    "version": "1.0.0"
  },
  "paths": {
    "/api/v1/ping": {
      "get": {
        "summary": "Ping – returns one or more \"pong\" strings",
        "description": "Optional query-param `times` repeats the reply. Example: `/api/v1/ping?times=3`",
        "parameters": [
          {
            "name": "times",
            "in": "query",
            "schema": {
              "type": "integer",
              "minimum": 1
            }
          }
        ],
        "responses": {
          "200": {
            "description": "OK"
          }
        }
      }
    }
  }
}
```

To toggle metadata-driven docs generation at build-time, wrap your 
`RegisterRoutes` behind a build tag:

```go
//go:build swagger
```

---

> **Tip:** Want *even more* type safety? Attach a strongly-typed struct 
> instead of a raw `map[string]any`, and type-assert it when you consume it.
