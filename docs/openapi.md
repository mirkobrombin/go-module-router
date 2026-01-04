# OpenAPI Generation

Generate OpenAPI/Swagger documents from your HTTP handlers.

## Basic Usage

```go
import "github.com/mirkobrombin/go-module-router/v2/pkg/swagger"

doc, err := swagger.Build("My API", "1.0.0", 
    &GetUser{},
    &CreateUser{},
)
fmt.Println(string(doc))
```

## Adding Metadata

Implement `OpenAPIMeta()` on your handler:

```go
type GetUser struct {
    Meta core.Pattern `method:"GET" path:"/users/{id}"`
    ID   string       `path:"id"`
}

func (e *GetUser) OpenAPIMeta() map[string]any {
    return map[string]any{
        "summary":     "Get a user by ID",
        "description": "Returns the user object for the given ID.",
        "responses": map[int]any{
            200: "User",
            404: "Not Found",
        },
    }
}
```

## CLI Integration

```bash
go run main.go swagger > openapi.json
go run main.go --meta > openapi.json
```
