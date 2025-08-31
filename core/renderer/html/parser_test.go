package html

import (
	"errors"
	"testing"

	"github.com/poteto0/jagaimo/core/renderer/dom"
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

		t.Run("if startTag & tag is html, insert node & set BeforeHead mode, return next", func(t *testing.T) {
			// Arrange
			parser := init()
			parser.t = NewHtmlTokenizer(
				"<html><head></head><body></body></html>",
			)

			// Act
			next, _ := parser.parseBeforeHtml(parser.t.Next())

			// Assert
			assert.Equal(t, &HtmlToken{
				StartTag: &StartTag{
					Tag:        "head",
					Attributes: []types.Attribute{},
				},
			}, next)
			assert.Equal(t, BeforeHead, parser.mode)
		})

		t.Run("if EOFToken, return nil & finish parsing", func(t *testing.T) {
			// Arrange
			parser := init()

			// Act
			next, isFinished := parser.parseBeforeHtml(newEOFToken())

			// Assert
			assert.Nil(t, next)
			assert.True(t, isFinished)
		})

		t.Run("if other token, insert html node & change mode BeforeHead & return nil", func(t *testing.T) {
			// Arrange
			parser := init()

			// Act
			next, _ := parser.parseBeforeHtml(newRuneToken('/'))

			// Assert
			assert.Nil(t, next)
		})
	})

	t.Run("panic case", func(t *testing.T) {
		t.Run("the mode is not BeforeHtml", func(t *testing.T) {
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
			assert.Equal(t, Initial, parser.mode)
		})
	})
}

func TestHtmlParser_parseBeforeHead(t *testing.T) {
	t.Run("normal case", func(t *testing.T) {
		init := func() *HtmlParser {
			parser := NewHtmlParser(NewHtmlTokenizer(
				" <head></head><body></body>",
			)).(*HtmlParser)
			parser.mode = BeforeHead
			return parser
		}

		t.Run("if rune token of ' ' or \n , return next token", func(t *testing.T) {
			// Arrange
			parser := init()
			token := parser.t.Next()

			// Act
			next, _ := parser.parseBeforeHead(token)
			assert.Equal(t, &HtmlToken{
				StartTag: &StartTag{
					Tag:           "head",
					IsSelfClosing: false,
					Attributes:    []types.Attribute{},
				},
			}, next)
		})

		t.Run("if startTag & tag is head, insert head node & set InHead mode, return next", func(t *testing.T) {
			// Arrange
			parser := init()
			parser.t = NewHtmlTokenizer(
				"<head></head><body></body>",
			)

			// Act
			next, _ := parser.parseBeforeHead(parser.t.Next())

			// Assert
			assert.Equal(t, &HtmlToken{
				EndTag: &EndTag{
					Tag: "head",
				},
			}, next)
			assert.Equal(t, InHead, parser.mode)
		})

		t.Run("if EOFToken, return nil & finish parsing", func(t *testing.T) {
			// Arrange
			parser := init()

			// Act
			next, isFinished := parser.parseBeforeHead(newEOFToken())

			// Assert
			assert.Nil(t, next)
			assert.True(t, isFinished)
		})

		t.Run("if other token, insert head node & change mode InHead & return nil", func(t *testing.T) {
			// Arrange
			parser := init()

			// Act
			next, _ := parser.parseBeforeHead(newRuneToken('/'))

			// Assert
			assert.Nil(t, next)
			assert.Equal(t, InHead, parser.mode)
		})
	})

	t.Run("panic case", func(t *testing.T) {
		t.Run("the mode is not BeforeHead", func(t *testing.T) {
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
			parser.parseBeforeHead(nil)

			// Assert
			assert.Error(t, err)
		})
	})
}

