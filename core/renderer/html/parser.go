package html

import (
	"fmt"
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
	Text
	AfterHead
	InBody
	TextAfterBody
	AfterBody
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
			next, isFinished := parser.parseBeforeHtml(token)
			if isFinished {
				return parser.window.(*dom.Window)
			}

			if next != nil {
				token = next
			}

		case BeforeHead:
			next, isFinished := parser.parseBeforeHead(token)
			if isFinished {
				return parser.window.(*dom.Window)
			}

			if next != nil {
				token = next
			}

		case InHead:
			next, isFinished := parser.parseInHead(token)
			if isFinished {
				return parser.window.(*dom.Window)
			}

			if next != nil {
				token = next
			}

		case AfterHead:
			next, isFinished := parser.parseAfterHead(token)
			if isFinished {
				return parser.window.(*dom.Window)
			}

			if next != nil {
				token = next
			}

		case InBody:
			return nil
		case TextAfterBody:
			return nil
		case AfterBody:
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

func (parser *HtmlParser) parseBeforeHead(token *HtmlToken) (next *HtmlToken, IsFinished bool) {
	if parser.mode != BeforeHead {
		panic("unexpected insertion mode")
	}

	if r := token.Rune; r != rune(0) {
		if r == ' ' || r == '\n' {
			return parser.t.Next(), false
		}
	}

	if token.IsStartTag() {
		tag, _, attributes := token.StartTag.Take()
		if tag == "head" {
			parser.insertElement(tag, attributes)
			parser.mode = InHead
			return parser.t.Next(), false
		}
	}

	if token.IsEOF() {
		return nil, true
	}

	// auto insert head token
	parser.insertElement("head", []types.Attribute{})
	parser.mode = InHead
	return nil, false
}

func (parser *HtmlParser) parseInHead(token *HtmlToken) (next *HtmlToken, IsFinished bool) {
	if parser.mode != InHead {
		panic("unexpected insertion mode")
	}

	if r := token.Rune; r != rune(0) {
		if r == ' ' || r == '\n' {
			// if currentNode is not Text, do nothing
			parser.insertRune(r)
			return parser.t.Next(), false
		}
	}

	if token.IsStartTag() {
		tag, _, attributes := token.StartTag.Take()
		if tag == "style" || tag == "script" {
			parser.insertElement(tag, attributes)
			parser.originalInsertionMode = parser.mode
			parser.mode = Text
			return parser.t.Next(), false
		}

		// deal w hanging infinite loop by omission of <head>
		// <head> is omitted, cannot move AfterHead from Input
		if tag == "body" {
			parser.popUntil(types.Head)
			parser.mode = AfterHead
			return nil, false
		}

		// !skip not supported element kind
		// !user can use on html, but doesn't work.
		// EX) <meta> <title>
		parser.popUntil(types.Head)
		parser.mode = AfterHead
		return nil, false
	}

	if token.IsEndTag() {
		tag := token.EndTag.Tag
		if tag == "head" {
			parser.mode = AfterHead
			parser.popUntil(types.Head)
			return parser.t.Next(), false
		}
	}

	if token.IsEOF() {
		return nil, true
	}

	// !skip not supported element kind
	// !user can use on html, but doesn't work.
	// EX) <meta> <title>
	return parser.t.Next(), false
}

func (parser *HtmlParser) parseAfterHead(token *HtmlToken) (next *HtmlToken, IsFinished bool) {
	if parser.mode != AfterHead {
		panic("unexpected insertion mode")
	}

	if r := token.Rune; r != rune(0) {
		if r == ' ' || r == '\n' {
			parser.insertRune(r)
			return parser.t.Next(), false
		}
	}

	if token.IsStartTag() {
		tag, _, attributes := token.StartTag.Take()
		if tag == "body" {
			parser.insertElement(tag, attributes)
			parser.mode = InBody
			return parser.t.Next(), false
		}
	}

	if token.IsEOF() {
		return nil, true
	}

	// if not has body, auto append body to DOM
	parser.insertElement("body", []types.Attribute{})
	parser.mode = InBody
	return nil, false
}

func (parser *HtmlParser) parseInBody(token *HtmlToken) (next *HtmlToken, IsFinished bool) {
	if parser.mode != InBody {
		panic("unexpected insertion mode")
	}

	if token.IsEndTag() {
		switch token.EndTag.Tag {
		case "body":
			parser.mode = AfterBody
			token := parser.t.Next()
			// if failed parse, skip token
			if !parser.hasKindInStack(types.Body) {
				return token, false
			}
			parser.popUntil(types.Body)
			return token, false

		// if skipped body
		case "html":
			// auto inserted body
			if parser.tryPopCurrentNode(types.Body) {
				parser.mode = AfterBody

				if !parser.tryPopCurrentNode(types.Html) {
					panic("unexpected html element")
				}

				return nil, false
			}
			return parser.t.Next(), false

		// TODO: other end tags
		default:
			return parser.t.Next(), false
		}
	}

	// TODO: star tag

	if token.IsEOF() {
		return nil, true
	}

	return parser.t.Next(), false
}

// insert element node into node's last child
//   - link parent to child
func (parser *HtmlParser) insertElement(tag string, attributes []types.Attribute) {
	currentNode := parser.currentNode()
	node := createElementNode(tag, attributes)

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

func (parser *HtmlParser) insertRune(r rune) {
	currentNode := parser.currentNode()

	if currentNode.Kind.IsText() {
		currentNode.Kind.Text += string(r)
		return
	}

	if r == '\n' || r == ' ' {
		return
	}

	node := createRune(r)

	defer func() {
		currentNode.LastChild = weak.Make(node)
		node.Parent = weak.Make(currentNode)
		parser.stackOfOpenElements = append(parser.stackOfOpenElements, node)
	}()

	if currentNode.FirstChild == nil {
		currentNode.FirstChild = node
		return
	}

	currentNode.FirstChild.NextSibling = node
}

// pop stacked open all elements until target element kind
func (parser *HtmlParser) popUntil(kind types.ElementKind) {
	if !parser.hasKindInStack(kind) {
		panic(fmt.Sprintf("unexpected stack doesn't have %s", kind))
	}

	for {
		node := parser.popCurrentNode()
		if node.ElementKind() == kind {
			return
		}
	}
}

func (parser *HtmlParser) hasKindInStack(kind types.ElementKind) bool {
	if len(parser.stackOfOpenElements) == 0 {
		return false
	}

	for _, n := range parser.stackOfOpenElements {
		if n.ElementKind() == kind {
			return true
		}
	}

	return false
}

func (parser *HtmlParser) popCurrentNode() *dom.Node {
	if len(parser.stackOfOpenElements) == 0 {
		panic("unexpected empty stack")
	}

	node := parser.currentNode()
	parser.stackOfOpenElements = parser.stackOfOpenElements[:len(parser.stackOfOpenElements)-1]
	return node
}

// try to pop by element Kind.
// if currentNode is target element, pop & return true
func (parser *HtmlParser) tryPopCurrentNode(kind types.ElementKind) bool {
	if len(parser.stackOfOpenElements) == 0 {
		return false
	}

	node := parser.currentNode()
	if node.ElementKind() == kind {
		parser.popCurrentNode()
		return true
	}

	return false
}
