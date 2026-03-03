package log

import (
	"io"

	"github.com/sirupsen/logrus"
)

// TestLogger returns logger which supposed to be used in tests
func TestLogger() *logrus.Logger {
	testLogger := logrus.New()
	testLogger.SetOutput(io.Discard)

	return testLogger
}
