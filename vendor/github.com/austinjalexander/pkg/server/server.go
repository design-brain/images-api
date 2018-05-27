package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

const (
	addressStringFormat = `0.0.0.0:%d`
	gracefulShutdownMsg = `gracefully shutting down server...`
	processName         = `server`
)

var (
	errListenAndServe = errors.New(`unable to listen and serve`)
	errParseTimeout   = errors.New(`unable to parse timeout`)
	errProcessEnv     = errors.New(`unable to process environment`)
	errShutdown       = errors.New(`unable to shutdown gracefully`)

	c    config
	once sync.Once
	s    *http.Server
)

// config represents the configuration necessary for this pkg.
type config struct {
	ServerPort    int    `envconfig:"SERVER_PORT"`
	ServerTimeout string `envconfig:"SERVER_TIMEOUT"`
	timeout       time.Duration
}

// Init configures the current pkg using env vars (via initialize).
// While using init would be nice, Init allows us to be more
// explicit in different environments and helps to enforce
// one-time initialization.
func Init() {
	once.Do(initialize)
}

func initialize() {
	if c == (config{}) {
		err := envconfig.Process(processName, &c)
		if err != nil {
			log.Fatal(errors.Wrap(err, errProcessEnv.Error()))
		}
		d, err := time.ParseDuration(c.ServerTimeout)
		if err != nil {
			log.Fatal(errors.Wrap(err, errParseTimeout.Error()))
		}
		c.timeout = d
		s = &http.Server{
			Addr:         fmt.Sprintf(addressStringFormat, c.ServerPort),
			IdleTimeout:  c.timeout,
			ReadTimeout:  c.timeout,
			WriteTimeout: c.timeout,
		}
	}
}

// Run runs this pkg's configured server.
func Run(h http.Handler) {
	s.Handler = h
	go func() {
		err := s.ListenAndServe()
		if err != nil {
			log.Fatal(errors.Wrap(err, errListenAndServe.Error()))
		}
	}()

	nc := make(chan os.Signal, 1)
	signal.Notify(nc, os.Interrupt)
	<-nc

	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	err := s.Shutdown(ctx)
	if err != nil {
		log.Fatal(errors.Wrap(err, errShutdown.Error()))
	}
	log.Print(gracefulShutdownMsg)
	os.Exit(0)
}
