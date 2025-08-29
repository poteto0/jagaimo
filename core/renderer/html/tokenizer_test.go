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
	t.Run("empty string return nil", func(t *testing.T) {
		// Arrange
		tokenizer := NewHtmlTokenizer("").(*HtmlTokenizer)

		// Act & Assert
		assert.Nil(t, tokenizer.Next())
	})

	t.Run("iter case", func(t *testing.T) {

		tests := []struct {
			name     string
			input    string
			expected []*HtmlToken
		}{
			{
				name:  "start & end tag",
				input: "<body></body>",
				expected: []*HtmlToken{
					{
						StartTag: &StartTag{
							Tag:           "body",
							IsSelfClosing: false,
							Attributes:    []Attribute{},
						},
					},
					{
						EndTag: &EndTag{
							Tag: "body",
						},
					},
				},
			},
			{
				name:  "attributes",
				input: "<p class=\"A\" id='B' foo=bar fizz=buzz></p>",
				expected: []*HtmlToken{
					{
						StartTag: &StartTag{
							Tag: "p",
							Attributes: []Attribute{
								{
									name:  "class",
									value: "A",
								},
								{
									name:  "id",
									value: "B",
								},
								{
									name:  "foo",
									value: "bar",
								},
								{
									name:  "fizz",
									value: "buzz",
								},
							},
							IsSelfClosing: false,
						},
					},
					{
						EndTag: &EndTag{
							Tag: "p",
						},
					},
				},
			},
			{
				name:  "quoted attribute",
				input: "<div id=\"div\"><p class=\" A\" id=\"BC\"/></div>",
				expected: []*HtmlToken{
					{
						StartTag: &StartTag{
							Tag: "div",
							Attributes: []Attribute{
								{
									name:  "id",
									value: "div",
								},
							},
							IsSelfClosing: false,
						},
					},
					{
						StartTag: &StartTag{
							Tag: "p",
							Attributes: []Attribute{
								{
									name:  "class",
									value: " A",
								},
								{
									name:  "id",
									value: "BC",
								},
							},
							IsSelfClosing: true,
						},
					},
					{
						EndTag: &EndTag{
							Tag: "div",
						},
					},
				},
			},
			{
				name:  "self closing & empty tag case",
				input: "<img />",
				expected: []*HtmlToken{
					{
						StartTag: &StartTag{
							Tag:           "img",
							IsSelfClosing: true,
							Attributes:    []Attribute{},
						},
					},
				},
			},
			{
				name:  "script tag",
				input: "<script> console.log(\"hello\")</script>",
				expected: []*HtmlToken{
					{
						StartTag: &StartTag{
							Tag:           "script",
							IsSelfClosing: false,
							Attributes:    []Attribute{},
						},
					},
					newRuneToken(' '),
					newRuneToken('c'),
					newRuneToken('o'),
					newRuneToken('n'),
					newRuneToken('s'),
					newRuneToken('o'),
					newRuneToken('l'),
					newRuneToken('e'),
					newRuneToken('.'),
					newRuneToken('l'),
					newRuneToken('o'),
					newRuneToken('g'),
					newRuneToken('('),
					newRuneToken('"'),
					newRuneToken('h'),
					newRuneToken('e'),
					newRuneToken('l'),
					newRuneToken('l'),
					newRuneToken('o'),
					newRuneToken('"'),
					newRuneToken(')'),
					{
						EndTag: &EndTag{
							Tag: "script",
						},
					},
				},
			},
		}

		for _, it := range tests {
			t.Run(it.name, func(t *testing.T) {
				// Arrange
				tokenizer := NewHtmlTokenizer(it.input).(*HtmlTokenizer)

				for _, e := range it.expected {
					// Act
					token := tokenizer.Next()

					// Assert
					assert.Equal(t, e, token)
				}
			})
		}
	})
}

func TestHtmlTokenizer_tokenizeData(t *testing.T) {
	t.Run("normal case", func(t *testing.T) {
		init := func() *HtmlTokenizer {
			tokenizer := NewHtmlTokenizer("").(*HtmlTokenizer)
			tokenizer.State = Data
			return tokenizer
		}

		t.Run("if EOF, return EOF token", func(t *testing.T) {
			// Arrange
			tokenizer := init()
			tokenizer.Pos = 1

			// Act
			result := tokenizer.tokenizeData(rune(0))

			// Assert
			assert.Equal(t, newEOFToken(), result)
		})

		t.Run("if <, change State TagOpen", func(t *testing.T) {
			// Arrange
			tokenizer := init()

			// Act & Assert
			assert.Nil(t, tokenizer.tokenizeData('<'))
			assert.Equal(t, TagOpen, tokenizer.State)
		})

		t.Run("if other input, return runeToken", func(t *testing.T) {
			// Arrange
			tokenizer := init()

			// Act
			result := tokenizer.tokenizeData('a')

			// Assert
			assert.Equal(t, newRuneToken('a'), result)
		})
	})

	t.Run("panic case", func(t *testing.T) {
		tests := []struct {
			name  string
			state State
			r     rune
		}{
			{
				name:  "the state is not Data",
				state: TagOpen,
				r:     'a',
			},
		}

		for _, it := range tests {
			t.Run(it.name, func(t *testing.T) {
				var err error
				defer func() {
					if r := recover(); r != nil {
						err = errors.New("panic")
					}
				}()

				// Arrange
				tokenizer := NewHtmlTokenizer("").(*HtmlTokenizer)
				tokenizer.State = it.state

				// Act
				tokenizer.tokenizeData(it.r)

				// Assert
				assert.Error(t, err)
			})
		}
	})
}

