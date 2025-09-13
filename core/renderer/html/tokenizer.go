package html

import (
	"unicode"

	"github.com/poteto0/jagaimo/core/renderer/html/types"
)

type IHtmlTokenizer interface {
	/*
		tokenize html to make token
			- increment letter one by one
			- change state
			- create token & return it
	*/
	Next() *HtmlToken
}

type HtmlTokenizer struct {
	State       State
	Pos         uint
	ReConsume   bool
	LatestToken *HtmlToken
	Input       []rune
	Buf         string
}

func NewHtmlTokenizer(html string) IHtmlTokenizer {
	return &HtmlTokenizer{
		State:       Data,
		Pos:         0,
		ReConsume:   false,
		LatestToken: nil,
		Input:       []rune(html),
		Buf:         "",
	}
}

func (tokenizer *HtmlTokenizer) Next() *HtmlToken {
	// ! not EOF
	// ! starts w/ last token
	if tokenizer.Pos >= uint(len(tokenizer.Input)) {
		return newEOFToken()
	}

	for {
		r := func() rune {
			if tokenizer.ReConsume {
				return tokenizer.reConsumeInput()
			}
			return tokenizer.consumeNextInput()
		}()

		// TODO: state stays case doc
		switch tokenizer.State {
		// Data ---> TagOpen
		case Data:
			if token := tokenizer.tokenizeData(r); token != nil {
				return token
			}

		// TagOpen ---> EndTagOpen | TagName
		case TagOpen:
			// this just return nil
			tokenizer.tokenizeTagOpen(r)

		// EndTagOpen ---> TagName
		case EndTagOpen:
			if token := tokenizer.tokenizeEndTagOpen(r); token != nil {
				return token
			}

		// TagName ---> BeforeAttributeName | SelfClosingStartTag | Data
		case TagName:
			if token := tokenizer.tokenizeTagName(r); token != nil {
				return token
			}

		// BeforeAttributeName ---> AfterAttributeName | AttributeName
		case BeforeAttributeName:
			if token := tokenizer.tokenizeBeforeAttributeName(r); token != nil {
				return token
			}

		// AttributeName ---> AfterAttributeName | BeforeAttributeValue
		case AttributeName:
			// just return nil
			tokenizer.tokenizeAttributeName(r)

		// AfterAttributeName ---> SelfClosingStartTag | Data | BeforeAttributeValue
		case AfterAttributeName:
			if token := tokenizer.tokenizeAfterAttributeName(r); token != nil {
				return token
			}

		//  BeforeAttributeValue ---> AttributeValueDoubleQuoted | AttributeValueSingleQuoted | AttributeValueUnquoted
		case BeforeAttributeValue:
			// just return nil
			tokenizer.tokenizeBeforeAttributeValue(r)

		// AttributeValueDoubleQuoted ---> AfterAttributeValueQuoted
		case AttributeValueDoubleQuoted:
			if token := tokenizer.tokenizeAttributeValueDoubleQuoted(r); token != nil {
				return token
			}

		// AttributeValueSingleQuoted ---> AfterAttributeValueQuoted
		case AttributeValueSingleQuoted:
			tokenizer.tokenizeAttributeValueSingleQuoted(r)

		// AttributeValueUnquoted ---> BeforeAttributeName
		case AttributeValueUnquoted:
			if token := tokenizer.tokenizeAttributeValueUnquoted(r); token != nil {
				return token
			}

		// AfterAttributeValueQuoted ---> BeforeAttributeName | SelfClosingStartTag | Data
		case AfterAttributeValueQuoted:
			if token := tokenizer.tokenizeAfterAttributeValueQuoted(r); token != nil {
				return token
			}

		// SelfClosingTag ---> Data
		case SelfClosingStartTag:
			if token := tokenizer.tokenizeSelfClosingStartTag(r); token != nil {
				return token
			}

		// ScriptData ---> ScriptDataLessThanSign
		case ScriptData:
			if token := tokenizer.tokenizeScriptData(r); token != nil {
				return token
			}

		// ScriptDataLessThanSign ---> ScriptDataEndTagOpen | ScriptData
		case ScriptDataLessThanSign:
			if token := tokenizer.tokenizeScriptDataLessThanSign(r); token != nil {
				return token
			}

		// ScriptDataEndTagOpen ---> ScriptDataEndTagName | ScriptData
		case ScriptDataEndTagOpen:
			if token := tokenizer.tokenizeScriptDataEndTagOpen(r); token != nil {
				return token
			}

		// ScriptDataEndTagName ---> Data | TemporaryBuffer
		case ScriptDataEndTagName:
			if token := tokenizer.tokenizeScriptDataEndTagName(r); token != nil {
				return token
			}

		// temporary for develop
		case TemporaryBuffer:
			if token := tokenizer.tokenizeTemporaryBuffer(r); token != nil {
				return token
			}

		default:
			panic("unexpected state")
		}
	}
}

