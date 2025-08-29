package dom

import (
	"weak"

	"github.com/poteto0/jagaimo/core/renderer/html"
)

type NodeKind struct {
	Document int
	Element  *html.Element
	Text     string
}

func (nk *NodeKind) IsDocument() bool {
	return nk.Document != 0
}

func (nk *NodeKind) IsElement() bool {
	return nk.Element != nil
}

func (nk *NodeKind) IsText() bool {
	return nk.Text != ""
}

type INode interface {
	SetWindow(window weak.Pointer[Window])
	Kind() NodeKind

	GetElement() *html.Element
	ElementKind() html.ElementKind
}

type Node struct {
	kind NodeKind
	// escape loop reference
	window      weak.Pointer[Window]
	Parent      weak.Pointer[Node]
	FirstChild  *Node
	LastChild   weak.Pointer[Node]
	NextSibling weak.Pointer[Node]
	PrevSibling *Node
}

func NewNode(kind NodeKind) INode {
	return &Node{
		kind:        kind,
		window:      weak.Pointer[Window]{},
		Parent:      weak.Pointer[Node]{},
		FirstChild:  nil,
		LastChild:   weak.Pointer[Node]{},
		NextSibling: weak.Pointer[Node]{},
		PrevSibling: nil,
	}
}

func (node *Node) SetWindow(window weak.Pointer[Window]) {
	node.window = window
}

func (node *Node) Kind() NodeKind {
	return node.kind
}

func (node *Node) GetElement() *html.Element {
	if node.kind.IsElement() {
		return node.kind.Element
	}

	return nil
}

func (node *Node) ElementKind() html.ElementKind {
	if node.kind.IsElement() {
		return node.kind.Element.Kind()
	}

	return html.NilElement
}