func TestHtmlTokenizer_tokenizeTagOpen(t *testing.T) {
	t.Run("normal case", func(t *testing.T) {
		init := func() *HtmlTokenizer {
			tokenizer := NewHtmlTokenizer("").(*HtmlTokenizer)
			tokenizer.State = TagOpen
			return tokenizer
		}

		t.Run("if EOF, return newEOFToken", func(t *testing.T) {
			// Arrange
			tokenizer := init()
			tokenizer.Pos = 1

			// Act
			result := tokenizer.tokenizeTagOpen(rune(0))

			// Assert
			assert.Equal(t, newEOFToken(), result)
		})

		t.Run("if /, change State EndTagOpen", func(t *testing.T) {
			// Arrange
			tokenizer := init()

			// Act & Assert
			assert.Nil(t, tokenizer.tokenizeTagOpen('/'))
			assert.Equal(t, EndTagOpen, tokenizer.State)
		})

		t.Run("if a-zA-Z, set re-consume & change State TagName & create StartTag", func(t *testing.T) {
			// Arrange
			tokenizer := init()

			// Act & Assert
			assert.Nil(t, tokenizer.tokenizeTagOpen('a'))
			assert.Equal(t, TagName, tokenizer.State)
			assert.True(t, tokenizer.ReConsume)
			assert.NotNil(t, tokenizer.LatestToken)
		})

		t.Run("others input, set re-consume & change State Data", func(t *testing.T) {
			// Arrange
			tokenizer := init()

			// Act & Assert
			assert.Nil(t, tokenizer.tokenizeTagOpen('@'))
			assert.Equal(t, tokenizer.State, Data)
			assert.True(t, tokenizer.ReConsume)
		})
	})

	t.Run("panic case", func(t *testing.T) {
		tests := []struct {
			name  string
			state State
			r     rune
		}{
			{
				name:  "the state is not TagOpen",
				state: Data,
				r:     'a',
			},
		}

		for _, it := range tests {
			t.Run(it.name, func(t *testing.T) {
				var err error
				defer func() {
					if r := recover(); r != nil {
						err = errors.New("panic")
					}
				}()

				// Arrange
				tokenizer := NewHtmlTokenizer("").(*HtmlTokenizer)
				tokenizer.State = it.state

				// Act
				tokenizer.tokenizeTagOpen(it.r)

				// Assert
				assert.Error(t, err)
			})
		}
	})
}

func TestHtmlTokenizer_tokenizeEndTagOpen(t *testing.T) {
	t.Run("normal case", func(t *testing.T) {
		init := func() *HtmlTokenizer {
			tokenizer := NewHtmlTokenizer("").(*HtmlTokenizer)
			tokenizer.State = EndTagOpen
			return tokenizer
		}

		t.Run("if EOF, return newEOFToken", func(t *testing.T) {
			// Arrange
			tokenizer := init()
			tokenizer.Pos = 1

			// Act
			result := tokenizer.tokenizeEndTagOpen(rune(0))

			// Assert
			assert.Equal(t, newEOFToken(), result)
		})

		t.Run("if a-zA-Z, set re-consume & change State TagName & create EndTag", func(t *testing.T) {
			// Arrange
			tokenizer := init()

			// Act & Assert
			assert.Nil(t, tokenizer.tokenizeEndTagOpen('a'))
			assert.Equal(t, TagName, tokenizer.State)
			assert.NotNil(t, tokenizer.LatestToken)
		})

		t.Run("if others input, do nothing", func(t *testing.T) {
			// Arrange
			tokenizer := init()

			// Act & Assert
			assert.Nil(t, tokenizer.tokenizeEndTagOpen('/'))
			assert.Equal(t, EndTagOpen, tokenizer.State)
		})
	})

	t.Run("panic case", func(t *testing.T) {
		tests := []struct {
			name  string
			state State
			r     rune
		}{
			{
				name:  "the state is not EndTagOpen",
				state: Data,
				r:     'a',
			},
		}

		for _, it := range tests {
			t.Run(it.name, func(t *testing.T) {
				var err error
				defer func() {
					if r := recover(); r != nil {
						err = errors.New("panic")
					}
				}()

				// Arrange
				tokenizer := NewHtmlTokenizer("").(*HtmlTokenizer)
				tokenizer.State = it.state

				// Act
				tokenizer.tokenizeEndTagOpen(it.r)

				// Assert
				assert.Error(t, err)
			})
		}
	})
}

