package browser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewPage(t *testing.T) {
	// Act & Assert
	assert.IsType(t, &Page{}, NewPage())
}
