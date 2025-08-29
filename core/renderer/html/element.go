package html

type ElementKind string

const (
	NilElement = ElementKind("")

	// refs: https://html.spec.whatwg.org/multipage/semantics.html#the-html-element
	Html = ElementKind("html")

	// refs: https://html.spec.whatwg.org/multipage/semantics.html#the-head-element
	Head = ElementKind("head")

	// refs: https://html.spec.whatwg.org/multipage/semantics.html#the-style-element
	Style = ElementKind("style")

	// refs: https://html.spec.whatwg.org/multipage/scripting.html#the-script-element
	Script = ElementKind("script")

	// refs: https://html.spec.whatwg.org/multipage/sections.html#the-body-element
	Body = ElementKind("body")
)

func convertElementKind(elementName string) ElementKind {
	switch elementName {
	case "html":
		return Html
	case "head":
		return Head
	case "style":
		return Style
	case "script":
		return Script
	case "body":
		return Body
	default:
		panic("unexpected element name")
	}
}

type IElement interface {
	Kind() ElementKind
}

type Element struct {
	kind       ElementKind
	attributes []Attribute
}

func NewElement(elementName string, attributes []Attribute) IElement {
	return &Element{
		kind:       convertElementKind(elementName),
		attributes: append(attributes, Attribute{}),
	}
}

func (element *Element) Kind() ElementKind {
	return element.kind
}