// Data ---> Self | TagOpen
func (tokenizer *HtmlTokenizer) tokenizeData(r rune) *HtmlToken {
	if tokenizer.State != Data {
		panic("unexpected state")
	}

	if tokenizer.isEOF() {
		return newEOFToken()
	}

	if r == '<' {
		tokenizer.State = TagOpen
		return nil
	}

	return newRuneToken(r)
}

// TagOpen ---> Self | EndTagOpen | TagName
//
// just return nil
func (tokenizer *HtmlTokenizer) tokenizeTagOpen(r rune) *HtmlToken {
	if tokenizer.State != TagOpen {
		panic("unexpected state")
	}

	if tokenizer.isEOF() {
		return newEOFToken()
	}

	if r == '/' {
		tokenizer.State = EndTagOpen
		return nil
	}

	if isAsciiAlphabetic(r) {
		tokenizer.ReConsume = true
		tokenizer.State = TagName
		tokenizer.createTag(true)
		return nil
	}

	tokenizer.ReConsume = true
	tokenizer.State = Data
	return nil
}

// EndTagOpen ---> Self | TagName
func (tokenizer *HtmlTokenizer) tokenizeEndTagOpen(r rune) *HtmlToken {
	if tokenizer.State != EndTagOpen {
		panic("unexpected state")
	}

	if tokenizer.isEOF() {
		return newEOFToken()
	}

	if isAsciiAlphabetic(r) {
		tokenizer.ReConsume = true
		tokenizer.State = TagName
		tokenizer.createTag(false)
		return nil
	}

	return nil
}

// TagName ---> Self | BeforeAttributeName | SelfClosingStartTag | Data
func (tokenizer *HtmlTokenizer) tokenizeTagName(r rune) *HtmlToken {
	if tokenizer.State != TagName {
		panic("unexpected state")
	}

	if r == ' ' {
		tokenizer.State = BeforeAttributeName
		return nil
	}

	// EX: <img />
	if r == '/' {
		tokenizer.State = SelfClosingStartTag
		return nil
	}

	if r == '>' {
		tokenizer.State = Data
		return tokenizer.takeLastToken()
	}

	if tokenizer.isEOF() {
		return newEOFToken()
	}

	if isAsciiAlphabetic(r) {
		tokenizer.appendTagName(unicode.ToLower(r))
		return nil
	}

	tokenizer.appendTagName(r)
	return nil
}

// BeforeAttributeName ---> AfterAttributeName | AttributeName
func (tokenizer *HtmlTokenizer) tokenizeBeforeAttributeName(r rune) *HtmlToken {
	if tokenizer.State != BeforeAttributeName {
		panic("unexpected state")
	}

	if r == '/' || r == '>' || tokenizer.isEOF() {
		tokenizer.ReConsume = true
		tokenizer.State = AfterAttributeName
		return nil
	}

	tokenizer.ReConsume = true
	tokenizer.State = AttributeName
	tokenizer.startNewAttribute()
	return nil
}

// AttributeName ---> Self | AfterAttributeName | BeforeAttributeValue
//
// just return nil
func (tokenizer *HtmlTokenizer) tokenizeAttributeName(r rune) *HtmlToken {
	if tokenizer.State != AttributeName {
		panic("unexpected state")
	}

	if r == '/' || r == '>' || tokenizer.isEOF() {
		tokenizer.ReConsume = true
		tokenizer.State = AfterAttributeName
		return nil
	}

	if r == '=' {
		tokenizer.State = BeforeAttributeValue
		return nil
	}

	if isAsciiAlphabetic(r) {
		tokenizer.appendAttribute(
			unicode.ToLower(r),
			/*isName = */ true,
		)
		return nil
	}

	tokenizer.appendAttribute(
		r,
		/*isName = */ true,
	)
	return nil
}

// AfterAttributeName ---> Self | SelfClosingStartTag | Data | BeforeAttributeValue
func (tokenizer *HtmlTokenizer) tokenizeAfterAttributeName(r rune) *HtmlToken {
	if tokenizer.State != AfterAttributeName {
		panic("unexpected state")
	}

	if r == ' ' {
		return nil
	}

	if r == '/' {
		tokenizer.State = SelfClosingStartTag
		return nil
	}

	if r == '>' {
		tokenizer.State = Data
		return tokenizer.takeLastToken()
	}

	if r == '=' {
		tokenizer.State = BeforeAttributeValue
		return nil
	}

	if tokenizer.isEOF() {
		return newEOFToken()
	}

	tokenizer.ReConsume = true
	tokenizer.State = Data
	tokenizer.startNewAttribute()
	return nil
}

