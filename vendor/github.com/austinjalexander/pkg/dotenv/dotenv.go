package dotenv

import (
	"bufio"
	"log"
	"os"
	"strings"

	"github.com/pkg/errors"
)

// Run looks for the passed filepath and sets the environment variables contained
// within.
func Run(filepath string) error {
	// Open file.
	f, err := os.Open(filepath)
	if err != nil {
		return errors.Wrapf(err, "unable to open file %s", filepath)
	}
	defer func() {
		err := f.Close()
		if err != nil {
			log.Printf("%+v", errors.Wrapf(err, "closing file %s", filepath))
		}
	}()

	// Read and process file.
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" && string(line[0]) != "#" {
			keyValuePair := strings.SplitN(line, "=", 2)
			err := os.Setenv(keyValuePair[0], keyValuePair[1])
			if err != nil {
				return errors.Wrapf(err, "unable to set env var pair k=%s, v=%s", keyValuePair[0], keyValuePair[1])
			}
		}
	}
	err = scanner.Err()
	if err != nil {
		return errors.Wrapf(err, "unable to process file %s", filepath)
	}
	return nil
}
