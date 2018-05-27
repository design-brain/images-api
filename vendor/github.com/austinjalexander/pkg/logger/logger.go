package logger

import (
	"log"
	"os"
	"sync"

	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

const (
	processName = `logger`
)

var (
	errParseLevel = errors.New(`unable to parse log level`)
	errProcessEnv = errors.New(`unable to process environment`)

	c    config
	once sync.Once
)

// config represents the configuration necessary for this pkg.
type config struct {
	Level string `envconfig:"LOG_LEVEL" default:"warn"`
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

		lvl, err := logrus.ParseLevel(c.Level)
		if err != nil {
			log.Fatal(errors.Wrap(err, errParseLevel.Error()))
		}
		logrus.SetLevel(lvl)
		logrus.SetFormatter(&logrus.JSONFormatter{})
		logrus.SetOutput(os.Stderr)
	}
}