// BeforeAttributeValue ---> Self | AttributeValueDoubleQuoted | AttributeValueSingleQuoted | AttributeValueUnquoted
//
// just return nil
func (tokenizer *HtmlTokenizer) tokenizeBeforeAttributeValue(r rune) *HtmlToken {
	if tokenizer.State != BeforeAttributeValue {
		panic("unexpected state")
	}

	if r == ' ' {
		return nil
	}

	if r == '"' {
		tokenizer.State = AttributeValueDoubleQuoted
		return nil
	}

	if r == '\'' {
		tokenizer.State = AttributeValueSingleQuoted
		return nil
	}

	tokenizer.ReConsume = true
	tokenizer.State = AttributeValueUnquoted
	return nil
}

// AttributeValueDoubleQuoted ---> AfterAttributeValueQuoted
func (tokenizer *HtmlTokenizer) tokenizeAttributeValueDoubleQuoted(r rune) *HtmlToken {
	if tokenizer.State != AttributeValueDoubleQuoted {
		panic("unexpected state")
	}

	if r == '"' {
		tokenizer.State = AfterAttributeValueQuoted
		return nil
	}

	if tokenizer.isEOF() {
		return newEOFToken()
	}

	tokenizer.appendAttribute(
		r,
		/*isName = */ false,
	)
	return nil
}

// AttributeValueSingleQuoted ---> Self | AfterAttributeValueQuoted
func (tokenizer *HtmlTokenizer) tokenizeAttributeValueSingleQuoted(r rune) *HtmlToken {
	if tokenizer.State != AttributeValueSingleQuoted {
		panic("unexpected state")
	}

	if r == '\'' {
		tokenizer.State = AfterAttributeValueQuoted
		return nil
	}

	if tokenizer.isEOF() {
		return newEOFToken()
	}

	tokenizer.appendAttribute(
		r,
		/*isName = */ false,
	)
	return nil
}

// AttributeValueUnquoted ---> Self | BeforeAttributeName
func (tokenizer *HtmlTokenizer) tokenizeAttributeValueUnquoted(r rune) *HtmlToken {
	if tokenizer.State != AttributeValueUnquoted {
		panic("unexpected state")
	}

	if r == ' ' {
		tokenizer.State = BeforeAttributeName
		return nil
	}

	if r == '>' {
		tokenizer.State = Data
		return tokenizer.takeLastToken()
	}

	tokenizer.appendAttribute(
		r,
		/*isName = */ false,
	)
	return nil
}

// AfterAttributeValueQuoted ---> Self | BeforeAttributeName | SelfClosingStartTag | Data
func (tokenizer *HtmlTokenizer) tokenizeAfterAttributeValueQuoted(r rune) *HtmlToken {
	if tokenizer.State != AfterAttributeValueQuoted {
		panic("unexpected state")
	}

	if r == ' ' {
		tokenizer.State = BeforeAttributeName
		return nil
	}

	if r == '/' {
		tokenizer.State = SelfClosingStartTag
		return nil
	}

	if r == '>' {
		tokenizer.State = Data
		return tokenizer.takeLastToken()
	}

	if tokenizer.isEOF() {
		return newEOFToken()
	}

	tokenizer.ReConsume = true
	tokenizer.State = BeforeAttributeName
	return nil
}

// SelfClosingTag ---> Self | Data
func (tokenizer *HtmlTokenizer) tokenizeSelfClosingStartTag(r rune) *HtmlToken {
	if tokenizer.State != SelfClosingStartTag {
		panic("unexpected state")
	}

	if r == '>' {
		tokenizer.setSelfClosingFlag()
		tokenizer.State = Data
		return tokenizer.takeLastToken()
	}

	if tokenizer.isEOF() {
		return newEOFToken()
	}

	return nil
}

// ScriptData ---> Self | ScriptDataLessThanSign
func (tokenizer *HtmlTokenizer) tokenizeScriptData(r rune) *HtmlToken {
	if tokenizer.State != ScriptData {
		panic("unexpected state")
	}

	if r == '<' {
		tokenizer.State = ScriptDataLessThanSign
		return nil
	}

	if tokenizer.isEOF() {
		return newEOFToken()
	}

	return newRuneToken(r)
}

// ScriptDataLessThanSign ---> ScriptDataEndTagOpen | ScriptData
func (tokenizer *HtmlTokenizer) tokenizeScriptDataLessThanSign(r rune) *HtmlToken {
	if tokenizer.State != ScriptDataLessThanSign {
		panic("unexpected state")
	}

	if r == '/' {
		// reset buffer
		tokenizer.Buf = ""
		tokenizer.State = ScriptDataEndTagOpen
		return nil
	}

	tokenizer.ReConsume = true
	tokenizer.State = ScriptData
	return newRuneToken('<')
}

