package browser

import (
	"weak"

	"github.com/poteto0/jagaimo/core/renderer/dom"
)

// One Tab
type Page struct {
	Browser weak.Pointer[Browser]
	frame   *dom.Window
}

func NewPage() *Page {
	return &Page{
		Browser: weak.Pointer[Browser]{},
		frame:   nil,
	}
}
