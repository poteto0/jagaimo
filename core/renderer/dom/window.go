package dom

import "weak"

type IWindow interface {
	Document() *Node
}

type Window struct {
	document *Node
}

func NewWindow() IWindow {
	window := &Window{
		document: NewNode(NodeKind{
			Document: 1, // not zero
		}).(*Node),
	}

	window.document.SetWindow(weak.Make(window))

	return window
}

func (window *Window) Document() *Node {
	return window.document
}
