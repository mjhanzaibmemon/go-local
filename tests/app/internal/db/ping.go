package db

import (
	"errors"

	"gorm.io/gorm"
)

// Ping checks whether the underlying DB connection is alive.
func Ping(db *gorm.DB) error {
	if pinger, ok := db.ConnPool.(interface{ Ping() error }); ok {
		return pinger.Ping()
	}

	return errors.New("DB connection does not implement Pinger interface")
}
