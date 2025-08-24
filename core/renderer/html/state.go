package html

// refs: https://html.spec.whatwg.org/multipage/parsing.html
type State int

const (
	Data State = iota
	TagOpen
	EndTagOpen
	TagName
	BeforeAttributeName
	AttributeName
	AfterAttributeName
	BeforeAttributeValue
	AttributeValueDoubleQuoted
	AttributeValueSingleQuoted
	AttributeValueUnquoted
	AfterAttributeValueQuoted
	SelfClosingStartTag
	ScriptData
	ScriptDataLessThanSign
	ScriptDataEndTagOpen
	ScriptDataEndTagName
	TemporaryBuffer
)
