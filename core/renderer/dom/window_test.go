package dom

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewWindow(t *testing.T) {
	// Act & Assert
	assert.IsType(t, &Window{}, NewWindow())
}

func TestWindow_Document(t *testing.T) {
	// Arrange
	window := NewWindow().(*Window)

	// Act & Assert
	assert.IsType(t, &Node{}, window.Document())
}
