package main

import (
	"flag"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/austinjalexander/pkg/db"
	"github.com/austinjalexander/pkg/dotenv"
	"github.com/austinjalexander/pkg/logger"
	"github.com/austinjalexander/pkg/server"
	"github.com/design-brain/images-api/internal/api/handlers/healthcheck"
	"github.com/design-brain/images-api/internal/api/services/images"
	ipb "github.com/design-brain/images-api/rpc/images"
)

var (
	dotenvrun = flag.Bool("dotenv", false, `run dotenv.Run(".env")`)
)

func init() {
	// Set environment variables for local development.
	if *dotenvrun {
		err := dotenv.Run(".env")
		if err != nil {
			log.Print(err)
		}
	}

	// Initialize logrus logger configuration.
	logger.Init()

	// Initialize database interface configuration.
	db.Init()

	// Initialize non-Twirp handlers.
	healthcheck.Init()

	// Initialize Twirp services.
	images.Init()

	// Initialize server configuration.
	server.Init()
}

func main() {
	// Configure routes.
	mux := http.NewServeMux()

	// Non-Twirp routes.
	mux.Handle(healthcheck.Path(), healthcheck.Handler(time.Now().UTC()))

	// Twirp routes.
	mux.Handle(ipb.ManagePathPrefix, ipb.NewManageServer(images.Svc(), nil))

	// Run server.
	server.Run(mux)
}