func TestHtmlTokenizer_tokenizeTagName(t *testing.T) {
	t.Run("normal case", func(t *testing.T) {
		init := func() *HtmlTokenizer {
			tokenizer := NewHtmlTokenizer("").(*HtmlTokenizer)
			tokenizer.State = TagName
			return tokenizer
		}

		t.Run("if ' ', change State BeforeAttributeName", func(t *testing.T) {
			// Arrange
			tokenizer := init()

			// Act & Assert
			assert.Nil(t, tokenizer.tokenizeTagName(' '))
			assert.Equal(t, BeforeAttributeName, tokenizer.State)
		})

		t.Run("if /, change State SelfClosingStartTag", func(t *testing.T) {
			// Arrange
			tokenizer := init()

			// Act & Assert
			assert.Nil(t, tokenizer.tokenizeTagName('/'))
			assert.Equal(t, SelfClosingStartTag, tokenizer.State)
		})

		t.Run("if >, change State Data & return latestToken", func(t *testing.T) {
			// Arrange
			expected := newRuneToken('a')
			tokenizer := init()
			tokenizer.LatestToken = expected

			// Act
			result := tokenizer.tokenizeTagName('>')

			// Assert
			assert.Equal(t, Data, tokenizer.State)
			assert.Equal(t, expected, result)
		})

		t.Run("if EOF, return newEOFToken", func(t *testing.T) {
			// Arrange
			tokenizer := init()
			tokenizer.Pos = 1

			// Act & Assert
			assert.Equal(t, newEOFToken(), tokenizer.tokenizeTagName(rune(0)))
		})

		t.Run("if a-zA-Z, append to the tag name", func(t *testing.T) {
			// Arrange
			tokenizer := init()
			tokenizer.LatestToken = &HtmlToken{
				StartTag: &StartTag{
					Tag: "",
				},
			}

			// Act & Assert
			assert.Nil(t, tokenizer.tokenizeTagName('a'))
			assert.Equal(t, "a", tokenizer.LatestToken.StartTag.Tag)
		})

		t.Run("if others input, append to the tag name", func(t *testing.T) {
			// Arrange
			tokenizer := init()
			tokenizer.LatestToken = &HtmlToken{
				StartTag: &StartTag{
					Tag: "",
				},
			}

			// Act & Assert
			assert.Nil(t, tokenizer.tokenizeTagName('@'))
			assert.Equal(t, "@", tokenizer.LatestToken.StartTag.Tag)
		})
	})

	t.Run("panic case", func(t *testing.T) {
		tests := []struct {
			name  string
			state State
			r     rune
		}{
			{
				name:  "the state is not TagName",
				state: Data,
				r:     'a',
			},
		}

		for _, it := range tests {
			t.Run(it.name, func(t *testing.T) {
				var err error
				defer func() {
					if r := recover(); r != nil {
						err = errors.New("panic")
					}
				}()

				// Arrange
				tokenizer := NewHtmlTokenizer("").(*HtmlTokenizer)
				tokenizer.State = it.state

				// Act
				tokenizer.tokenizeTagName(it.r)

				// Assert
				assert.Error(t, err)
			})
		}
	})
}

func TestHtmlTokenizer_tokenizeBeforeAttributeName(t *testing.T) {
	t.Run("normal case", func(t *testing.T) {
		init := func() *HtmlTokenizer {
			tokenizer := NewHtmlTokenizer("").(*HtmlTokenizer)
			tokenizer.State = BeforeAttributeName
			return tokenizer
		}

		t.Run("if / | > , set re-consume & change State AfterAttributeName", func(t *testing.T) {
			tests := []struct {
				r rune
			}{{r: '/'}, {r: '>'}}

			for _, it := range tests {
				// Arrange
				tokenizer := init()

				// Act & Assert
				assert.Nil(t, tokenizer.tokenizeBeforeAttributeName(it.r))
				assert.Equal(t, tokenizer.State, AfterAttributeName)
			}
		})

		t.Run("if EOF, re-consume & change State AfterAttributeName", func(t *testing.T) {
			// Arrange
			tokenizer := init()
			tokenizer.Pos = 1

			// Act & Assert
			assert.Nil(t, tokenizer.tokenizeBeforeAttributeName(rune(0)))
			assert.Equal(t, tokenizer.State, AfterAttributeName)
			assert.True(t, tokenizer.ReConsume)
		})

		t.Run("if others input, set re-consume & change State AttributeName & StartNewAttribute", func(t *testing.T) {
			// Arrange
			tokenizer := init()
			tokenizer.LatestToken = &HtmlToken{
				StartTag: &StartTag{
					Attributes: []Attribute{},
				},
			}

			// Act & Assert
			assert.Nil(t, tokenizer.tokenizeBeforeAttributeName('a'))
			assert.Equal(t, tokenizer.State, AttributeName)
			assert.True(t, tokenizer.ReConsume)
			assert.Equal(t, len(tokenizer.LatestToken.StartTag.Attributes), 1)
		})
	})

	t.Run("panic case", func(t *testing.T) {
		tests := []struct {
			name  string
			state State
			r     rune
		}{
			{
				name:  "the state is not BeforeAttributeName",
				state: Data,
				r:     'a',
			},
		}

		for _, it := range tests {
			t.Run(it.name, func(t *testing.T) {
				var err error
				defer func() {
					if r := recover(); r != nil {
						err = errors.New("panic")
					}
				}()

				// Arrange
				tokenizer := NewHtmlTokenizer("").(*HtmlTokenizer)
				tokenizer.State = it.state

				// Act
				tokenizer.tokenizeBeforeAttributeName(it.r)

				// Assert
				assert.Error(t, err)
			})
		}
	})
}

func TestHtmlTokenizer_tokenizeAttributeName(t *testing.T) {
	t.Run("normal case", func(t *testing.T) {
		init := func() *HtmlTokenizer {
			tokenizer := NewHtmlTokenizer("").(*HtmlTokenizer)
			tokenizer.State = AttributeName
			return tokenizer
		}

		t.Run("if / | > | EOF, set re-consume & change State AfterAttributeName", func(t *testing.T) {
			// Arrange
			tokenizer := init()

			// Act & Assert
			assert.Nil(t, tokenizer.tokenizeAttributeName('/'))
			assert.Equal(t, AfterAttributeName, tokenizer.State)
			assert.True(t, tokenizer.ReConsume)
		})

		t.Run("if =, change State BeforeAttributeValue", func(t *testing.T) {
			// Arrange
			tokenizer := init()

			// Act & Assert
			assert.Nil(t, tokenizer.tokenizeAttributeName('='))
			assert.Equal(t, BeforeAttributeValue, tokenizer.State)
		})

		t.Run("if a-zA-Z, append lower to the attribute name", func(t *testing.T) {
			// Arrange
			tokenizer := init()
			tokenizer.LatestToken = &HtmlToken{
				StartTag: &StartTag{
					Attributes: []Attribute{
						{name: "", value: ""},
					},
				},
			}

			// Act & Assert
			assert.Nil(t, tokenizer.tokenizeAttributeName('A'))
			assert.Equal(t, "a", tokenizer.LatestToken.StartTag.Attributes[0].name)
		})

		t.Run("if others input, append to the attribute name", func(t *testing.T) {
			// Arrange
			tokenizer := init()
			tokenizer.LatestToken = &HtmlToken{
				StartTag: &StartTag{
					Attributes: []Attribute{
						{name: "", value: ""},
					},
				},
			}

			// Act
			tokenizer.tokenizeAttributeName('@')

			// Assert
			assert.Equal(t, "@", tokenizer.LatestToken.StartTag.Attributes[0].name)
		})
	})

	t.Run("panic case", func(t *testing.T) {
		tests := []struct {
			name  string
			state State
			r     rune
		}{
			{
				name:  "the state is not AttributeName",
				state: Data,
				r:     'a',
			},
		}

		for _, it := range tests {
			t.Run(it.name, func(t *testing.T) {
				var err error
				defer func() {
					if r := recover(); r != nil {
						err = errors.New("panic")
					}
				}()

				// Act
				tokenizer := NewHtmlTokenizer("").(*HtmlTokenizer)
				tokenizer.State = it.state
				tokenizer.tokenizeAttributeName(it.r)

				// Assert
				assert.Error(t, err)
			})
		}
	})
}

