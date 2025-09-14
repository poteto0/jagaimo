package utils

import (
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/poteto0/jagaimo/core/renderer/dom"
	htmlTypes "github.com/poteto0/jagaimo/core/renderer/html/types"
	"github.com/stretchr/testify/assert"
)

func TestConvertDomToString(t *testing.T) {
	// Arrange
	root := dom.NewNode(dom.NodeKind{
		Element: &htmlTypes.Element{},
	}).(*dom.Node)

	t.Run("call convertDomToStringInternal", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		// Mock
		isCalled := false
		patches.ApplyFunc(convertDomToStringInternal,
			func(_ *dom.Node, _ uint16) {
				isCalled = true
				return
			},
		)

		// Act
		ConvertDomToString(root)

		// Assert
		assert.True(t, isCalled)
	})
}

var expected = `
Element(Element {kind: html, Attributes: []})
  Element(Element {kind: head, Attributes: []})
  Element(Element {kind: body, Attributes: []})
    Element(Element {kind: p, Attributes: [{foo bar}]})
      Text(text)
`

func Test_convertDomToStringInternal(t *testing.T) {
	t.Run("if node is nil, return", func(t *testing.T) {
		result := ""

		// Act
		convertDomToStringInternal(nil, 0)

		// Assert
		assert.Equal(t, result, "")
	})

	t.Run("convert node tree", func(t *testing.T) {
		// Arrange
		result = "\n"

		// <html><head></head><body><p foo=bar>text</p></body></html>
		root := dom.NewNode(dom.NodeKind{
			Element: htmlTypes.NewElement(
				"html", []htmlTypes.Attribute{},
			).(*htmlTypes.Element),
		}).(*dom.Node)
		head := dom.NewNode(dom.NodeKind{
			Element: htmlTypes.NewElement(
				"head", []htmlTypes.Attribute{},
			).(*htmlTypes.Element),
		}).(*dom.Node)
		root.FirstChild = head
		body := dom.NewNode(dom.NodeKind{
			Element: htmlTypes.NewElement(
				"body", []htmlTypes.Attribute{},
			).(*htmlTypes.Element),
		}).(*dom.Node)
		head.NextSibling = body
		p := dom.NewNode(dom.NodeKind{
			Element: htmlTypes.NewElement(
				"p", []htmlTypes.Attribute{
					*htmlTypes.NewAttribute("foo", "bar").(*htmlTypes.Attribute),
				},
			).(*htmlTypes.Element),
		}).(*dom.Node)
		body.FirstChild = p
		text := dom.NewNode(dom.NodeKind{
			Text:    "text",
			HasText: true,
		}).(*dom.Node)
		p.FirstChild = text

		// Act
		convertDomToStringInternal(root, 0)

		// Assert
		assert.Equal(t, expected, result)
	})
}
