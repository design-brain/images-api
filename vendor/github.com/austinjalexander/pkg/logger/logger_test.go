package logger

import (
	"os"
	"testing"

	"github.com/sirupsen/logrus"
)

func TestInit(t *testing.T) {
	testSetup(t)
	defer testTeardown(t)

	Init()

	lvl := logrus.GetLevel()
	if lvl != logrus.DebugLevel {
		t.Errorf("expected lvl to be Debug, got %s", lvl)
	}
}

func testEnv() ([][]string, error) {
	return [][]string{
		{"LOG_LEVEL", "debug"},
	}, nil
}

func testSetup(t *testing.T) {
	env, err := testEnv()
	if err != nil {
		t.Error(err)
	} else {
		for _, e := range env {
			err := os.Setenv(e[0], e[1])
			if err != nil {
				t.Error(err)
			}
		}
	}
}

func testTeardown(t *testing.T) {
	env, err := testEnv()
	if err != nil {
		t.Error(err)
	} else {
		for _, e := range env {
			err := os.Unsetenv(e[0])
			if err != nil {
				t.Error(err)
			}
		}
	}
}