func TestHtmlTokenizer_tokenizeAfterAttributeName(t *testing.T) {
	t.Run("normal case", func(t *testing.T) {
		init := func() *HtmlTokenizer {
			tokenizer := NewHtmlTokenizer("").(*HtmlTokenizer)
			tokenizer.State = AfterAttributeName
			return tokenizer
		}

		t.Run("if ' ', do nothing", func(t *testing.T) {
			// Arrange
			tokenizer := init()

			// Act & Assert
			assert.Nil(t, tokenizer.tokenizeAfterAttributeName(' '))
			assert.Equal(t, AfterAttributeName, tokenizer.State)
		})

		t.Run("if /, change State SelfClosingStartTag", func(t *testing.T) {
			// Arrange
			tokenizer := init()

			// Act & Assert
			assert.Nil(t, tokenizer.tokenizeAfterAttributeName('/'))
			assert.Equal(t, SelfClosingStartTag, tokenizer.State)
		})

		t.Run("if >, change State Data & return latestToken", func(t *testing.T) {
			// Arrange
			expected := newRuneToken('a')
			tokenizer := init()
			tokenizer.LatestToken = expected

			// Act
			result := tokenizer.tokenizeAfterAttributeName('>')

			// Assert
			assert.Equal(t, Data, tokenizer.State)
			assert.Equal(t, expected, result)
		})

		t.Run("if =, change State BeforeAttributeValue", func(t *testing.T) {
			// Arrange
			tokenizer := init()

			// Act & Assert
			assert.Nil(t, tokenizer.tokenizeAfterAttributeName('='))
			assert.Equal(t, BeforeAttributeValue, tokenizer.State)
		})

		t.Run("if EOF, return EOF token", func(t *testing.T) {
			// Arrange
			tokenizer := init()
			tokenizer.Pos = 1

			// Act
			result := tokenizer.tokenizeAfterAttributeName(rune(0))

			// Assert
			assert.Equal(t, newEOFToken(), result)
		})

		t.Run("if a-zA-Z, set re-consume & change State Data", func(t *testing.T) {
			// Arrange
			tokenizer := init()
			tokenizer.Input = []rune("a")
			tokenizer.LatestToken = &HtmlToken{
				StartTag: &StartTag{
					Attributes: []Attribute{},
				},
			}

			// Act & Assert
			assert.Nil(t, tokenizer.tokenizeAfterAttributeName('a'))
			assert.Equal(t, Data, tokenizer.State)
			assert.True(t, tokenizer.ReConsume)
			assert.Equal(t, len(tokenizer.LatestToken.StartTag.Attributes), 1)
		})
	})

	t.Run("panic case", func(t *testing.T) {
		tests := []struct {
			name  string
			state State
			r     rune
		}{
			{
				name:  "the state is not AttributeName",
				state: Data,
				r:     'a',
			},
		}

		for _, it := range tests {
			t.Run(it.name, func(t *testing.T) {
				var err error
				defer func() {
					if r := recover(); r != nil {
						err = errors.New("panic")
					}
				}()

				// Act
				tokenizer := NewHtmlTokenizer("").(*HtmlTokenizer)
				tokenizer.State = it.state
				tokenizer.tokenizeAfterAttributeName(it.r)

				// Assert
				assert.Error(t, err)
			})
		}
	})
}

func TestHtmlTokenizer_tokenizeBeforeAttributeValue(t *testing.T) {
	t.Run("normal case", func(t *testing.T) {
		init := func() *HtmlTokenizer {
			tokenizer := NewHtmlTokenizer("").(*HtmlTokenizer)
			tokenizer.State = BeforeAttributeValue
			return tokenizer
		}

		t.Run("if ' ', do nothing", func(t *testing.T) {
			// Arrange
			tokenizer := init()

			// Act & Assert
			assert.Nil(t, tokenizer.tokenizeBeforeAttributeValue(' '))
			assert.Equal(t, BeforeAttributeValue, tokenizer.State)
		})

		t.Run("if \", change State AttributeValueDoubleQuoted", func(t *testing.T) {
			// Arrange
			tokenizer := init()

			// Act & Assert
			assert.Nil(t, tokenizer.tokenizeBeforeAttributeValue('"'))
			assert.Equal(t, AttributeValueDoubleQuoted, tokenizer.State)
		})

		t.Run("if ', change State AttributeValueSingleQuoted", func(t *testing.T) {
			// Arrange
			tokenizer := init()

			// Act & Assert
			assert.Nil(t, tokenizer.tokenizeBeforeAttributeValue('\''))
			assert.Equal(t, AttributeValueSingleQuoted, tokenizer.State)
		})

		t.Run("if others input, set re-consume & change State AttributeValueUnquoted", func(t *testing.T) {
			// Arrange
			tokenizer := init()

			// Act & Assert
			assert.Nil(t, tokenizer.tokenizeBeforeAttributeValue('a'))
			assert.Equal(t, AttributeValueUnquoted, tokenizer.State)
			assert.True(t, tokenizer.ReConsume)
		})
	})

	t.Run("panic case", func(t *testing.T) {
		tests := []struct {
			name  string
			state State
			r     rune
		}{
			{
				name:  "the state is not BeforeAttributeValue",
				state: Data,
				r:     'a',
			},
		}

		for _, it := range tests {
			t.Run(it.name, func(t *testing.T) {
				var err error
				defer func() {
					if r := recover(); r != nil {
						err = errors.New("panic")
					}
				}()

				// Act
				tokenizer := NewHtmlTokenizer("").(*HtmlTokenizer)
				tokenizer.State = it.state
				tokenizer.tokenizeBeforeAttributeValue(it.r)

				// Assert
				assert.Error(t, err)
			})
		}
	})
}

