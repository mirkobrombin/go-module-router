# Go Module Router

A lightweight, zero-codegen library that lets you compose **modular HTTP
applications** in Go.

Define repositories, services, middleware and routes in plain Go code and
have them autowired at start-up—no reflection on the hot path, no generated
files to commit.

## Features

* **Module Auto-Discovery:** Register repositories, services, handlers,
  middleware, models and routes via `init()` functions—imports are enough.
* **Dependency Injection (No Codegen):** Factory functions receive their
  dependencies by name at runtime; use `SkipAutoWire` when you prefer to build
  instances yourself.
* **Multiple HTTP Engines:** Ship with fasthttp and net/http adapters—drop in
  your own by implementing a tiny `http.Engine` interface.
* **Pluggable Logging:** Bring any logger that satisfies a five-method
  interface; a Zap implementation (and a nop) are already included.
* **Route Introspection:** Build the route table without starting a server—
  handy for tests, docs, or static analysis tooling.
* **Zero Magic:** No struct tags, no global singletons, no hidden goroutines—
  just idiomatic Go.

## Getting Started

### Installation

```bash
go get github.com/mirkobrombin/go-module-router/v1
```

### Basic Usage

```go
package main

import (
	"time"

	"github.com/mirkobrombin/go-module-router/v1/http"
	"github.com/mirkobrombin/go-module-router/v1/logger"
	"github.com/mirkobrombin/go-module-router/v1/registry"
	"github.com/mirkobrombin/go-module-router/v1/router"

	_ "example.com/project/core/modules/ping" // ⚙ self-registering module

	"go.uber.org/zap"
)

func main() {
    zapL, _ := zap.NewDevelopment()
	defer zapL.Sync()

	eng := http.NewFastHTTP()
	lg := &logger.Zap{L: zapL}

	router.New(
		registry.Global(), // collected during imports
		nil,               // extra services you built manually
		eng,               // chosen HTTP engine
		router.Options{
			SessionDuration: 24 * time.Hour,
			Logger:          lg,
		},
	)

	if err := eng.Serve(":8080"); err != nil {
		lg.Fatal("server terminated", "err", err)
	}
}
```

> **Note:** Importing the package `.../modules/ping` is all that’s required; its
> `init()` registers a service, handler and route which the router wires up
> automatically.

For more detailed information, please refer to the documentation files in the
[docs/](docs/) directory.

## Documentation

* [Module Discovery](docs/modules.md)
* [Dependency Injection](docs/di.md)
* [HTTP Engines](docs/engines.md)
* [Logging](docs/logging.md)
* [Router Options](docs/options.md)
* [Route Introspection](docs/introspection.md)
* [Route Metadata](docs/router_metadata.md)

## License

Go Module Router is released under the MIT license.
See the [LICENSE](LICENSE) file for the full text.
