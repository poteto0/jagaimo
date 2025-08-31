package html

import (
	"testing"

	"github.com/poteto0/jagaimo/core/renderer/dom"
	"github.com/poteto0/jagaimo/core/renderer/html/types"
	"github.com/stretchr/testify/assert"
)

func Test_createElementNode(t *testing.T) {
	// Act & Assert
	assert.IsType(
		t,
		&dom.Node{},
		createElementNode(
			"html",
			[]types.Attribute{},
		),
	)
}

func Test_createRune(t *testing.T) {
	// Act & Assert
	assert.IsType(
		t,
		&dom.Node{},
		createRune('a'),
	)
}

func Test_isAsciiAlphabetic(t *testing.T) {
	// Act
	result := isAsciiAlphabetic('a')

	// Assert
	assert.True(t, result)

	// Act
	result = isAsciiAlphabetic('A')

	// Assert
	assert.True(t, result)

	// Act
	result = isAsciiAlphabetic('<')

	// Assert
	assert.False(t, result)
}