func TestHtmlTokenizer_tokenizeAttributeValueDoubleQuoted(t *testing.T) {
	t.Run("normal case", func(t *testing.T) {
		init := func() *HtmlTokenizer {
			tokenizer := NewHtmlTokenizer("").(*HtmlTokenizer)
			tokenizer.State = AttributeValueDoubleQuoted
			tokenizer.LatestToken = &HtmlToken{
				StartTag: &StartTag{
					Attributes: []Attribute{{}},
				},
			}
			return tokenizer
		}

		t.Run("if \", change state to AfterAttributeValueQuoted", func(t *testing.T) {
			// Arrange
			tokenizer := init()

			// Act
			result := tokenizer.tokenizeAttributeValueDoubleQuoted('"')

			// Assert
			assert.Nil(t, result)
			assert.Equal(t, AfterAttributeValueQuoted, tokenizer.State)
		})

		t.Run("if EOF, return EOF token", func(t *testing.T) {
			// Arrange
			tokenizer := init()
			tokenizer.Pos = 1

			// Act
			result := tokenizer.tokenizeAttributeValueDoubleQuoted(rune(0))

			// Assert
			assert.Equal(t, newEOFToken(), result)
		})

		t.Run("otherwise, append to attribute value", func(t *testing.T) {
			// Arrange
			tokenizer := init()

			// Act
			result := tokenizer.tokenizeAttributeValueDoubleQuoted('a')

			// Assert
			assert.Nil(t, result)
			assert.Equal(t, "a", tokenizer.LatestToken.StartTag.Attributes[0].value)
		})
	})

	t.Run("panic case", func(t *testing.T) {
		t.Run("the state is not AttributeValueDoubleQuoted", func(t *testing.T) {
			var err error
			defer func() {
				if r := recover(); r != nil {
					err = errors.New("panic")
				}
			}()

			// Arrange
			tokenizer := NewHtmlTokenizer("").(*HtmlTokenizer)
			tokenizer.State = Data

			// Act
			tokenizer.tokenizeAttributeValueDoubleQuoted('a')

			// Assert
			assert.Error(t, err)
		})
	})
}

func TestHtmlTokenizer_tokenizeAttributeValueSingleQuoted(t *testing.T) {
	t.Run("normal case", func(t *testing.T) {
		init := func() *HtmlTokenizer {
			tokenizer := NewHtmlTokenizer("").(*HtmlTokenizer)
			tokenizer.State = AttributeValueSingleQuoted
			return tokenizer
		}

		t.Run("if ', change State AfterAttributeValueQuoted", func(t *testing.T) {
			// Arrange
			tokenizer := init()

			// Act & Assert
			assert.Nil(t, tokenizer.tokenizeAttributeValueSingleQuoted('\''))
			assert.Equal(t, AfterAttributeValueQuoted, tokenizer.State)
		})

		t.Run("if EOF, return newEOFToken", func(t *testing.T) {
			// Arrange
			tokenizer := init()
			tokenizer.Pos = 1

			// Act & Assert
			assert.Equal(t, newEOFToken(), tokenizer.tokenizeAttributeValueSingleQuoted(rune(0)))
		})

		t.Run("if other input, add attribute value", func(t *testing.T) {
			// Arrange
			tokenizer := init()
			tokenizer.LatestToken = &HtmlToken{
				StartTag: &StartTag{
					Attributes: []Attribute{{name: "", value: ""}},
				},
			}

			// Act & Assert
			assert.Nil(t, tokenizer.tokenizeAttributeValueSingleQuoted('a'))
			assert.Equal(t, "a", tokenizer.LatestToken.StartTag.Attributes[0].value)
		})
	})

	t.Run("panic case", func(t *testing.T) {
		t.Run("the state is not AttributeValueSingleQuoted", func(t *testing.T) {
			var err error
			defer func() {
				if r := recover(); r != nil {
					err = errors.New("panic")
				}
			}()

			// Arrange
			tokenizer := NewHtmlTokenizer("").(*HtmlTokenizer)
			tokenizer.State = Data

			// Act
			tokenizer.tokenizeAttributeValueSingleQuoted('a')

			// Assert
			assert.Error(t, err)
		})
	})
}

