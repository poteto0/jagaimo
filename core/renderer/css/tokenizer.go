package css

type ICssTokenizer interface {
	Next() CssToken
}

type CssTokenizer struct {
	pos   uint
	input []rune
}

func NewCssTokenizer(input string) ICssTokenizer {
	return &CssTokenizer{
		pos:   0,
		input: []rune(input),
	}
}

func (tokenizer *CssTokenizer) Next() CssToken {
	for {
		if tokenizer.pos >= uint(len(tokenizer.input)) {
			return NilToken
		}

		defer func() {
			tokenizer.pos += 1
		}()

		r := tokenizer.input[tokenizer.pos]
		token, isSkip := tokenizer.decideToken(r)
		if isSkip {
			continue
		}

		tokenizer.pos += 1
		return token
	}
}

func (tokenizer *CssTokenizer) decideToken(r rune) (token CssToken, isSkip bool) {
	switch r {
	case '(':
		return OpenParenthesis, false
	case ')':
		return CloseParenthesis, false
	case ',':
		return DelimComma, false
	case '.':
		return DelimPeriod, false
	case ':':
		return Colon, false
	case ';':
		return Semicolon, false
	case '{':
		return OpenCurly, false
	case '}':
		return CloseCurly, false
	case ' ', '\n':
		return NilToken, true
	default:
		return NilToken, false
	}
}
