package html

import (
	"iter"
	"unicode"
)

type IHtmlTokenizer interface {
	Iter() iter.Seq[*HtmlToken]

	consumeNextInput() rune
	reConsumeInput() rune
	isEOF() bool
	createTag(isStartTagToken bool)
	appendTagName(r rune)
	takeLastToken() *HtmlToken
	startNewAttribute()
	appendAttribute(r rune, isName bool)
	setSelfClosingFlag()
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

func (tokenizer *HtmlTokenizer) Iter() iter.Seq[*HtmlToken] {
	return func(yield func(*HtmlToken) bool) {
		// ! not EOF
		// ! starts w/ last token
		if tokenizer.Pos >= uint(len(tokenizer.Input)) {
			yield(nil)
			return
		}

		for {
			r := func() rune {
				if tokenizer.ReConsume {
					return tokenizer.reConsumeInput()
				}
				return tokenizer.consumeNextInput()
			}()

			switch tokenizer.State {
			case Data:
				if r == '<' {
					tokenizer.State = TagOpen
					continue
				}

				if tokenizer.isEOF() {
					yield(newEOFToken())
					return
				}

				yield(newRuneToken(r))
				return

			case TagOpen:
				if r == '/' {
					tokenizer.State = EndTagOpen
					continue
				}

				if isAsciiAlphabetic(r) {
					tokenizer.ReConsume = true
					tokenizer.State = TagName
					tokenizer.createTag(true)
					continue
				}

				if tokenizer.isEOF() {
					yield(newEOFToken())
					return
				}

				tokenizer.ReConsume = true
				tokenizer.State = Data
				continue

			case EndTagOpen:
				if tokenizer.isEOF() {
					yield(newEOFToken())
					return
				}

				if isAsciiAlphabetic(r) {
					tokenizer.ReConsume = true
					tokenizer.State = TagName
					tokenizer.createTag(false)
					continue
				}

			case TagName:
				if r == ' ' {
					tokenizer.State = BeforeAttributeName
					continue
				}

				// <img />
				if r == '/' {
					tokenizer.State = SelfClosingStartTag
					continue
				}

				if r == '>' {
					tokenizer.State = Data
					yield(tokenizer.takeLastToken())
					return
				}

				if isAsciiAlphabetic(r) {
					tokenizer.appendTagName(unicode.ToLower(r))
					continue
				}

				if tokenizer.isEOF() {
					yield(newEOFToken())
					return
				}

				tokenizer.appendTagName(r)

			case BeforeAttributeName:
				if r == '/' || r == '>' || tokenizer.isEOF() {
					tokenizer.ReConsume = true
					tokenizer.State = AfterAttributeName
					continue
				}

				tokenizer.ReConsume = true
				tokenizer.State = AttributeName
				tokenizer.startNewAttribute()
				continue

			case AttributeName:
				if r == '/' || r == '>' || tokenizer.isEOF() {
					tokenizer.ReConsume = true
					tokenizer.State = AfterAttributeName
					continue
				}

				if r == '=' {
					tokenizer.State = BeforeAttributeValue
					continue
				}

				if isAsciiAlphabetic(r) {
					tokenizer.appendAttribute(
						unicode.ToLower(r),
						/*isName = */ true,
					)
				}

				tokenizer.appendAttribute(
					r,
					/*isName = */ true,
				)

			case AfterAttributeName:
				if r == ' ' {
					continue
				}

				if r == '/' {
					tokenizer.State = SelfClosingStartTag
					continue
				}

				if r == '>' {
					tokenizer.State = Data
					yield(tokenizer.takeLastToken())
					return
				}

				if r == '=' {
					tokenizer.State = BeforeAttributeValue
					continue
				}

				if tokenizer.isEOF() {
					yield(newEOFToken())
					return
				}

				tokenizer.ReConsume = true
				tokenizer.State = Data
				tokenizer.startNewAttribute()
				continue

			case BeforeAttributeValue:
				if r == ' ' {
					continue
				}

				if r == '"' {
					tokenizer.State = AttributeValueDoubleQuoted
					continue
				}

				if r == '\'' {
					tokenizer.State = AttributeValueSingleQuoted
					continue
				}

				tokenizer.ReConsume = true
				tokenizer.State = AttributeValueUnquoted
				continue

			case AttributeValueDoubleQuoted:
				if r == '"' {
					tokenizer.State = AfterAttributeValueQuoted
					continue
				}

				if tokenizer.isEOF() {
					yield(newEOFToken())
					return
				}

				tokenizer.appendAttribute(
					r,
					/*isName = */ false,
				)

			case AttributeValueSingleQuoted:
				if r == '\'' {
					tokenizer.State = AfterAttributeValueQuoted
					continue
				}

				if tokenizer.isEOF() {
					yield(newEOFToken())
					return
				}

				tokenizer.appendAttribute(
					r,
					/*isName = */ false,
				)

			case AttributeValueUnquoted:
				if r == ' ' {
					tokenizer.State = BeforeAttributeName
					continue
				}

				if r == '>' {
					tokenizer.State = Data
					yield(tokenizer.takeLastToken())
					return
				}

				if tokenizer.isEOF() {
					yield(newEOFToken())
					return
				}

				tokenizer.appendAttribute(
					r,
					/*isName = */ false,
				)

			case AfterAttributeValueQuoted:
				if r == ' ' {
					tokenizer.State = BeforeAttributeName
					continue
				}

				if r == '/' {
					tokenizer.State = SelfClosingStartTag
					continue
				}

				if r == '>' {
					tokenizer.State = Data
					yield(tokenizer.takeLastToken())
					return
				}

				if tokenizer.isEOF() {
					yield(newEOFToken())
					return
				}

				tokenizer.ReConsume = true
				tokenizer.State = BeforeAttributeName
				continue

			case SelfClosingStartTag:
				if r == '>' {
					tokenizer.setSelfClosingFlag()
					tokenizer.State = Data
					yield(tokenizer.takeLastToken())
					return
				}

				if tokenizer.isEOF() {
					yield(newEOFToken())
					return
				}

			case ScriptData:
				if r == '<' {
					tokenizer.State = ScriptDataLessThanSign
					continue
				}

				if tokenizer.isEOF() {
					yield(newEOFToken())
					return
				}

				yield(newRuneToken(r))
				return

			case ScriptDataLessThanSign:
				if r == '/' {
					// reset buffer
					tokenizer.Buf = ""
					tokenizer.State = ScriptDataEndTagOpen
					continue
				}

				tokenizer.ReConsume = true
				tokenizer.State = ScriptData
				yield(newRuneToken('<'))
				return

			case ScriptDataEndTagOpen:
				if isAsciiAlphabetic(r) {
					tokenizer.ReConsume = true
					tokenizer.State = ScriptDataEndTagName
					tokenizer.createTag(false)
					continue
				}

				tokenizer.ReConsume = true
				tokenizer.State = ScriptData
				// return "</", in the specifications
				yield(newRuneToken('<'))
				return

			case ScriptDataEndTagName:
				if r == '>' {
					tokenizer.State = Data
					yield(tokenizer.takeLastToken())
					return
				}

				if isAsciiAlphabetic(r) {
					tokenizer.Buf += string(r)
					tokenizer.appendTagName(unicode.ToLower(r))
					continue
				}

				tokenizer.State = TemporaryBuffer
				tokenizer.Buf = "</" + tokenizer.Buf
				tokenizer.Buf += string(r)
				continue

			// temporary for develop
			case TemporaryBuffer:
				tokenizer.ReConsume = true

				if len(tokenizer.Buf) == 0 {
					tokenizer.State = ScriptData
					continue
				}

				// delete first letter
				rr := []rune(tokenizer.Buf)
				if len(rr) <= 0 {
					panic("unexpected empty buffer")
				}
				tokenizer.Buf = string(rr[1:])
				yield(newRuneToken(r))
				return

			default:
				panic("unexpected state")
			}
		}
	}
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
				Attributes:    []Attribute{},
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
		token.StartTag.Attributes = append(token.StartTag.Attributes, Attribute{})
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

func isAsciiAlphabetic(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z')
}