func TestHtmlTokenizer_tokenizeAttributeValueUnQuoted(t *testing.T) {
	t.Run("normal case", func(t *testing.T) {
		init := func() *HtmlTokenizer {
			tokenizer := NewHtmlTokenizer("").(*HtmlTokenizer)
			tokenizer.State = AttributeValueUnquoted
			return tokenizer
		}

		t.Run("if ' ', change State BeforeAttributeName", func(t *testing.T) {
			// Arrange
			tokenizer := init()

			// Act & Assert
			assert.Nil(t, tokenizer.tokenizeAttributeValueUnquoted(' '))
			assert.Equal(t, BeforeAttributeName, tokenizer.State)
		})

		t.Run("if >, change State Data & return latestToken", func(t *testing.T) {
			// Arrange
			expected := newRuneToken('a')
			tokenizer := init()
			tokenizer.LatestToken = expected

			// Act & Assert
			assert.Equal(t, tokenizer.tokenizeAttributeValueUnquoted('>'), expected)
			assert.Equal(t, Data, tokenizer.State)
		})

		t.Run("if other input, add attribute value", func(t *testing.T) {
			// Arrange
			tokenizer := init()
			tokenizer.LatestToken = &HtmlToken{
				StartTag: &StartTag{
					Attributes: []Attribute{{name: "", value: ""}},
				},
			}

			// Act & Assert
			assert.Nil(t, tokenizer.tokenizeAttributeValueUnquoted('a'))
			assert.Equal(t, "a", tokenizer.LatestToken.StartTag.Attributes[0].value)
		})
	})

	t.Run("panic case", func(t *testing.T) {
		t.Run("the state is not AttributeValueUnQuoted", func(t *testing.T) {
			var err error
			defer func() {
				if r := recover(); r != nil {
					err = errors.New("panic")
				}
			}()

			// Arrange
			tokenizer := NewHtmlTokenizer("").(*HtmlTokenizer)
			tokenizer.State = Data

			// Act
			tokenizer.tokenizeAttributeValueUnquoted('a')

			// Assert
			assert.Error(t, err)
		})
	})
}

func TestHtmlTokenizer_tokenizeAfterAttributeValueQuoted(t *testing.T) {
	t.Run("normal case", func(t *testing.T) {
		init := func() *HtmlTokenizer {
			tokenizer := NewHtmlTokenizer("").(*HtmlTokenizer)
			tokenizer.State = AfterAttributeValueQuoted
			return tokenizer
		}

		t.Run("if /, change State SelfClosingStartTag", func(t *testing.T) {
			// Arrange
			tokenizer := init()

			// Act & Assert
			assert.Nil(t, tokenizer.tokenizeAfterAttributeValueQuoted('/'))
			assert.Equal(t, SelfClosingStartTag, tokenizer.State)
		})

		t.Run("if >, change State Data & return latestToken", func(t *testing.T) {
			// Arrange
			expected := newRuneToken('a')
			tokenizer := init()
			tokenizer.LatestToken = expected

			// Act & Assert
			assert.Equal(t, expected, tokenizer.tokenizeAfterAttributeValueQuoted('>'))
			assert.Equal(t, Data, tokenizer.State)
		})

		t.Run("if EOF, return newEOFToken", func(t *testing.T) {
			// Arrange
			tokenizer := init()
			tokenizer.Pos = 1

			// Act & Assert
			assert.Equal(t, newEOFToken(), tokenizer.tokenizeAfterAttributeValueQuoted(rune(0)))
		})

		t.Run("if other input, set re-consume & change State BeforeAttributeName", func(t *testing.T) {
			// Arrange
			tokenizer := init()

			// Act & Assert
			assert.Nil(t, tokenizer.tokenizeAfterAttributeValueQuoted('a'))
			assert.Equal(t, tokenizer.State, BeforeAttributeName)
			assert.True(t, tokenizer.ReConsume)
		})
	})

	t.Run("panic case", func(t *testing.T) {
		t.Run("the state is not AfterAttributeValueQuoted", func(t *testing.T) {
			var err error
			defer func() {
				if r := recover(); r != nil {
					err = errors.New("panic")
				}
			}()

			// Arrange
			tokenizer := NewHtmlTokenizer("").(*HtmlTokenizer)
			tokenizer.State = Data

			// Act
			tokenizer.tokenizeAfterAttributeValueQuoted('a')

			// Assert
			assert.Error(t, err)
		})
	})
}

func TestHtmlTokenizer_tokenizeSelfClosingStartTag(t *testing.T) {
	t.Run("normal case", func(t *testing.T) {
		init := func() *HtmlTokenizer {
			tokenizer := NewHtmlTokenizer("").(*HtmlTokenizer)
			tokenizer.State = SelfClosingStartTag
			return tokenizer
		}

		t.Run("if >, set SelfClosingStartTag & change State Data & return latestToken", func(t *testing.T) {
			// Arrange
			tokenizer := init()
			tokenizer.LatestToken = &HtmlToken{
				StartTag: &StartTag{
					IsSelfClosing: false,
				},
			}

			// Act
			token := tokenizer.tokenizeSelfClosingStartTag('>')

			// Assert
			assert.Equal(t, Data, tokenizer.State)
			assert.True(t, token.StartTag.IsSelfClosing)
		})

		t.Run("if EOF, return EOF token", func(t *testing.T) {
			// Arrange
			tokenizer := init()
			tokenizer.Pos = 1

			// Act & Assert
			assert.Equal(t, newEOFToken(), tokenizer.tokenizeSelfClosingStartTag(rune(0)))
		})

		t.Run("if other input, do nothing", func(t *testing.T) {
			// Arrange
			tokenizer := init()

			// Act & Assert
			assert.Nil(t, tokenizer.tokenizeSelfClosingStartTag('a'))
			assert.Equal(t, SelfClosingStartTag, tokenizer.State)
		})
	})

	t.Run("panic case", func(t *testing.T) {
		t.Run("the state is not SelfClosingStartTag", func(t *testing.T) {
			var err error
			defer func() {
				if r := recover(); r != nil {
					err = errors.New("panic")
				}
			}()

			// Arrange
			tokenizer := NewHtmlTokenizer("").(*HtmlTokenizer)
			tokenizer.State = Data

			// Act
			tokenizer.tokenizeSelfClosingStartTag('a')

			// Assert
			assert.Error(t, err)
		})
	})
}

