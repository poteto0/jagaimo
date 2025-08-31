package dom

import (
	"weak"

	htmlTypes "github.com/poteto0/jagaimo/core/renderer/html/types"
)

type NodeKind struct {
	Document int
	Element  *htmlTypes.Element
	Text     string
	HasText  bool
}

func (nk *NodeKind) IsDocument() bool {
	return nk.Document != 0
}

func (nk *NodeKind) IsElement() bool {
	return nk.Element != nil
}

func (nk *NodeKind) IsText() bool {
	return nk.HasText
}

type INode interface {
	SetWindow(window weak.Pointer[Window])

	GetElement() *htmlTypes.Element
	ElementKind() htmlTypes.ElementKind
}

type Node struct {
	Kind NodeKind
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
		Kind:        kind,
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

func (node *Node) GetElement() *htmlTypes.Element {
	if node.Kind.IsElement() {
		return node.Kind.Element
	}

	return nil
}

func (node *Node) ElementKind() htmlTypes.ElementKind {
	if node.Kind.IsElement() {
		return node.Kind.Element.Kind()
	}

	return htmlTypes.NilElement
}
