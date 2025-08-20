package html_test

import (
	"testing"

	"github.com/poteto0/jagaimo/core/renderer/html"
	"github.com/stretchr/testify/assert"
)

func TestNewHtmlTokenizer(t *testing.T) {
	// Arrange
	input := "<html><body><h1>Hello, World!</h1></body></html>"

	// Act
	tokenizer := html.NewHtmlTokenizer(input).(*html.HtmlTokenizer)

	// Assert
	assert.NotNil(t, tokenizer)
	assert.Equal(t, html.Data, tokenizer.State)
	assert.Equal(t, uint(0), tokenizer.Pos)
	assert.False(t, tokenizer.ReConsume)
	assert.Nil(t, tokenizer.LatestToken)
	assert.Equal(t, input, string(tokenizer.Input))
	assert.Equal(t, "", tokenizer.Buf)
}