func TestHtmlTokenizer_tokenizeScriptData(t *testing.T) {
	t.Run("normal case", func(t *testing.T) {
		init := func() *HtmlTokenizer {
			tokenizer := NewHtmlTokenizer("").(*HtmlTokenizer)
			tokenizer.State = ScriptData
			return tokenizer
		}

		t.Run("if <, change State ScriptDataLessThanSign", func(t *testing.T) {
			// Arrange
			tokenizer := init()

			// Act & Assert
			assert.Nil(t, tokenizer.tokenizeScriptData('<'))
			assert.Equal(t, ScriptDataLessThanSign, tokenizer.State)
		})

		t.Run("if EOF, return newEOFToken", func(t *testing.T) {
			// Arrange
			tokenizer := init()
			tokenizer.Pos = 1

			// Act & Assert
			assert.Equal(t, newEOFToken(), tokenizer.tokenizeScriptData(rune(0)))
		})

		t.Run("if other input, return runeToken", func(t *testing.T) {
			// Arrange
			tokenizer := init()

			// Act & Assert
			assert.Equal(t, newRuneToken('a'), tokenizer.tokenizeScriptData('a'))
		})
	})

	t.Run("panic case", func(t *testing.T) {
		t.Run("the state is not ScriptData", func(t *testing.T) {
			var err error
			defer func() {
				if r := recover(); r != nil {
					err = errors.New("panic")
				}
			}()

			// Arrange
			tokenizer := NewHtmlTokenizer("").(*HtmlTokenizer)
			tokenizer.State = Data

			// Act
			tokenizer.tokenizeScriptData('a')

			// Assert
			assert.Error(t, err)
		})
	})
}

func TestHtmlTokenizer_tokenizeScriptDataLessThanSign(t *testing.T) {
	t.Run("normal case", func(t *testing.T) {
		init := func() *HtmlTokenizer {
			tokenizer := NewHtmlTokenizer("").(*HtmlTokenizer)
			tokenizer.State = ScriptDataLessThanSign
			return tokenizer
		}

		t.Run("if /, reset Buf & change State ScriptDataEndTagOpen", func(t *testing.T) {
			// Arrange
			tokenizer := init()
			tokenizer.Buf = "temp"

			// Act & Assert
			assert.Nil(t, tokenizer.tokenizeScriptDataLessThanSign('/'))
			assert.Equal(t, ScriptDataEndTagOpen, tokenizer.State)
			assert.Equal(t, "", tokenizer.Buf)
		})

		t.Run("if other input, set re-consume & change State ScriptData & return runeToken of <", func(t *testing.T) {
			// Arrange
			tokenizer := init()

			// Act & Assert
			assert.Equal(t, tokenizer.tokenizeScriptDataLessThanSign('@'), newRuneToken('<'))
			assert.Equal(t, tokenizer.State, ScriptData)
			assert.True(t, tokenizer.ReConsume)
		})
	})

	t.Run("panic case", func(t *testing.T) {
		t.Run("the state is not ScriptDataLessThanSign", func(t *testing.T) {
			var err error
			defer func() {
				if r := recover(); r != nil {
					err = errors.New("panic")
				}
			}()

			// Arrange
			tokenizer := NewHtmlTokenizer("").(*HtmlTokenizer)
			tokenizer.State = Data

			// Act
			tokenizer.tokenizeScriptDataLessThanSign('a')

			// Assert
			assert.Error(t, err)
		})
	})
}

func TestHtmlTokenizer_tokenizeScriptDataEndTagOpen(t *testing.T) {
	t.Run("normal case", func(t *testing.T) {
		init := func() *HtmlTokenizer {
			tokenizer := NewHtmlTokenizer("").(*HtmlTokenizer)
			tokenizer.State = ScriptDataEndTagOpen
			return tokenizer
		}

		t.Run("if a-zA-Z, set re-consume & change State ScriptDataEndTagName & create EndTag", func(t *testing.T) {
			// Arrange
			tokenizer := init()

			// Act & Assert
			assert.Nil(t, tokenizer.tokenizeScriptDataEndTagOpen('a'))
			assert.Equal(t, tokenizer.State, ScriptDataEndTagName)
			assert.True(t, tokenizer.ReConsume)
			assert.NotNil(t, tokenizer.LatestToken)
		})

		t.Run("if others input, set re-consume & change State ScriptData & return runeToken of <", func(t *testing.T) {
			// Arrange
			tokenizer := init()

			// Act & Assert
			assert.Equal(t, tokenizer.tokenizeScriptDataEndTagOpen('@'), newRuneToken('<'))
			assert.Equal(t, tokenizer.State, ScriptData)
			assert.True(t, tokenizer.ReConsume)
		})
	})

	t.Run("panic case", func(t *testing.T) {
		t.Run("the state is not ScriptDataEndTagOpen", func(t *testing.T) {
			var err error
			defer func() {
				if r := recover(); r != nil {
					err = errors.New("panic")
				}
			}()

			// Arrange
			tokenizer := NewHtmlTokenizer("").(*HtmlTokenizer)
			tokenizer.State = Data

			// Act
			tokenizer.tokenizeScriptDataEndTagOpen('a')

			// Assert
			assert.Error(t, err)
		})
	})
}

