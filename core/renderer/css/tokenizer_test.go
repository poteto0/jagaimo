package css

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewCssTokenizer(t *testing.T) {
	// Act & Assert
	assert.IsType(t, &CssTokenizer{}, NewCssTokenizer("").(*CssTokenizer))
}

func TestCssTokenizer_decideToken(t *testing.T) {
	t.Run("decide CssToken by rune", func(t *testing.T) {
		// Arrange
		tokenizer := NewCssTokenizer("").(*CssTokenizer)
		tests := []struct {
			input           rune
			expected        CssToken
			expectedSkipped bool
		}{
			{
				input:           '(',
				expected:        OpenParenthesis,
				expectedSkipped: false,
			},
			{
				input:           ')',
				expected:        CloseParenthesis,
				expectedSkipped: false,
			},
			{
				input:           ',',
				expected:        DelimComma,
				expectedSkipped: false,
			},
			{
				input:           '.',
				expected:        DelimPeriod,
				expectedSkipped: false,
			},
			{
				input:           ':',
				expected:        Colon,
				expectedSkipped: false,
			},
			{
				input:           ';',
				expected:        Semicolon,
				expectedSkipped: false,
			},
			{
				input:           '{',
				expected:        OpenCurly,
				expectedSkipped: false,
			},
			{
				input:           '}',
				expected:        CloseCurly,
				expectedSkipped: false,
			},
			{
				input:           ' ',
				expected:        NilToken,
				expectedSkipped: true,
			},
			{
				input:           '\n',
				expected:        NilToken,
				expectedSkipped: true,
			},
			{
				input:           'k',
				expected:        NilToken,
				expectedSkipped: false,
			},
		}

		for _, it := range tests {
			// Act
			actual, isSkip := tokenizer.decideToken(it.input)

			// Assert
			assert.Equal(t, it.expected, actual)
			assert.Equal(t, it.expectedSkipped, isSkip)
		}
	})
}
