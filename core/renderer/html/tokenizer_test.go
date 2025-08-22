package html

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewHtmlTokenizer(t *testing.T) {
	// Arrange
	input := "<html><body><h1>Hello, World!</h1></body></html>"

	// Act
	tokenizer := NewHtmlTokenizer(input).(*HtmlTokenizer)

	// Assert
	assert.NotNil(t, tokenizer)
	assert.Equal(t, Data, tokenizer.State)
	assert.Equal(t, uint(0), tokenizer.Pos)
	assert.False(t, tokenizer.ReConsume)
	assert.Nil(t, tokenizer.LatestToken)
	assert.Equal(t, input, string(tokenizer.Input))
	assert.Equal(t, "", tokenizer.Buf)
}

func TestHtmlTokenizer_Iter(t *testing.T) {
	// Arrange
	tokenizer := NewHtmlTokenizer("html")
	expected := &HtmlToken{
		StartTag: StartTag{},
		EndTag:   EndTag{},
		EOF:      EOF(0),
		Rune:     'h',
	}

	// Act
	for token := range tokenizer.Iter() {
		assert.Equal(t, expected, token)
	}
}

func TestHtmlTokenizer_consumeNextInput(t *testing.T) {
	// Arrange
	tokenizer := NewHtmlTokenizer("<html><body><h1>Hello, World!</h1></body></html>").(*HtmlTokenizer)

	// Act
	result := tokenizer.consumeNextInput()

	// Assert
	assert.Equal(t, '<', result)
}

func TestHtmlTokenizer_isEOF(t *testing.T) {
	// Arrange
	tokenizer := NewHtmlTokenizer("a").(*HtmlTokenizer)

	// Act & Assert
	assert.False(t, tokenizer.isEOF())

	// Arrange
	tokenizer.Pos = 2

	// Act & Assert
	assert.True(t, tokenizer.isEOF())
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
