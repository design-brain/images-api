package images

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"

	ipb "github.com/design-brain/images-api/rpc/images"
)

const (
	processName = `images`
)

var (
	errParseTimeout = errors.New(`unable to parse timeout`)
	errProcessEnv   = errors.New(`unable to process environment`)

	c    config
	once sync.Once
	s    *Service
)

// config represents the configuration necessary for this pkg.
type config struct {
	ServerTimeout string `envconfig:"SERVER_TIMEOUT"`
	timeout       time.Duration
}

// Service implements the ipb.Manage interface.
type Service struct {
}

// Init configures the current pkg (via initialize).
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
		if s == nil {
			s = &Service{}
		}
	}
}

// Svc returns this pkg's configured Service.
func Svc() *Service {
	return s
}

// Fetch fetches...
func (s *Service) Fetch(ctx context.Context, i *ipb.Image) (*ipb.Image, error) {
	return i, nil
}

// Upload uploads...
func (s *Service) Upload(ctx context.Context, i *ipb.Image) (*ipb.Image, error) {
	return i, nil
}
