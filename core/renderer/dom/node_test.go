package dom

import (
	"testing"
	"weak"

	htmlTypes "github.com/poteto0/jagaimo/core/renderer/html/types"
	"github.com/stretchr/testify/assert"
)

func TestNodeKind_IsDocument(t *testing.T) {
	// Arrange
	kind := NodeKind{
		Document: 1,
	}

	// Act & Assert
	assert.True(t, kind.IsDocument())
}

func TestNodeKind_IsElement(t *testing.T) {
	// Arrange
	kind := NodeKind{
		Element: &htmlTypes.Element{},
	}

	// Act & Assert
	assert.True(t, kind.IsElement())
}

func TestNodeKind_IsText(t *testing.T) {
	// Arrange
	kind := NodeKind{
		Text:    "text",
		HasText: true,
	}

	// Act & Assert
	assert.True(t, kind.IsText())
}

func TestNewNode(t *testing.T) {
	// Act & Assert
	assert.IsType(t, &Node{}, NewNode(NodeKind{}))
}

func TestNode_SetWindow(t *testing.T) {
	// Arrange
	window := weak.Make(NewWindow().(*Window))
	node := NewNode(NodeKind{}).(*Node)

	// Act
	node.SetWindow(window)

	// Assert
	assert.Equal(t, window, node.window)
}

func TestNode_GetElement(t *testing.T) {
	t.Run("if nodeKind is element, return element", func(t *testing.T) {
		// Arrange
		kind := NodeKind{
			Element: &htmlTypes.Element{},
		}
		node := NewNode(kind).(*Node)

		// Act & Assert
		assert.Equal(t, kind.Element, node.GetElement())
	})

	t.Run("otherwise, return nil", func(t *testing.T) {
		// Arrange
		kind := NodeKind{}
		node := NewNode(kind).(*Node)

		// Act & Assert
		assert.Nil(t, node.GetElement())
	})
}

func TestNode_ElementKind(t *testing.T) {
	t.Run("if nodeKind is element, return element kind", func(t *testing.T) {
		// Arrange
		kind := NodeKind{
			Element: htmlTypes.NewElement("html", []htmlTypes.Attribute{}).(*htmlTypes.Element),
		}
		node := NewNode(kind).(*Node)

		// Act & Assert
		assert.Equal(t, kind.Element.Kind(), node.ElementKind())
	})

	t.Run("otherwise, return NilElement", func(t *testing.T) {
		// Arrange
		kind := NodeKind{}
		node := NewNode(kind).(*Node)

		// Act & Assert
		assert.Equal(t, htmlTypes.NilElement, node.ElementKind())
	})
}
