package html

import "github.com/poteto0/jagaimo/core/renderer/html/types"

type StartTag struct {
	Tag           string
	IsSelfClosing bool
	Attributes    []types.Attribute
}

func (st *StartTag) Take() (tag string, isSelfClosing bool, attributes []types.Attribute) {
	return st.Tag, st.IsSelfClosing, st.Attributes
}

type EndTag struct {
	Tag string
}

type EOF int

type HtmlToken struct {
	StartTag *StartTag
	EndTag   *EndTag
	EOF      EOF
	Rune     rune
}

func newEOFToken() *HtmlToken {
	return &HtmlToken{
		StartTag: nil,
		EndTag:   nil,
		EOF:      EOF(1), // ! not 0
		Rune:     0,
	}
}

func newRuneToken(r rune) *HtmlToken {
	return &HtmlToken{
		StartTag: nil,
		EndTag:   nil,
		EOF:      EOF(0), // not EOF
		Rune:     r,
	}
}

func (token *HtmlToken) IsStartTag() bool {
	return token.StartTag != nil
}

func (token *HtmlToken) IsEndTag() bool {
	return token.EndTag != nil
}

func (token *HtmlToken) IsEOF() bool {
	return token.EOF != 0
}

func (token *HtmlToken) IsRune() bool {
	return token.Rune != 0
}
