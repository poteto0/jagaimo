package css

type CssToken uint16

const (
	NilToken CssToken = iota

	// refs: https://www.w3.org/TR/css-syntax-3/#typedef-hash-token
	HashToken

	// refs: https://www.w3.org/TR/css-syntax-3/#typedef-delim-token
	DelimComma
	DelimPeriod

	// refs: https://www.w3.org/TR/css-syntax-3/#typedef-number-token
	Number

	// refs: https://www.w3.org/TR/css-syntax-3/#typedef-colon-token
	Colon

	// refs: https://www.w3.org/TR/css-syntax-3/#typedef-semicolon-token
	Semicolon

	// refs: https://www.w3.org/TR/css-syntax-3/#tokendef-open-paren
	OpenParenthesis

	// refs: https://www.w3.org/TR/css-syntax-3/#tokendef-close-paren
	CloseParenthesis

	// refs: https://www.w3.org/TR/css-syntax-3/#tokendef-open-curly
	OpenCurly

	// refs: https://www.w3.org/TR/css-syntax-3/#tokendef-close-curly
	CloseCurly

	// refs: https://www.w3.org/TR/css-syntax-3/#typedef-ident-token
	Ident

	// refs: https://www.w3.org/TR/css-syntax-3/#typedef-string-token
	StringToken

	// refs: https://www.w3.org/TR/css-syntax-3/#typedef-at-keyword-token
	AtKeyword
)
