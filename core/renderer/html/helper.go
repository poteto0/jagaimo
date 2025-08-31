package html

import (
	"github.com/poteto0/jagaimo/core/renderer/dom"
	"github.com/poteto0/jagaimo/core/renderer/html/types"
)

func createElementNode(tag string, attributes []types.Attribute) *dom.Node {
	return dom.NewNode(
		dom.NodeKind{
			Element: types.NewElement(
				tag, attributes,
			).(*types.Element),
		},
	).(*dom.Node)
}

func createRune(r rune) *dom.Node {
	return dom.NewNode(
		dom.NodeKind{
			Text:    string(r),
			HasText: true,
		},
	).(*dom.Node)
}

func isAsciiAlphabetic(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z')
}
