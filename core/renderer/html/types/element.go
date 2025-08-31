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
