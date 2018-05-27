// Package db wraps.
package db

import (
	"fmt"
	"sync"

	"github.com/jmoiron/sqlx"
	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
	"github.com/rubenv/sql-migrate"
	log "github.com/sirupsen/logrus"
	// Register postgres driver.
	_ "github.com/lib/pq"
)

const (
	appliedMigrationsStringFormat = `applied %d migration(s)`
	connStringFormat              = `host=%s dbname=%s password=%s port=%d user=%s sslmode=%s`
	driverName                    = `postgres`
	processName                   = `db`
	sslEnabled                    = `require`
	sslDisabled                   = `disable`
)

var (
	errDbConnect  = errors.New(`unable to connect to database`)
	errMigrateUp  = errors.New(`unable to migrate up`)
	errProcessEnv = errors.New(`unable to process environment`)

	c config
	// dba is a singleton database abstraction.
	dba  *sqlx.DB
	once sync.Once
)

// config represents the configuration necessary for this pkg.
type config struct {
	Host string `envconfig:"DB_HOST"`
	Name string `envconfig:"DB_NAME"`
	Pass string `envconfig:"DB_PASS"`
	Port int    `envconfig:"DB_PORT"`
	User string `envconfig:"DB_USER"`

	SslEnabled bool `envconfig:"DB_SSL_ENABLED"`

	MaxIdleConns int `envconfig:"DB_MAX_IDLE_CONNS"`
	MaxOpenConns int `envconfig:"DB_MAX_OPEN_CONNS"`

	MigrationsDir   string `envconfig:"DB_MIGRATIONS_DIR"`
	MigrationsRun   bool   `envconfig:"DB_MIGRATIONS_RUN"`
	MigrationsTable string `envconfig:"DB_MIGRATIONS_TABLE"`

	DriverName, ProcessName string
}

// DB returns a pointer to the configured database abstraction for this pkg.
func DB() *sqlx.DB {
	return dba
}

// DownMigrate down-migrates all migration files in migrationsDir.
func DownMigrate(migrationsDir string) error {
	n, err := migrate.Exec(dba.DB, driverName, &migrate.FileMigrationSource{Dir: migrationsDir}, migrate.Down)
	if err != nil {
		return err
	}
	log.Printf("down-migrated %d migrations", n)
	return nil
}

// Init configures the current pkg using env vars (via initialize).
// While using init would be nice, Init allows us to be more
// explicit in different environments and helps to enforce
// one-time initialization.
func Init() {
	once.Do(initialize)
}

func initialize() {
	if c == (config{}) && dba == nil {
		err := envconfig.Process(processName, &c)
		if err != nil {
			log.Fatal(errors.Wrap(err, errProcessEnv.Error()))
		}

		// sqlx.Connect() tests the ability to connect with a db.Ping().
		db, err := sqlx.Connect(driverName, connString())
		if err != nil {
			log.Fatal(errors.Wrap(err, errDbConnect.Error()))
		}
		dba = db

		dba.SetMaxIdleConns(c.MaxIdleConns)
		dba.SetMaxOpenConns(c.MaxOpenConns)

		if c.MigrationsRun {
			migrate.SetTable(c.MigrationsTable)
			migrations := &migrate.FileMigrationSource{Dir: c.MigrationsDir}
			n, err := migrate.Exec(dba.DB, driverName, migrations, migrate.Up)
			if err != nil {
				log.Fatal(errors.Wrap(err, errMigrateUp.Error()))
			}
			log.Printf(appliedMigrationsStringFormat, n)
		}
	}
}

func connString() string {
	return fmt.Sprintf(connStringFormat,
		c.Host,
		c.Name,
		c.Pass,
		c.Port,
		c.User,
		sslMode(),
	)
}

func sslMode() string {
	if c.SslEnabled {
		return sslEnabled
	}
	return sslDisabled
}
