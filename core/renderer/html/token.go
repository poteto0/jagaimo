package html

type StartTag struct {
	Tag           string
	IsSelfClosing bool
	Attributes    []Attribute
}

type EndTag struct {
	Tag string
}

type EOF int

type HtmlToken struct {
	StartTag StartTag
	EndTag   EndTag
	EOF      EOF
	Rune     rune
}

func newEOFToken() *HtmlToken {
	return &HtmlToken{
		StartTag: StartTag{},
		EndTag:   EndTag{},
		EOF:      EOF(1), // ! not 0
		Rune:     0,
	}
}

func newRuneToken(r rune) *HtmlToken {
	return &HtmlToken{
		StartTag: StartTag{},
		EndTag:   EndTag{},
		EOF:      EOF(0), // not EOF
		Rune:     r,
	}
}

func (token *HtmlToken) IsStartTag() bool {
	return token.StartTag.Tag != ""
}

func (token *HtmlToken) IsEndTag() bool {
	return token.EndTag.Tag != ""
}

func (token *HtmlToken) IsEOF() bool {
	return token.EOF != 0
}

func (token *HtmlToken) IsRune() bool {
	return token.Rune != 0
}