// ScriptDataEndTagOpen ---> ScriptDataEndTagName | ScriptData
func (tokenizer *HtmlTokenizer) tokenizeScriptDataEndTagOpen(r rune) *HtmlToken {
	if tokenizer.State != ScriptDataEndTagOpen {
		panic("unexpected state")
	}

	if isAsciiAlphabetic(r) {
		tokenizer.ReConsume = true
		tokenizer.State = ScriptDataEndTagName
		tokenizer.createTag(false)
		return nil
	}

	tokenizer.ReConsume = true
	tokenizer.State = ScriptData
	// return "</", in the specifications
	return newRuneToken('<')
}

// ScriptDataEndTagName ---> Self | Data | TemporaryBuffer
func (tokenizer *HtmlTokenizer) tokenizeScriptDataEndTagName(r rune) *HtmlToken {
	if tokenizer.State != ScriptDataEndTagName {
		panic("unexpected state")
	}

	if r == '>' {
		tokenizer.State = Data
		return tokenizer.takeLastToken()
	}

	if isAsciiAlphabetic(r) {
		tokenizer.Buf += string(r)
		tokenizer.appendTagName(unicode.ToLower(r))
		return nil
	}

	tokenizer.State = TemporaryBuffer
	tokenizer.Buf = "</" + tokenizer.Buf
	tokenizer.Buf += string(r)
	return nil
}

// TemporaryBuffer ---> Self | ScriptData
func (tokenizer *HtmlTokenizer) tokenizeTemporaryBuffer(r rune) *HtmlToken {
	if tokenizer.State != TemporaryBuffer {
		panic("unexpected state")
	}

	tokenizer.ReConsume = true

	if len(tokenizer.Buf) == 0 {
		tokenizer.State = ScriptData
		return nil
	}

	// delete first letter
	rr := []rune(tokenizer.Buf)
	tokenizer.Buf = string(rr[1:])
	return newRuneToken(r)
}

func (tokenizer *HtmlTokenizer) consumeNextInput() rune {
	r := tokenizer.Input[tokenizer.Pos]
	tokenizer.Pos++
	return r
}

func (tokenizer *HtmlTokenizer) reConsumeInput() rune {
	tokenizer.ReConsume = false
	tokenizer.Pos--
	return tokenizer.consumeNextInput()
}

func (tokenizer *HtmlTokenizer) isEOF() bool {
	return tokenizer.Pos > uint(len(tokenizer.Input))
}

func (tokenizer *HtmlTokenizer) createTag(isStartTagToken bool) {
	if isStartTagToken {
		tokenizer.LatestToken = &HtmlToken{
			StartTag: &StartTag{
				Tag:           "",
				IsSelfClosing: false,
				Attributes:    []types.Attribute{},
			},
		}
		return
	}

	tokenizer.LatestToken = &HtmlToken{
		EndTag: &EndTag{
			Tag: "",
		},
	}
}

func (tokenizer *HtmlTokenizer) appendTagName(r rune) {
	if tokenizer.LatestToken == nil {
		panic("unexpected nil latest token")
	}

	token := tokenizer.LatestToken
	if token.IsStartTag() {
		token.StartTag.Tag += string(r)
		return
	}

	if token.IsEndTag() {
		token.EndTag.Tag += string(r)
		return
	}

	panic("unexpected latest token, only expect StartTag or EndTag")
}

func (tokenizer *HtmlTokenizer) takeLastToken() *HtmlToken {
	if tokenizer.LatestToken == nil {
		panic("unexpected nil latest token")
	}

	token := tokenizer.LatestToken
	tokenizer.LatestToken = nil

	return token
}

func (tokenizer *HtmlTokenizer) startNewAttribute() {
	if tokenizer.LatestToken == nil {
		panic("unexpected nil latest token")
	}

	token := tokenizer.LatestToken
	if token.IsStartTag() {
		token.StartTag.Attributes = append(token.StartTag.Attributes, types.Attribute{})
		return
	}

	panic("unexpected latest token, only expect StartTag")
}

func (tokenizer *HtmlTokenizer) appendAttribute(r rune, isName bool) {
	if tokenizer.LatestToken == nil {
		panic("unexpected nil latest token")
	}

	token := tokenizer.LatestToken
	if token.IsStartTag() {
		attrLens := len(token.StartTag.Attributes)
		if attrLens <= 0 {
			panic("unexpected empty attribute list")
		}

		token.StartTag.Attributes[len(token.StartTag.Attributes)-1].AddRune(r, isName)
		return
	}

	panic("unexpected latest token, only expect StartTag")
}

func (tokenizer *HtmlTokenizer) setSelfClosingFlag() {
	if tokenizer.LatestToken == nil {
		panic("unexpected nil latest token")
	}

	token := tokenizer.LatestToken
	if token.IsStartTag() {
		token.StartTag.IsSelfClosing = true
		return
	}

	panic("unexpected latest token, only expect StartTag")
}
