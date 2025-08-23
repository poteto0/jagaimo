package html

import (
	"errors"
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
	expected := newRuneToken('h')

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
	assert.Equal(t, uint(1), tokenizer.Pos)
}

func TestHtmlTokenizer_reConsumeInput(t *testing.T) {
	// Arrange
	tokenizer := NewHtmlTokenizer("<html><body><h1>Hello, World!</h1></body></html>").(*HtmlTokenizer)
	tokenizer.Pos = 1

	// Act
	result := tokenizer.reConsumeInput()

	// Assert
	assert.Equal(t, '<', result)
	assert.Equal(t, uint(1), tokenizer.Pos)
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

func TestHtmlTokenizer_createTag(t *testing.T) {
	tests := []struct {
		name             string
		isStartTagToken  bool
		expectedStartTag *StartTag
		expectedEndTag   *EndTag
	}{
		{
			name:            "Create start tag token",
			isStartTagToken: true,
			expectedStartTag: &StartTag{
				Tag:           "",
				IsSelfClosing: false,
				Attributes:    []Attribute{},
			},
			expectedEndTag: nil,
		},
		{
			name:             "Create end tag token",
			isStartTagToken:  false,
			expectedStartTag: nil,
			expectedEndTag: &EndTag{
				Tag: "",
			},
		},
	}

	for _, it := range tests {
		t.Run(it.name, func(t *testing.T) {
			// Arrange
			tokenizer := NewHtmlTokenizer("").(*HtmlTokenizer)

			// Act
			tokenizer.createTag(it.isStartTagToken)

			// Assert
			if it.isStartTagToken {
				assert.Equal(t, it.expectedStartTag, tokenizer.LatestToken.StartTag)
			} else {
				assert.Equal(t, it.expectedEndTag, tokenizer.LatestToken.EndTag)
			}
		})
	}
}

func TestHtmlTokenizer_appendTagName(t *testing.T) {
	t.Run("normal case", func(t *testing.T) {
		tests := []struct {
			name        string
			r           rune
			latestToken *HtmlToken
			expected    string
		}{
			{
				name: "append to start tag",
				r:    'a',
				latestToken: &HtmlToken{
					StartTag: &StartTag{
						Tag:           "",
						IsSelfClosing: false,
						Attributes:    []Attribute{},
					},
				},
				expected: "a",
			},
			{
				name: "append to end tag",
				r:    'a',
				latestToken: &HtmlToken{
					EndTag: &EndTag{
						Tag: "",
					},
				},
				expected: "a",
			},
		}

		// Arrange
		tokenizer := NewHtmlTokenizer("").(*HtmlTokenizer)

		for _, it := range tests {
			t.Run(it.name, func(t *testing.T) {
				tokenizer.LatestToken = it.latestToken

				// Act
				tokenizer.appendTagName(it.r)

				// Assert
				if it.latestToken.IsStartTag() {
					assert.Equal(t, it.expected, tokenizer.LatestToken.StartTag.Tag)
				}

				if it.latestToken.IsEndTag() {
					assert.Equal(t, it.expected, tokenizer.LatestToken.EndTag.Tag)
				}
			})
		}
	})

	t.Run("panic case", func(t *testing.T) {
		tests := []struct {
			name        string
			latestToken *HtmlToken
		}{
			{
				name:        "nil latest token",
				latestToken: nil,
			},
			{
				name:        "unexpected eof token",
				latestToken: newEOFToken(),
			},
		}

		// Arrange
		tokenizer := NewHtmlTokenizer("").(*HtmlTokenizer)

		for _, it := range tests {
			t.Run(it.name, func(t *testing.T) {
				var err error
				defer func() {
					if r := recover(); r != nil {
						err = errors.New("panic")
					}
				}()

				// Arrange
				tokenizer.LatestToken = it.latestToken

				// Act
				tokenizer.appendTagName('a')

				// Assert
				assert.Error(t, err)
			})
		}
	})
}

func TestHtmlTokenizer_takeLastToken(t *testing.T) {
	t.Run("normal case, return latest token & set latest token to nil", func(t *testing.T) {
		tests := []struct {
			name        string
			latestToken *HtmlToken
			expected    *HtmlToken
		}{
			{
				name:        "take last start tag token",
				latestToken: newRuneToken('a'),
				expected:    newRuneToken('a'),
			},
		}

		// Arrange
		tokenizer := NewHtmlTokenizer("").(*HtmlTokenizer)

		for _, it := range tests {
			t.Run(it.name, func(t *testing.T) {
				// Arrange
				tokenizer.LatestToken = it.latestToken

				// Act
				result := tokenizer.takeLastToken()

				// Assert
				assert.Equal(t, it.expected, result)
				assert.Nil(t, tokenizer.LatestToken)
			})
		}
	})

	t.Run("panic case", func(t *testing.T) {
		tests := []struct {
			name        string
			latestToken *HtmlToken
		}{
			{
				name:        "nil latest token",
				latestToken: nil,
			},
		}

		// Arrange
		tokenizer := NewHtmlTokenizer("").(*HtmlTokenizer)

		for _, it := range tests {
			t.Run(it.name, func(t *testing.T) {
				var err error
				defer func() {
					if r := recover(); r != nil {
						err = errors.New("panic")
					}
				}()

				// Arrange
				tokenizer.LatestToken = it.latestToken

				// Act
				tokenizer.takeLastToken()

				// Assert
				assert.Error(t, err)
			})
		}
	})
}

func TestHtmlTokenizer_startNewAttribute(t *testing.T) {
	t.Run("normal case, append attribute", func(t *testing.T) {
		tests := []struct {
			name        string
			latestToken *HtmlToken
			expected    []Attribute
		}{
			{
				name: "append to start tag",
				latestToken: &HtmlToken{
					StartTag: &StartTag{},
				},
				expected: []Attribute{
					{},
				},
			},
		}

		// Arrange
		tokenizer := NewHtmlTokenizer("").(*HtmlTokenizer)

		for _, it := range tests {
			t.Run(it.name, func(t *testing.T) {
				// Arrange
				tokenizer.LatestToken = it.latestToken

				// Act
				tokenizer.startNewAttribute()

				// Assert
				assert.Equal(t, it.expected, tokenizer.LatestToken.StartTag.Attributes)
			})
		}
	})

	t.Run("panic case", func(t *testing.T) {
		tests := []struct {
			name        string
			latestToken *HtmlToken
		}{
			{
				name:        "nil latest token",
				latestToken: nil,
			},
			{
				name:        "unexpected eof token",
				latestToken: newEOFToken(),
			},
		}

		// Arrange
		tokenizer := NewHtmlTokenizer("").(*HtmlTokenizer)

		for _, it := range tests {
			t.Run(it.name, func(t *testing.T) {
				var err error
				defer func() {
					if r := recover(); r != nil {
						err = errors.New("panic")
					}
				}()

				// Arrange
				tokenizer.LatestToken = it.latestToken

				// Act
				tokenizer.startNewAttribute()

				// Assert
				assert.Error(t, err)
			})
		}
	})
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