func TestHtmlParser_parseInHead(t *testing.T) {
	t.Run("normal case", func(t *testing.T) {
		init := func() *HtmlParser {
			parser := NewHtmlParser(NewHtmlTokenizer(
				" <script></script>",
			)).(*HtmlParser)
			parser.mode = InHead
			return parser
		}

		t.Run("if token has ' ' or '\n', return next (if current is text append rune)", func(t *testing.T) {
			// Arrange
			parser := init()
			node := createRune('a')
			parser.stackOfOpenElements = []*dom.Node{
				node,
			}
			token := parser.t.Next()

			// Act
			next, _ := parser.parseInHead(token)

			// Assert
			assert.Equal(t, &HtmlToken{
				StartTag: &StartTag{
					Tag:           "script",
					IsSelfClosing: false,
					Attributes:    []types.Attribute{},
				},
			}, next)
			assert.Equal(t, parser.currentNode().Kind.Text, "a ")
		})

		t.Run("if token is start script, return token & insert Element & set Text mode", func(t *testing.T) {
			// Arrange
			parser := init()
			parser.t = NewHtmlTokenizer(
				"<script></script>",
			)

			// Act
			next, _ := parser.parseInHead(parser.t.Next())

			// Assert
			assert.Equal(t, &HtmlToken{
				EndTag: &EndTag{
					Tag: "script",
				},
			}, next)
			assert.Equal(t, Text, parser.mode)
			assert.Equal(t, InHead, parser.originalInsertionMode)
		})

		t.Run("if token is start style, return token & insert Element & set Text mode", func(t *testing.T) {
			// Arrange
			parser := init()
			parser.t = NewHtmlTokenizer(
				"<style></style>",
			)

			// Act
			next, _ := parser.parseInHead(parser.t.Next())

			// Assert
			assert.Equal(t, &HtmlToken{
				EndTag: &EndTag{
					Tag: "style",
				},
			}, next)
			assert.Equal(t, Text, parser.mode)
			assert.Equal(t, InHead, parser.originalInsertionMode)
		})

		t.Run("if token is start body, call popUntil & set AfterHead mode, return nil", func(t *testing.T) {
			// Arrange
			parser := init()
			parser.t = NewHtmlTokenizer(
				"<body></body>",
			)
			node1 := dom.NewNode(dom.NodeKind{
				Element: types.NewElement("html", []types.Attribute{}).(*types.Element),
			}).(*dom.Node)
			node2 := dom.NewNode(dom.NodeKind{}).(*dom.Node)
			node3 := dom.NewNode(dom.NodeKind{
				Element: types.NewElement("head", []types.Attribute{}).(*types.Element),
			}).(*dom.Node)
			node4 := dom.NewNode(dom.NodeKind{}).(*dom.Node)
			parser.stackOfOpenElements = []*dom.Node{
				node1,
				node2,
				node3,
				node4,
			}

			// Act
			next, _ := parser.parseInHead(parser.t.Next())

			// Assert
			assert.Nil(t, next)
			assert.Equal(t, AfterHead, parser.mode)
			assert.Equal(t, []*dom.Node{
				node1,
				node2,
			}, parser.stackOfOpenElements)
		})

		t.Run("if startTag not supported tag, call popUntil & set AfterHead mode, return nil", func(t *testing.T) {
			// Arrange
			parser := init()
			parser.t = NewHtmlTokenizer(
				"<meta></meta>",
			)
			node1 := dom.NewNode(dom.NodeKind{
				Element: types.NewElement("html", []types.Attribute{}).(*types.Element),
			}).(*dom.Node)
			node2 := dom.NewNode(dom.NodeKind{}).(*dom.Node)
			node3 := dom.NewNode(dom.NodeKind{
				Element: types.NewElement("head", []types.Attribute{}).(*types.Element),
			}).(*dom.Node)
			node4 := dom.NewNode(dom.NodeKind{}).(*dom.Node)
			parser.stackOfOpenElements = []*dom.Node{
				node1,
				node2,
				node3,
				node4,
			}

			// Act
			next, _ := parser.parseInHead(parser.t.Next())

			// Asset
			assert.Nil(t, next)
			assert.Equal(t, AfterHead, parser.mode)
			assert.Equal(t, []*dom.Node{
				node1,
				node2,
			}, parser.stackOfOpenElements)
		})

		t.Run("if token is end head tag, return next & set AfterHead & call popUntil", func(t *testing.T) {
			// Arrange
			parser := init()
			parser.t = NewHtmlTokenizer(
				"</head><body></body>",
			)
			node1 := dom.NewNode(dom.NodeKind{
				Element: types.NewElement("html", []types.Attribute{}).(*types.Element),
			}).(*dom.Node)
			node2 := dom.NewNode(dom.NodeKind{}).(*dom.Node)
			node3 := dom.NewNode(dom.NodeKind{
				Element: types.NewElement("head", []types.Attribute{}).(*types.Element),
			}).(*dom.Node)
			node4 := dom.NewNode(dom.NodeKind{}).(*dom.Node)
			parser.stackOfOpenElements = []*dom.Node{
				node1,
				node2,
				node3,
				node4,
			}

			// Act
			next, _ := parser.parseInHead(parser.t.Next())

			// Assert
			assert.Equal(t, next, &HtmlToken{
				StartTag: &StartTag{
					Tag:           "body",
					IsSelfClosing: false,
					Attributes:    []types.Attribute{},
				},
			})
			assert.Equal(t, AfterHead, parser.mode)
			assert.Equal(t, []*dom.Node{
				node1,
				node2,
			}, parser.stackOfOpenElements)
		})

		t.Run("if not supported tag, just skip", func(t *testing.T) {
			// Arrange
			parser := init()
			parser.t = NewHtmlTokenizer(
				"</meta></head>",
			)

			// Act
			next, _ := parser.parseInHead(parser.t.Next())

			// Asset
			assert.Equal(t, next, &HtmlToken{
				EndTag: &EndTag{
					Tag: "head",
				},
			})
		})

		t.Run("if EOF, return nil & finished", func(t *testing.T) {
			// Arrange
			parser := init()

			// Act
			next, isFinished := parser.parseInHead(newEOFToken())

			// Assert
			assert.Nil(t, next)
			assert.True(t, isFinished)
		})
	})

	t.Run("panic case", func(t *testing.T) {
		t.Run("the mode is not InHead", func(t *testing.T) {
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
			parser.parseInHead(nil)

			// Assert
			assert.Error(t, err)
		})
	})
}

