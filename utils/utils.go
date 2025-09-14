package utils

import (
	"fmt"
	"strings"

	"github.com/poteto0/jagaimo/core/renderer/dom"
)

var result = ""

func ConvertDomToString(root *dom.Node) string {
	result = "\n"
	convertDomToStringInternal(root, 0)
	return result
}

func convertDomToStringInternal(node *dom.Node, depth uint16) {
	if node == nil {
		return
	}

	result += strings.Repeat("  ", int(depth))
	result += convertNodeKindToString(node.Kind)
	result += "\n"

	convertDomToStringInternal(node.FirstChild, depth+1)
	convertDomToStringInternal(node.NextSibling, depth)
}

func convertNodeKindToString(kind dom.NodeKind) string {
	if kind.IsElement() {
		return fmt.Sprintf(
			"Element(Element {kind: %s, Attributes: %v})",
			kind.Element.Kind(),
			kind.Element.Attributes(),
		)
	}

	if kind.IsText() {
		return fmt.Sprintf(
			"Text(%s)",
			kind.Text,
		)
	}

	return ""
}