func TestHtmlTokenizer_tokenizeScriptDataEndTagName(t *testing.T) {
	t.Run("normal case", func(t *testing.T) {
		init := func() *HtmlTokenizer {
			tokenizer := NewHtmlTokenizer("").(*HtmlTokenizer)
			tokenizer.State = ScriptDataEndTagName
			return tokenizer
		}

		t.Run("if >, change Sate Data & return token", func(t *testing.T) {
			// Arrange
			tokenizer := init()
			expected := newRuneToken('a')
			tokenizer.LatestToken = expected

			// Act & Assert
			assert.Equal(t, expected, tokenizer.tokenizeScriptDataEndTagName('>'))
			assert.Equal(t, Data, tokenizer.State)
		})

		t.Run("if a-zA-z, append lower to tag name", func(t *testing.T) {
			// Arrange
			tokenizer := init()
			tokenizer.LatestToken = &HtmlToken{
				EndTag: &EndTag{
					Tag: "",
				},
			}

			// Act & Assert
			assert.Nil(t, tokenizer.tokenizeScriptDataEndTagName('a'))
			assert.Equal(t, "a", tokenizer.LatestToken.EndTag.Tag)
		})

		t.Run("if others input, change state TemporaryBuffer & add </ to Buf, append to tag name", func(t *testing.T) {
			// Arrange
			tokenizer := init()
			tokenizer.LatestToken = &HtmlToken{
				EndTag: &EndTag{
					Tag: "",
				},
			}

			// Act & Assert
			assert.Nil(t, tokenizer.tokenizeScriptDataEndTagName('@'))
			assert.Equal(t, TemporaryBuffer, tokenizer.State)
			assert.Equal(t, "</@", tokenizer.Buf)
		})
	})

	t.Run("panic case", func(t *testing.T) {
		t.Run("the state is not ScriptDataEndTagName", func(t *testing.T) {
			var err error
			defer func() {
				if r := recover(); r != nil {
					err = errors.New("panic")
				}
			}()

			// Arrange
			tokenizer := NewHtmlTokenizer("").(*HtmlTokenizer)
			tokenizer.State = Data

			// Act
			tokenizer.tokenizeScriptDataEndTagName('a')

			// Assert
			assert.Error(t, err)
		})
	})
}

func TestHtmlTokenizer_tokenizeTemporaryBuffer(t *testing.T) {
	t.Run("normal case, set re-consume", func(t *testing.T) {
		init := func() *HtmlTokenizer {
			tokenizer := NewHtmlTokenizer("").(*HtmlTokenizer)
			tokenizer.State = TemporaryBuffer
			return tokenizer
		}

		t.Run("if Buf is empty, change State ScriptData", func(t *testing.T) {
			// Arrange
			tokenizer := init()
			tokenizer.Buf = ""

			// Act & Assert
			assert.Nil(t, tokenizer.tokenizeTemporaryBuffer('a'))
			assert.Equal(t, ScriptData, tokenizer.State)
			assert.True(t, tokenizer.ReConsume)
		})

		t.Run("if Buf is valid, return runeToken", func(t *testing.T) {
			// Arrange
			tokenizer := init()
			tokenizer.Buf = "hello"

			// Act & Assert
			assert.Equal(t, tokenizer.tokenizeTemporaryBuffer('a'), newRuneToken('a'))
			assert.True(t, tokenizer.ReConsume)
			assert.Equal(t, tokenizer.State, TemporaryBuffer)
			assert.Equal(t, "ello", tokenizer.Buf)
		})
	})

	t.Run("panic case", func(t *testing.T) {
		t.Run("the state is not TemporaryBuffer", func(t *testing.T) {
			var err error
			defer func() {
				if r := recover(); r != nil {
					err = errors.New("panic")
				}
			}()

			// Arrange
			tokenizer := NewHtmlTokenizer("").(*HtmlTokenizer)
			tokenizer.State = Data

			// Act
			tokenizer.tokenizeTemporaryBuffer('a')

			// Assert
			assert.Error(t, err)
		})
	})
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

func TestHtmlTokenizer_appendAttribute(t *testing.T) {
	t.Run("normal case, append rune to target attribute", func(t *testing.T) {
		tests := []struct {
			name        string
			r           rune
			isName      bool
			latestToken *HtmlToken
			expected    []Attribute
		}{
			{
				name:   "append to start tag",
				r:      'a',
				isName: true,
				latestToken: &HtmlToken{
					StartTag: &StartTag{
						Attributes: []Attribute{
							{},
						},
					},
				},
				expected: []Attribute{
					{
						name:  "a",
						value: "",
					},
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
				tokenizer.appendAttribute(it.r, it.isName)

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
				name: "StartTag w/ 0 length attribute list",
				latestToken: &HtmlToken{
					StartTag: &StartTag{
						Attributes: []Attribute{},
					},
				},
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
				tokenizer.appendAttribute('a', true)

				// Assert
				assert.Error(t, err)
			})
		}
	})
}

func TestHtmlTokenizer_setSelfClosingFlag(t *testing.T) {
	t.Run("normal case, set IsSelfClosing flag to true", func(t *testing.T) {
		tests := []struct {
			name        string
			latestToken *HtmlToken
			expected    bool
		}{
			{
				name: "set IsSelfClosing flag to true",
				latestToken: &HtmlToken{
					StartTag: &StartTag{
						IsSelfClosing: false,
					},
				},
				expected: true,
			},
		}

		// Arrange
		tokenizer := NewHtmlTokenizer("").(*HtmlTokenizer)

		for _, it := range tests {
			t.Run(it.name, func(t *testing.T) {
				// Arrange
				tokenizer.LatestToken = it.latestToken

				// Act
				tokenizer.setSelfClosingFlag()

				// Assert
				assert.Equal(t, it.expected, tokenizer.LatestToken.StartTag.IsSelfClosing)
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
				tokenizer.setSelfClosingFlag()

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
