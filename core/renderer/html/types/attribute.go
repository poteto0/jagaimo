package types

type IAttribute interface {
	/*
		add rune -> name/value
	*/
	AddRune(r rune, isName bool)

	Name() string
	Value() string
}

type Attribute struct {
	name  string
	value string
}

func NewAttribute(name string, value string) IAttribute {
	return &Attribute{
		name:  name,
		value: value,
	}
}

func (attr *Attribute) AddRune(r rune, isName bool) {
	if isName {
		attr.name += string(r)
		return
	}

	attr.value += string(r)
}

func (attr *Attribute) Name() string {
	return attr.name
}

func (attr *Attribute) Value() string {
	return attr.value
}
