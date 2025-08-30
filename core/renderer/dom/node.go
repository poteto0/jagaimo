package dom

import (
	"weak"

	htmlTypes "github.com/poteto0/jagaimo/core/renderer/html/types"
)

type NodeKind struct {
	Document int
	Element  *htmlTypes.Element
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

	GetElement() *htmlTypes.Element
	ElementKind() htmlTypes.ElementKind
}

type Node struct {
	kind NodeKind
	// escape loop reference
	window      weak.Pointer[Window]
	Parent      weak.Pointer[Node]
	FirstChild  *Node
	LastChild   weak.Pointer[Node]
	NextSibling *Node
	PrevSibling weak.Pointer[Node]
}

func NewNode(kind NodeKind) INode {
	return &Node{
		kind:        kind,
		window:      weak.Pointer[Window]{},
		Parent:      weak.Pointer[Node]{},
		FirstChild:  nil,
		LastChild:   weak.Pointer[Node]{},
		NextSibling: nil,
		PrevSibling: weak.Pointer[Node]{},
	}
}

func (node *Node) SetWindow(window weak.Pointer[Window]) {
	node.window = window
}

func (node *Node) Kind() NodeKind {
	return node.kind
}

func (node *Node) GetElement() *htmlTypes.Element {
	if node.kind.IsElement() {
		return node.kind.Element
	}

	return nil
}

func (node *Node) ElementKind() htmlTypes.ElementKind {
	if node.kind.IsElement() {
		return node.kind.Element.Kind()
	}

	return htmlTypes.NilElement
}
