package healthcheck

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/austinjalexander/pkg/db"
	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

const (
	processName = `healthcheck`
	statusOkMsg = `OK`
)

var (
	errDbPing      = errors.New(`unable to ping database`)
	errMarshalResp = errors.New(`problem marshaling response`)
	errProcessEnv  = errors.New(`unable to process environment`)
	errWriteResp   = errors.New(`problem writing response`)

	c    config
	once sync.Once
)

// config represents the configuration necessary for this pkg.
type config struct {
	Path string `envconfig:"HEALTHCHECK_PATH"`
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
	}
}

// Handler takes a time constant and handles healthchecks by pinging the database.
func Handler(t time.Time) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := db.DB().Ping()
		if err != nil {
			err = errors.Wrap(err, errDbPing.Error())
			log.Printf("%+v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		// Construct response body.
		resp := struct {
			DbAvailable bool   `json:"db_available"`
			StartupTime string `json:"startup_time"`
			Status      string `json:"status"`
		}{
			DbAvailable: true,
			StartupTime: t.Format(time.RFC3339),
			Status:      statusOkMsg,
		}
		b, err := json.Marshal(resp)
		if err != nil {
			err = errors.Wrap(err, errMarshalResp.Error())
			log.Printf("%+v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		} else {
			_, err = w.Write(b)
			if err != nil {
				err = errors.Wrap(err, errWriteResp.Error())
				log.Printf("%+v", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}
	})
}

// Path returns the pkg's configured path.
func Path() string {
	return c.Path
}
