# Module Discovery

`go-module-router` is built around **self-registering modules**.  
A module is any Go package that, in an `init()` function, calls one or more
`registry.Register*` helpers.

```go
func init() {
    registry.RegisterService("PingService",
        func(_ map[string]any) any { return pingService{} })

    registry.RegisterHandler("PingHandler",
        func(s map[string]any) any {
            return &PingHandler{svc: s["PingService"].(PingService)}
        })

    registry.RegisterRoutes(func() []registry.Route{{
        Method:      http.MethodGet,
        Path:        "/api/v1/ping",
        HandlerName: "PingHandler.Pong",
        Meta: map[string]any{
            "summary": "Ping – returns one or more \"pong\" strings",
            "parameters": []any{
                map[string]any{
                    "name":"times","in":"query","schema":map[string]any{"type":"integer","minimum":1},
                },
            },
            "responses": map[int]any{200:"PingResponse"},
        },
    }})
}
```

### Naming convention

* **Service / Repository names** are arbitrary, but *unique* inside the
  process.
  The router uses them as keys to wire dependencies.
* **HandlerName** in a route is the fully-qualified form
  `"<HandlerStruct>.<Method>"`.

### Import side-effects

Simply importing the package is enough to make its pieces available:

```go
import _ "example.com/app/core/modules/ping"
```

> No reflection on the hot path – discovery only happens once, during
> program start-up, when Go runs `init()` functions.

### Toggling modules and handlers via environment variables

You can skip entire modules or individual handlers without touching your 
`init()` registrations by setting one of two env-vars before startup:

| **Toggle Type**      | **Environment Variable**                     | **Effect**                                                                 |
|-----------------------|---------------------------------------------|-----------------------------------------------------------------------------|
| **Module-level**      | `GMR_MOD_OFF_<MODULE_NAME>`                 | Skips *every* route whose `HandlerName` begins with `Ping…`.               |
| **Handler-level**     | `GMR_HAND_OFF_<HANDLER_KEY>`                | Skips only the `PingHandler.Pong` route.                                   |

**Notes**:  
- `<MODULE_NAME>` is the uppercased name of the module you want to disable.  
- `<HANDLER_KEY>` is your `HandlerName` with `.` replaced by `_` and uppercased.

The router checks these flags on startup and simply omits any matching routes.

### Compile-time toggles via Go build tags

Thanks to Go's build tags, you can also exclude entire modules at build time.
Like for example, if you have a `billing` module that you want to include
only in a specific build, you can use a build tag to conditionally include
it.

Of course, this is not a feature of this library, but a Go feature, so enjoy
the full power of Go's build system.

```go
//go:build !payments
// +build !payments

package payments

import "github.com/mirkobrombin/go-module-router/v1/registry"

func init() {
    registry.RegisterRoutes(func() []registry.Route {
        return []registry.Route{{
            Method:      http.MethodPost,
            Path:        "/api/v1/charge",
            HandlerName: "PaymentHandler.Charge",
        }}
    })
}
```

Then exclude it when building with:

```bash
go build -tags payments ./…
```

Just remember to add a stub file for the module in the main package, so that
you can import it without making the compiler complain about missing
dependencies.

For example in your hypothetical `billing` module, you can create a file
`billing_stub.go` with the following content:

```go
//go:build payments
// +build payments

package billing
```

This is the way I usually do it, so keeping all the modules enabled and pass
the build tags to the compiler only to disable the ones I don't want to
build.

I think this is the best way but you can also do the opposite, keeping all
the modules disabled and enabling the ones you want to build with the build
tags, to achieve this, simply invert the build tag in the stub file and all
the other files in the module.


