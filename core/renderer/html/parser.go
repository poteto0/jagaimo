package html

import (
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
	stackOfOpenElements []*types.Element

	t IHtmlTokenizer
}

func NewHtmlParser(tokenizer IHtmlTokenizer) IHtmlParser {
	return &HtmlParser{
		window:                dom.NewWindow(),
		mode:                  Initial,
		originalInsertionMode: Initial,
		stackOfOpenElements:   []*types.Element{},
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
			if next := parser.parseInitial(token); next != nil {
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
func (parser *HtmlParser) parseInitial(token *HtmlToken) *HtmlToken {
	if parser.mode != Initial {
		panic("unexpected insertion mode")
	}

	// ignore rune token
	if token.IsRune() {
		return parser.t.Next()
	}

	parser.mode = BeforeHtml
	return nil
}