func TestHtmlParser_insertElement(t *testing.T) {
	t.Run("set currentNode's last child & create node & append stackOfOpenElements", func(t *testing.T) {
		t.Run("if currentNode's first child is nil, set currentNode's FirstChild node", func(t *testing.T) {
			// Arrange
			tag1 := "html"
			attributes1 := []types.Attribute{}
			node1 := createElementNode(tag1, attributes1)

			tag2 := "html"
			attributes2 := []types.Attribute{}
			node2 := createElementNode(tag2, attributes2)

			parser := NewHtmlParser(nil).(*HtmlParser)
			parser.stackOfOpenElements = []*dom.Node{
				node1, // current node
			}

			// Act
			parser.insertElement(tag2, attributes2)

			// Assert
			assert.Equal(t, node2.ElementKind(), node1.LastChild.Value().ElementKind())
			createdNode := parser.stackOfOpenElements[len(parser.stackOfOpenElements)-1]
			assert.Equal(
				t,
				node2.ElementKind(),
				createdNode.ElementKind(),
			)
			assert.Equal(
				t,
				node2.ElementKind(),
				createdNode.Parent.Value().ElementKind(),
			)
			targetNode := parser.stackOfOpenElements[len(parser.stackOfOpenElements)-2]
			assert.Equal(
				t,
				node2.ElementKind(),
				targetNode.FirstChild.ElementKind(),
			)
		})

		t.Run("if currentNode's first child is not nil, set currentNode's FirstChild's last sibling node", func(t *testing.T) {
			// Arrange
			tag1 := "html"
			attributes1 := []types.Attribute{}
			node1 := createElementNode(tag1, attributes1)
			nodeChild := createElementNode(tag1, attributes1)
			lastSibling := createElementNode(tag1, attributes1)
			nodeChild.NextSibling = lastSibling
			node1.FirstChild = nodeChild

			tag2 := "html"
			attributes2 := []types.Attribute{}
			node2 := createElementNode(tag2, attributes2)

			parser := NewHtmlParser(nil).(*HtmlParser)
			parser.stackOfOpenElements = []*dom.Node{
				node1, // current node
			}

			// Act
			parser.insertElement(tag2, attributes2)

			// Assert
			assert.Equal(t, node2.ElementKind(), node1.LastChild.Value().ElementKind())
			createdNode := parser.stackOfOpenElements[len(parser.stackOfOpenElements)-1]
			assert.Equal(
				t,
				node2.ElementKind(),
				createdNode.ElementKind(),
			)
			assert.Equal(
				t,
				node2.ElementKind(),
				createdNode.Parent.Value().ElementKind(),
			)
			assert.Equal(
				t,
				node2.ElementKind(),
				createdNode.PrevSibling.Value().ElementKind(),
			)
			assert.Equal(
				t,
				node2.ElementKind(),
				lastSibling.NextSibling.ElementKind(),
			)
		})
	})
}

