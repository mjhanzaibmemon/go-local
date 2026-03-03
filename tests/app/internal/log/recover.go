package log

import (
	"runtime/debug"

	elasticlogrus "bitbucket.org/csgot/helis-elasticlogrus"
	"github.com/sirupsen/logrus"
)

//nolint:revive
func RecoverFromPanic() {
	if err := recover(); err != nil {
		logger := elasticlogrus.NewCensoringLogger(logrus.InfoLevel.String(), "panic", nil)
		logger.
			WithField(logrus.ErrorKey, err).
			WithField("trace", string(debug.Stack())).
			Fatal("panic occurred")
	}
}
