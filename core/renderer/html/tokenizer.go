package html

type IHtmlTokenizer interface{}

type HtmlTokenizer struct {
	State       State
	Pos         uint
	ReConsume   bool
	LatestToken *HtmlTokenizer
	Input       []byte
	Buf         string
}

func NewHtmlTokenizer(html string) IHtmlTokenizer {
	return &HtmlTokenizer{
		State:       Data,
		Pos:         0,
		ReConsume:   false,
		LatestToken: nil,
		Input:       []byte(html),
		Buf:         "",
	}
}