func TestHtmlParser_currentNode(t *testing.T) {
	// Arrange
	parser := NewHtmlParser(nil).(*HtmlParser)
	tests := []struct {
		name     string
		stack    []*dom.Node
		expected *dom.Node
	}{
		{
			name:     "if 0 length stack, return document",
			stack:    []*dom.Node{},
			expected: parser.window.Document(),
		},
		{
			name: "if last element is not nil, return last element",
			stack: []*dom.Node{
				{},
			},
			expected: &dom.Node{},
		},
		{
			name: "if last element is nil, return document",
			stack: []*dom.Node{
				nil,
			},
			expected: parser.window.Document(),
		},
	}

	for _, it := range tests {
		t.Run(it.name, func(t *testing.T) {
			parser.stackOfOpenElements = it.stack

			// Act & Assert
			assert.Equal(t, it.expected, parser.currentNode())
		})
	}
}

func TestHtmlParser_insertRune(t *testing.T) {
	t.Run("if text element, insert to text", func(t *testing.T) {
		// Arrange
		parser := NewHtmlParser(nil).(*HtmlParser)
		node := createRune('a')
		parser.stackOfOpenElements = []*dom.Node{
			node,
		}

		// Act
		parser.insertRune('b')

		// Assert
		assert.Equal(t, parser.currentNode().Kind.Text, "ab")
	})

	t.Run("if \n or ' ', & not text element, do nothing", func(t *testing.T) {
		// Arrange
		parser := NewHtmlParser(nil).(*HtmlParser)
		node := createElementNode("html", []types.Attribute{})
		parser.stackOfOpenElements = []*dom.Node{
			node,
		}

		for _, r := range []rune{' ', '\n'} {
			// Act
			parser.insertRune(r)

			// Assert
			assert.Equal(
				t,
				parser.currentNode().Kind.Element.Kind(),
				node.Kind.Element.Kind(),
			)
		}
	})

	t.Run("if not text element & other input, create runeNode & link", func(t *testing.T) {
		t.Run("if currentNode's FirstChild is nil, set currentNode's FirstChild node", func(t *testing.T) {
			// Arrange
			parser := NewHtmlParser(nil).(*HtmlParser)
			node := createElementNode("html", []types.Attribute{})
			parser.stackOfOpenElements = []*dom.Node{
				node,
			}

			// Act
			parser.insertRune('a')

			// Assert
			assert.True(t, node.LastChild.Value().Kind.IsText())
			assert.True(t, parser.currentNode().Kind.IsText())
			assert.Equal(t, node.FirstChild.ElementKind(), parser.currentNode().ElementKind())
		})

		t.Run("if currentNode's FirstChild is not nil, set FirstChild's nextSibling", func(t *testing.T) {
			// Arrange
			parser := NewHtmlParser(nil).(*HtmlParser)
			node := createElementNode("html", []types.Attribute{})
			node.FirstChild = createRune('a')
			parser.stackOfOpenElements = []*dom.Node{
				node,
			}

			// Act
			parser.insertRune('a')

			// Assert
			assert.True(t, node.LastChild.Value().Kind.IsText())
			assert.True(t, parser.currentNode().Kind.IsText())
			assert.True(t, node.FirstChild.NextSibling.Kind.IsText())
		})
	})
}

