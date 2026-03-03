package log

import (
	"github.com/stretchr/testify/assert"
	"io"
	"testing"
)

func TestTestLogger(t *testing.T) {
	assert.Equal(t, io.Discard, TestLogger().Out)
}
