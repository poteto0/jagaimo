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