func TestHtmlParser_popUntil(t *testing.T) {
	t.Run("normal case", func(t *testing.T) {
		t.Run("pop until the target elements", func(t *testing.T) {
			// Arrange
			parser := NewHtmlParser(nil).(*HtmlParser)
			node1 := dom.NewNode(dom.NodeKind{
				Element: types.NewElement("html", []types.Attribute{}).(*types.Element),
			}).(*dom.Node)
			node2 := dom.NewNode(dom.NodeKind{}).(*dom.Node)
			node3 := dom.NewNode(dom.NodeKind{
				Element: types.NewElement("head", []types.Attribute{}).(*types.Element),
			}).(*dom.Node)
			node4 := dom.NewNode(dom.NodeKind{}).(*dom.Node)
			parser.stackOfOpenElements = []*dom.Node{
				node1,
				node2,
				node3,
				node4,
			}

			// Act
			parser.popUntil(types.Head)

			// Assert
			assert.Equal(t, []*dom.Node{
				node1,
				node2,
			}, parser.stackOfOpenElements)
		})
	})

	t.Run("panic case", func(t *testing.T) {
		t.Run("if stack doesn't have target element, panic", func(t *testing.T) {
			var err error
			defer func() {
				if r := recover(); r != nil {
					err = errors.New("panic")
				}
			}()

			// Arrange
			parser := NewHtmlParser(nil).(*HtmlParser)

			// Act
			parser.popUntil(types.Html)

			// Assert
			assert.Error(t, err)
		})
	})
}

func TestHtmlParser_hasKindInStack(t *testing.T) {
	tests := []struct {
		name     string
		stack    []*dom.Node
		kind     types.ElementKind
		expected bool
	}{
		{
			name:     "if 0 length stack, return false",
			stack:    []*dom.Node{},
			kind:     types.Html,
			expected: false,
		},
		{
			name: "if has target element",
			stack: []*dom.Node{
				dom.NewNode(dom.NodeKind{
					Element: types.NewElement("html", []types.Attribute{}).(*types.Element),
				}).(*dom.Node),
				dom.NewNode(dom.NodeKind{
					Element: types.NewElement("head", []types.Attribute{}).(*types.Element),
				}).(*dom.Node),
			},
			kind:     types.Head,
			expected: true,
		},
		{
			name: "if doesn't have target element",
			stack: []*dom.Node{
				dom.NewNode(dom.NodeKind{
					Element: types.NewElement("html", []types.Attribute{}).(*types.Element),
				}).(*dom.Node),
			},
			kind:     types.Head,
			expected: false,
		},
	}

	for _, it := range tests {
		t.Run(it.name, func(t *testing.T) {
			// Arrange
			parser := NewHtmlParser(nil).(*HtmlParser)
			parser.stackOfOpenElements = it.stack

			// Act & Assert
			assert.Equal(t, it.expected, parser.hasKindInStack(it.kind))
		})
	}
}

func TestHtmlParser_popCurrentNode(t *testing.T) {
	t.Run("normal case", func(t *testing.T) {
		t.Run("pop node from stack", func(t *testing.T) {
			// Arrange
			parser := NewHtmlParser(nil).(*HtmlParser)
			node1 := &dom.Node{}
			node2 := &dom.Node{}
			parser.stackOfOpenElements = []*dom.Node{
				node1,
				node2,
			}

			// Act
			result := parser.popCurrentNode()

			// Assert
			assert.Equal(t, node2, result)
			assert.Equal(t, []*dom.Node{
				node1,
			}, parser.stackOfOpenElements)
		})
	})

	t.Run("panic case", func(t *testing.T) {
		t.Run("if stack is empty, panic", func(t *testing.T) {
			var err error
			defer func() {
				if r := recover(); r != nil {
					err = errors.New("panic")
				}
			}()

			// Arrange
			parser := NewHtmlParser(nil).(*HtmlParser)

			// Act
			parser.popCurrentNode()

			// Assert
			assert.Error(t, err)
		})
	})
}
