package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/mirkobrombin/go-module-router/v1/http"
	"github.com/mirkobrombin/go-module-router/v1/logger"
	"github.com/mirkobrombin/go-module-router/v1/registry"
	"github.com/mirkobrombin/go-module-router/v1/router"

	_ "github.com/mirkobrombin/go-module-router/examples/v1/core/modules/ping"
	"github.com/mirkobrombin/go-module-router/examples/v1/core/swagger"

	"go.uber.org/zap"
)

func main() {
	metaFlag := flag.Bool("meta", false, "generate OpenAPI JSON and exit")
	flag.Parse()

	if *metaFlag {
		doc, err := swagger.Build("Example API", "1.0.0")
		if err != nil {
			fmt.Fprintf(os.Stderr, "swagger build error: %v\n", err)
			os.Exit(1)
		}
		fmt.Println(string(doc))
		os.Exit(0)
	}

	zapL, _ := zap.NewDevelopment()
	defer zapL.Sync()

	eng := http.NewFastHTTP()
	lg := &logger.Zap{L: zapL}

	router.New(
		registry.Global(),
		nil,
		eng,
		router.Options{
			SessionDuration: 24 * time.Hour,
			Logger:          lg,
		},
	)

	if err := eng.Serve(":8080"); err != nil {
		lg.Fatal("server terminated", "err", err)
	}
}
