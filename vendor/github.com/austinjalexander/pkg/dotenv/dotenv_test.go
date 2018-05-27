package dotenv

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRun(t *testing.T) {
	tmpfilepath, envVars, err := setUp()
	if err != nil {
		t.Error(err)
	}
	defer func() {
		err := os.Remove(tmpfilepath)
		if err != nil {
			log.Printf("problem removing file %q: %s", tmpfilepath, err)
		}
		for k := range envVars {
			err := os.Unsetenv(k)
			if err != nil {
				log.Printf("problem unsetting env var%q: %s", k, err)
			}
		}
	}()
	// Check initial env-var values.
	for k := range envVars {
		assert.Empty(t, os.Getenv(k))
	}
	// Run Run().
	err = Run(tmpfilepath)
	if err != nil {
		t.Error(err)
	}
	// Check env vars again (post-Run).
	for k, v := range envVars {
		assert.Equal(t, os.Getenv(k), v)
	}
}

func setUp() (string, map[string]string, error) {
	envVars := map[string]string{
		"TEST_DB_HOST": "api-staging.abc123.us-west-2.rds.amazonaws.com",
		"TEST_DB_NAME": "api-staging",
		"TEST_DB_PASS": "notArealPa55word",
		"TEST_DB_PORT": "5432",
		"TEST_DB_USER": "api-staging",
	}
	var buf bytes.Buffer
	for k, v := range envVars {
		buf.WriteString(fmt.Sprintf("%s=%s\n", k, v))
	}
	tmpfile, err := ioutil.TempFile("", ".env-test")
	if err != nil {
		return "", nil, err
	}
	_, err = tmpfile.Write(buf.Bytes())
	if err != nil {
		return "", nil, err
	}
	err = tmpfile.Close()
	if err != nil {
		return "", nil, err
	}
	return tmpfile.Name(), envVars, nil
}
