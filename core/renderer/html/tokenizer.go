package html

import (
	"iter"
)

type IHtmlTokenizer interface {
	Iter() iter.Seq[*HtmlToken]

	consumeNextInput() rune
	isEOF() bool
	createTag(isStartTagToken bool)
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
			r := tokenizer.consumeNextInput()

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
			}
		}
	}
}

func (tokenizer *HtmlTokenizer) consumeNextInput() rune {
	r := tokenizer.Input[tokenizer.Pos]
	tokenizer.Pos++
	return r
}

func (tokenizer *HtmlTokenizer) isEOF() bool {
	return tokenizer.Pos > uint(len(tokenizer.Input))
}

func isAsciiAlphabetic(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z')
}

func (tokenizer *HtmlTokenizer) createTag(isStartTagToken bool) {
	if isStartTagToken {
		tokenizer.LatestToken = &HtmlToken{
			StartTag: StartTag{
				Tag:           "",
				IsSelfClosing: false,
				Attributes:    []Attribute{},
			},
		}
		return
	}

	tokenizer.LatestToken = &HtmlToken{
		EndTag: EndTag{
			Tag: "",
		},
	}
}
