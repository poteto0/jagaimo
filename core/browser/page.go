package browser

import (
	"weak"

	"github.com/poteto0/jagaimo/core"
	"github.com/poteto0/jagaimo/core/renderer/dom"
	"github.com/poteto0/jagaimo/core/renderer/html"
	"github.com/poteto0/jagaimo/utils"
)

type IPage interface {
	// DOM tree to string for debug
	ReceiveResponse(response core.HttpResponse) string
}

// One Tab
type Page struct {
	Browser weak.Pointer[Browser]
	frame   *dom.Window
}

func NewPage() IPage {
	return &Page{
		Browser: weak.Pointer[Browser]{},
		frame:   nil,
	}
}

func (p *Page) ReceiveResponse(response core.HttpResponse) string {
	frame := p.createFrame(response.Body)
	p.frame = frame

	if p.frame != nil {
		document := p.frame.Document()
		debug := utils.ConvertDomToString(document)
		return debug
	}

	return ""
}

func (p *Page) createFrame(rawHtml string) *dom.Window {
	htmlTokenizer := html.NewHtmlTokenizer(rawHtml)
	return html.NewHtmlParser(htmlTokenizer).ConstructTree()
}
