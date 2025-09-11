package types

import "errors"

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

	P = ElementKind("p")

	// refs: https://html.spec.whatwg.org/multipage/sections.html#the-h1,-h2,-h3,-h4,-h5,-and-h6-elements
	H1 = ElementKind("h1")
	H2 = ElementKind("h2")

	// refs: https://html.spec.whatwg.org/multipage/text-level-semantics.html#the-a-element
	A = ElementKind("a")
)

func ConvertToElementKind(elementName string) (ElementKind, error) {
	switch elementName {
	case "html":
		return Html, nil
	case "head":
		return Head, nil
	case "style":
		return Style, nil
	case "script":
		return Script, nil
	case "body":
		return Body, nil
	case "p":
		return P, nil
	case "h1":
		return H1, nil
	case "h2":
		return H2, nil
	case "a":
		return A, nil
	default:
		return NilElement, errors.New("unexpected element name")
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
	ele, err := ConvertToElementKind(elementName)
	if err != nil {
		panic(err)
	}

	return &Element{
		kind:       ele,
		attributes: append(attributes, Attribute{}),
	}
}

func (element *Element) Kind() ElementKind {
	return element.kind
}
