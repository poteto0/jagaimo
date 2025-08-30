package html

import (
	"errors"
	"testing"

	"github.com/poteto0/jagaimo/core/renderer/html/types"
	"github.com/stretchr/testify/assert"
)

func TestNewHtmlParser(t *testing.T) {
	// Act & Assert
	assert.IsType(t, &HtmlParser{}, NewHtmlParser(nil))
}

func TestHtmlParser_parseInitial(t *testing.T) {
	t.Run("normal case", func(t *testing.T) {
		t.Run("if rune token, return next token", func(t *testing.T) {
			// Arrange
			parser := NewHtmlParser(NewHtmlTokenizer(
				"<html><head></head><body></body></html>",
			)).(*HtmlParser)

			// Act & Assert
			assert.Equal(t, &HtmlToken{
				StartTag: &StartTag{
					Tag:           "html",
					IsSelfClosing: false,
					Attributes:    []types.Attribute{},
				},
			}, parser.parseInitial(newRuneToken('<')))
		})

		t.Run("if not rune token, change mode BeforeHtml & return nil", func(t *testing.T) {
			// Arrange
			parser := NewHtmlParser(
				NewHtmlTokenizer(""),
			).(*HtmlParser)

			// Act & Assert
			assert.Nil(t, parser.parseInitial(newEOFToken()))
			assert.Equal(t, BeforeHtml, parser.mode)
		})
	})

	t.Run("panic case", func(t *testing.T) {
		t.Run("the mode is not Initial", func(t *testing.T) {
			var err error
			defer func() {
				if r := recover(); r != nil {
					err = errors.New("panic")
				}
			}()

			// Arrange
			parser := NewHtmlParser(nil).(*HtmlParser)
			parser.mode = BeforeHtml

			// Act
			parser.parseInitial(nil)

			// Assert
			assert.Error(t, err)
		})
	})
}
