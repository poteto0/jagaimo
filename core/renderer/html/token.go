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

type HtmlToken interface {
	StartTag | EndTag | EOF | *rune
}
