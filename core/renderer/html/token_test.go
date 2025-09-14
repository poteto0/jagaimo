package html

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewEOFToken(t *testing.T) {
	// Arrange
	token := newEOFToken()

	// Assert
	assert.Equal(t, EOF(1), token.EOF)
}

func TestNewRuneToken(t *testing.T) {
	// Arrange
	r := 'a'
	token := newRuneToken(r)

	// Assert
	assert.Equal(t, r, token.Rune)
}

func TestHtmlToken_IsStartTag(t *testing.T) {
	// Arrange
	token := HtmlToken{
		StartTag: &StartTag{},
	}

	// Act & Assert
	assert.True(t, token.IsStartTag())

	// Arrange
	token = HtmlToken{
		StartTag: nil,
	}

	// Act & Assert
	assert.False(t, token.IsStartTag())
}

func TestIsEnToken(t *testing.T) {
	// Arrange
	token := HtmlToken{
		EndTag: &EndTag{},
	}

	// Act & Assert
	assert.True(t, token.IsEndTag())

	// Arrange
	token = HtmlToken{
		EndTag: nil,
	}

	// Act & Assert
	assert.False(t, token.IsEndTag())
}

func TestHtmlToken_IsEOF(t *testing.T) {
	// Arrange
	token := HtmlToken{
		EOF: EOF(1),
	}

	// Act & Assert
	assert.True(t, token.IsEOF())

	// Arrange
	token = HtmlToken{
		EOF: EOF(0),
	}

	// Act & Assert
	assert.False(t, token.IsEOF())
}

func TestHtmlToken_IsRune(t *testing.T) {
	// Arrange
	token := HtmlToken{
		Rune: 'a',
	}

	// Act & Assert
	assert.True(t, token.IsRune())

	// Arrange
	token = HtmlToken{
		Rune: 0,
	}

	// Act & Assert
	assert.False(t, token.IsRune())
}
