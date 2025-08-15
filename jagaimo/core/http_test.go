package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewHeader(t *testing.T) {
	// Act & Assert
	assert.Equal(t, Header{"name", "value"}, NewHeader("name", "value"))
}
