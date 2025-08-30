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

			// Act
			next, _ := parser.parseInitial(newRuneToken('<'))

			// Assert
			assert.Equal(t, &HtmlToken{
				StartTag: &StartTag{
					Tag:           "html",
					IsSelfClosing: false,
					Attributes:    []types.Attribute{},
				},
			}, next)
		})

		t.Run("if not rune token, change mode BeforeHtml & return nil", func(t *testing.T) {
			// Arrange
			parser := NewHtmlParser(
				NewHtmlTokenizer(""),
			).(*HtmlParser)

			// Act
			next, _ := parser.parseInitial(newEOFToken())

			// Act & Assert
			assert.Nil(t, next)
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

func TestHtmlParser_parseBeforeHtml(t *testing.T) {
	t.Run("normal case", func(t *testing.T) {
		init := func() *HtmlParser {
			parser := NewHtmlParser(NewHtmlTokenizer(
				" <html><head></head><body></body></html>",
			)).(*HtmlParser)
			parser.mode = BeforeHtml
			return parser
		}

		t.Run("if rune token of ' ' or \n , return next token", func(t *testing.T) {
			// Arrange
			parser := init()
			token := parser.t.Next()

			// Act
			next, _ := parser.parseBeforeHtml(token)
			assert.Equal(t, &HtmlToken{
				StartTag: &StartTag{
					Tag:           "html",
					IsSelfClosing: false,
					Attributes:    []types.Attribute{},
				},
			}, next)
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
			parser.mode = Initial

			// Act
			parser.parseBeforeHtml(nil)

			// Assert
			assert.Error(t, err)
		})
	})
}
