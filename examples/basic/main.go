package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/mirkobrombin/go-module-router/v2/examples/basic/core/modules/ping"
	"github.com/mirkobrombin/go-module-router/v2/pkg/logger"
	"github.com/mirkobrombin/go-module-router/v2/pkg/router"
	"github.com/mirkobrombin/go-module-router/v2/pkg/swagger"
)

func main() {
	metaFlag := flag.Bool("meta", false, "generate OpenAPI JSON and exit")
	flag.Parse()

	// Create endpoint prototype (for registration and swagger)
	pingEndpoint := &ping.PingEndpoint{}

	// Handle 'swagger' subcommand OR --meta flag
	if (len(os.Args) > 1 && os.Args[1] == "swagger") || *metaFlag {
		doc, err := swagger.Build("Example API", "1.0.0", pingEndpoint)
		if err != nil {
			fmt.Fprintf(os.Stderr, "swagger build error: %v\n", err)
			os.Exit(1)
		}
		fmt.Println(string(doc))
		os.Exit(0)
	}

	// Create router
	r := router.New()
	r.SetLogger(logger.NewSlog(slog.Default()))

	// Register dependencies (by field name)
	r.Provide("PingService", ping.NewPingService())

	// Register endpoints (explicit, no init() magic)
	r.Register(pingEndpoint)

	slog.Info("🚀 Server listening on :8080")
	slog.Info("Try: curl http://localhost:8080/api/v1/ping?times=3")

	if err := r.Listen(":8080"); err != nil {
		slog.Error("server terminated", "err", err)
	}
}
