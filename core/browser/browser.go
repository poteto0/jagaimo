package browser

import "weak"

type IBrowser interface {
	CurrentPage() *Page
}

// ! single tab browser for now
type Browser struct {
	// tab for watching: currently always 0
	currentPageIndex uint8

	// all active pages: currently lens is always 1
	pages []*Page
}

func NewBrowser() IBrowser {
	page := NewPage().(*Page)
	browser := &Browser{
		currentPageIndex: 0,
		pages:            []*Page{},
	}

	page.Browser = weak.Make(browser)
	browser.pages = append(browser.pages, page)

	return browser
}

func (b *Browser) CurrentPage() *Page {
	return b.pages[b.currentPageIndex]
}
