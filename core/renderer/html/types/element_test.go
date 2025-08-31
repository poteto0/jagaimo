package types

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConvertToElementKind(t *testing.T) {
	t.Run("normal case", func(t *testing.T) {
		tests := []struct {
			elementName string
			expected    ElementKind
		}{
			{
				elementName: "html",
				expected:    Html,
			},
			{
				elementName: "head",
				expected:    Head,
			},
			{
				elementName: "style",
				expected:    Style,
			},
			{
				elementName: "script",
				expected:    Script,
			},
			{
				elementName: "body",
				expected:    Body,
			},
		}

		for _, it := range tests {
			t.Run(fmt.Sprintf("can convert from %s", it.elementName), func(t *testing.T) {
				// Act
				ele, _ := ConvertToElementKind(it.elementName)

				// Assert
				assert.Equal(t, it.expected, ele)
			})
		}
	})

	t.Run("error case", func(t *testing.T) {
		// Act
		_, err := ConvertToElementKind("invalid")

		// Assert
		assert.Error(t, err)
	})
}

func TestNewElement(t *testing.T) {
	// Act & Assert
	assert.IsType(t, &Element{}, NewElement("html", []Attribute{}))
}

func TestElement_Kind(t *testing.T) {
	// Arrange
	element := NewElement("html", []Attribute{}).(*Element)

	// Act & Assert
	assert.Equal(t, Html, element.Kind())
}
