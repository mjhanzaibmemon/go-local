package cacheadapter

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDomainId(t *testing.T) {
	assert.Equal(t, "testNamespace/testKey", makeID("testNamespace", "testKey"))
}
