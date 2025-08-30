package html

import (
	"weak"

	"github.com/poteto0/jagaimo/core/renderer/dom"
	"github.com/poteto0/jagaimo/core/renderer/html/types"
)

// refs: https://html.spec.whatwg.org/multipage/parsing.html#parsing-main-inhtml
type InsertionMode int

const (
	Initial InsertionMode = iota
	BeforeHtml
	BeforeHead
	InHead
	AfterHead
	InBody
	TextAfterBody
	AfterAfterBody
)

type IHtmlParser interface {
	ConstructTree() *dom.Window
}

type HtmlParser struct {
	window dom.IWindow

	// current state
	mode InsertionMode

	// before state
	// refs: https://html.spec.whatwg.org/multipage/parsing.html#original-insertion-mode
	originalInsertionMode InsertionMode

	// refs: https://html.spec.whatwg.org/multipage/parsing.html#the-stack-of-open-elements
	stackOfOpenElements []*dom.Node

	t IHtmlTokenizer
}

func NewHtmlParser(tokenizer IHtmlTokenizer) IHtmlParser {
	return &HtmlParser{
		window:                dom.NewWindow(),
		mode:                  Initial,
		originalInsertionMode: Initial,
		stackOfOpenElements:   []*dom.Node{},
		t:                     tokenizer,
	}
}

/*
Initial       : Document
  │
BeforeHtml    : Document
  │               └ HtmlElement
BeforeHead    : Document
  │               └ HtmlElement
	│                    └ HeadElement
InHead
  │
AfterHead     : Document
  │               └ HtmlElement
	│                    ├ HeadElement
	│                    └ BodyElement
InBody───┐: Document
  │└───┘     └ HtmlElement
	│                    ├ HeadElement
	│                    └ BodyElement
	│                          └ H1Element
	│                                └ "hello"
AfterBody
  │
AfterAfterBody
  │
EOF
*/

func (parser *HtmlParser) ConstructTree() *dom.Window {
	token := parser.t.Next()
	for token != nil {
		switch parser.mode {
		// not support DOCTYPE
		case Initial:
			if next, _ := parser.parseInitial(token); next != nil {
				token = next
			}

		case BeforeHtml:
			return nil
		case BeforeHead:
			return nil
		case InHead:
			return nil
		case AfterHead:
			return nil
		case InBody:
			return nil
		case TextAfterBody:
			return nil
		case AfterAfterBody:
			return nil
		default:
			panic("unexpected mode")
		}
	}
	return nil
}

// not support DOCTYPE
func (parser *HtmlParser) parseInitial(token *HtmlToken) (next *HtmlToken, IsFinished bool) {
	if parser.mode != Initial {
		panic("unexpected insertion mode")
	}

	// ignore rune token
	if token.IsRune() {
		return parser.t.Next(), false
	}

	parser.mode = BeforeHtml
	return nil, false
}

// take <html>
func (parser *HtmlParser) parseBeforeHtml(token *HtmlToken) (next *HtmlToken, IsFinished bool) {
	if parser.mode != BeforeHtml {
		panic("unexpected insertion mode")
	}

	if r := token.Rune; r != rune(0) {
		if r == ' ' || r == '\n' {
			return parser.t.Next(), false
		}
	}

	if token.IsStartTag() {
		tag, _, attributes := token.StartTag.Take()
		if tag == "html" {
			parser.insertElement(tag, attributes)
			parser.mode = BeforeHead
			return parser.t.Next(), false
		}
	}

	if token.IsEOF() {
		return nil, true
	}

	// auto insert html token
	parser.insertElement("html", []types.Attribute{})
	parser.mode = BeforeHead
	return nil, false
}

func (parser *HtmlParser) insertElement(tag string, attributes []types.Attribute) {
	currentNode := parser.currentNode()
	node := parser.createElementNode(tag, attributes)

	defer func() {
		currentNode.LastChild = weak.Make(node)
		node.Parent = weak.Make(currentNode)
		parser.stackOfOpenElements = append(parser.stackOfOpenElements, node)
	}()

	if currentNode.FirstChild == nil {
		currentNode.FirstChild = node
		return
	}

	lastSibling := currentNode.FirstChild
	for {
		if lastSibling == nil {
			panic("lastSibling shouldn't be nil")
		}

		if lastSibling.NextSibling == nil {
			break
		}

		lastSibling = lastSibling.NextSibling
	}

	lastSibling.NextSibling = node
	node.PrevSibling = weak.Make(lastSibling)
}

func (parser *HtmlParser) currentNode() *dom.Node {
	if len(parser.stackOfOpenElements) == 0 {
		return parser.window.Document()
	}

	if n := parser.stackOfOpenElements[len(parser.stackOfOpenElements)-1]; n != nil {
		return n
	}

	return parser.window.Document()
}

func (parser *HtmlParser) createElementNode(tag string, attributes []types.Attribute) *dom.Node {
	return dom.NewNode(
		dom.NodeKind{
			Element: types.NewElement(
				tag, attributes,
			).(*types.Element),
		},
	).(*dom.Node)
}
