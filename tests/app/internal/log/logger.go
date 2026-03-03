// Package log provides a thin wrapper around logrus with support for field censorship.
package log

import (
	"encoding/json"
	"fmt"
	"sync"

	elasticlogrus "bitbucket.org/csgot/helis-elasticlogrus"
	"github.com/sirupsen/logrus"
)

// DefaultChannel is the default logger channel name used by the application.
const DefaultChannel = "app"

var (
	loggerOnce     sync.Once
	loggerInstance *logrus.Logger
)

// SensitiveFields is a set of field names that should be censored in logs.
type SensitiveFields map[string]struct{}

// NewLogger creates a singleton logrus.Logger configured with censorship and optional hooks.
func NewLogger(logLevel, channel string, fields SensitiveFields, additionalHooks ...logrus.Hook) *logrus.Logger {
	loggerOnce.Do(func() {
		loggerInstance = elasticlogrus.NewCensoringLogger(logLevel, channel, fields)

		for _, h := range additionalHooks {
			loggerInstance.AddHook(h)
		}
	})

	return loggerInstance
}

// ParseSensitiveFields parses a JSON map of field names to booleans into a SensitiveFields set.
func ParseSensitiveFields(rawSensitiveFields string) (SensitiveFields, error) {
	fields := make(SensitiveFields)

	if err := json.Unmarshal([]byte(rawSensitiveFields), &fields); err != nil {
		return nil, fmt.Errorf("failed to unmarshal sensitive fields: %w", err)
	}

	return fields, nil
}
