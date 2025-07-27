package boxes

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type PageContent interface {
	tview.Primitive

	Title() string
	Hotkeys() map[string]string
	SetSearch(term string)
}

type Page struct {
	*tview.Flex

	searchFlex     *tview.Flex
	content        PageContent
	isSearchActive bool
	searchTerm     string
	closeHandler   func()
}

func NewPage(content PageContent) *Page {
	contentFlex := tview.NewFlex()
	contentFlex.SetDirection(tview.FlexRow)
	contentFlex.AddItem(content, 0, 1, true)

	searchFlex := tview.NewFlex()
	searchFlex.SetDirection(tview.FlexRow)
	searchFlex.AddItem(contentFlex, 0, 1, true)

	p := &Page{
		Flex:       searchFlex,
		content:    content,
		searchFlex: searchFlex,
	}

	contentFlex.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEscape:
			if p.closeHandler != nil {
				p.closeHandler()
				return nil
			}
		case tcell.KeyRune:
			if event.Rune() == '/' {
				p.activateSearch()
				return nil
			}
		}
		return event
	})

	searchFlex.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if p.isSearchActive {
			switch event.Key() {
			case tcell.KeyEnter:
				p.deactivateSearch()
				return nil
			case tcell.KeyEscape:
				p.content.SetSearch("")
				p.searchTerm = ""
				p.deactivateSearch()
				return nil
			}
		} else if p.searchTerm != "" {
			switch event.Key() {
			case tcell.KeyEscape:
				p.content.SetSearch("")
				p.searchTerm = ""
				return nil
			}
		}

		return event
	})

	return p
}

// Modal flex -> Search Flex -> Content

// modes: modal mode, search mode, search term not empty mode
// when in no  mode then normal hotkeys
// modal and search and search term not empty not together
// modal input capture + search input capture + search term not empty above normal hotkeys

// escape to close page is part of normal hotkeys

func (p *Page) activateSearch() {
	if p.isSearchActive {
		return
	}

	p.isSearchActive = true
	search := tview.NewInputField()
	search.SetChangedFunc(func(text string) {
		p.content.SetSearch(text)
		p.searchTerm = text
	})
	search.SetPlaceholder("Search: ")

	contentFlex := p.searchFlex.GetItem(0)
	p.searchFlex.Clear()
	p.searchFlex.AddItem(search, 1, 1, true)
	p.searchFlex.AddItem(contentFlex, 0, 1, false)

	search.Focus(nil)
}

func (p *Page) deactivateSearch() {
	if !p.isSearchActive {
		return
	}

	p.isSearchActive = false
	contentFlex := p.searchFlex.GetItem(1)
	p.Clear()
	p.searchFlex.AddItem(contentFlex, 0, 1, true)
}

func (p *Page) Title() string {
	return p.content.Title()
}

func (p *Page) SetCloseHandler(handler func()) {
	p.closeHandler = handler
}
