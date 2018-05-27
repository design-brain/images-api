package db

import (
	"flag"
	"io/ioutil"
	"os"
	"testing"

	"github.com/rubenv/sql-migrate"

	"github.com/stretchr/testify/assert"
)

var (
	// Tests using this flag expect relevant env vars to be set
	// and the relevant services to be available.
	integration = flag.Bool("integration", false, "run integration tests")
)

func TestInit(t *testing.T) {
	flag.Parse()
	if !*integration {
		t.SkipNow()
	}
	err := os.Setenv("DB_MIGRATIONS_DIR", "../../_migrations")
	if err != nil {
		t.Error(err)
	}

	// Initialize pkg configuration.
	Init()

	// Check migration status against migration files.
	migrations, err := migrate.GetMigrationRecords(DB().DB, driverName)
	if err != nil {
		t.Error(err)
	}
	migrationsDir := os.Getenv("DB_MIGRATIONS_DIR")
	migrationDirEntries, err := ioutil.ReadDir(migrationsDir)
	if err != nil {
		t.Error(err)
	}
	var migrationFileCount int
	for _, f := range migrationDirEntries {
		if !f.IsDir() {
			migrationFileCount++
		}
	}
	assert.Equal(t, len(migrations), migrationFileCount)

	// Migrate down.
	err = DownMigrate(migrationsDir)
	if err != nil {
		t.Error(err)
	}
}

func TestConnString(t *testing.T) {
	t.Parallel()

	c = config{
		Host:       "localhost",
		Name:       "dbname",
		Pass:       "secret",
		Port:       8080,
		User:       "dbuser",
		SslEnabled: true,
	}

	const expectedConnStr = `host=localhost dbname=dbname password=secret port=8080 user=dbuser sslmode=require`
	assert.Equal(t, connString(), expectedConnStr)
}

func TestSslMode(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		cfg             config
		expectedSslMode string
	}{
		{
			cfg: config{
				SslEnabled: false,
			},
			expectedSslMode: sslDisabled,
		},
		{
			cfg: config{
				SslEnabled: true,
			},
			expectedSslMode: sslEnabled,
		},
	}
	for _, tc := range testCases {
		c = tc.cfg
		assert.Equal(t, sslMode(), tc.expectedSslMode)
	}
}
