package terminal

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/schidstorm/s3tool/pkg/s3lib"
)

type RootPage struct {
	*tview.Flex
	pages       *tview.Pages
	profileInfo *ProfileInfoBox
	hotkeyInfo  *HotkeyInfoBox

	pageStask  []*Page
	modalOpen  string
	statusText *tview.TextView
}

func NewRootPage() *RootPage {
	header := tview.NewFlex()
	header.SetDirection(tview.FlexColumn)

	profileInfo := NewProfileInfoBox()
	header.AddItem(profileInfo, 0, 1, false)

	hotkeyInfo := NewHotkeyInfoBox()
	hotkeyInfo.Update(nil)
	header.AddItem(hotkeyInfo, 40, 0, false)

	content := tview.NewPages()
	content.SetBorder(true)
	content.SetBorderStyle(DefaultTheme.PageBorder)
	content.SetBorderPadding(0, 1, 1, 1)

	statusText := tview.NewTextView().SetTextAlign(tview.AlignCenter)

	flex := tview.NewFlex()
	flex.SetDirection(tview.FlexRow)
	flex.AddItem(header, 8, 1, false)
	flex.AddItem(content, 0, 1, true)
	flex.AddItem(statusText, 1, 1, false)

	a := &RootPage{
		profileInfo: profileInfo,
		hotkeyInfo:  hotkeyInfo,
		pages:       content,
		Flex:        flex,
		statusText:  statusText,
	}

	flex.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEscape:
			if a.modalOpen != "" {
				a.CloseModal(a.modalOpen)
				return nil
			} else if a.statusText.GetText(true) != "" {
				a.statusText.SetText("")
				return nil
			}
		}
		return event
	})

	return a
}

func (a *RootPage) Modal(p ModalBuilder, name string, width, height int) {
	content := p(func() {
		a.CloseModal(name)
	})
	modal := tview.NewFlex().
		AddItem(nil, 0, 1, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(nil, 0, 1, false).
			AddItem(content, height, 1, true).
			AddItem(nil, 0, 1, false), width, 1, true).
		AddItem(nil, 0, 1, false)

	a.pages.AddPage(name, modal, true, true)
	a.modalOpen = name
}

func (a *RootPage) CloseModal(name string) {
	if a.pages.HasPage(name) {
		a.pages.RemovePage(name)
		a.modalOpen = ""
	}
}

func (a *RootPage) OpenPage(pageContent PageContent) {
	if pageContent == nil {
		return
	}

	page := NewPage(pageContent)
	page.SetCloseHandler(func() {
		a.closePage()
	})
	a.pageStask = append(a.pageStask, page)
	a.openPage(page)
}

func (a *RootPage) openPage(page *Page) {
	a.pages.AddPage(page.Title(), page, true, true)
	a.pages.SetTitle(" " + page.Title() + " ")
	titleFg, _, _ := DefaultTheme.PageTitlePrimary.Decompose()
	a.pages.SetTitleColor(titleFg)
	a.pages.SwitchToPage(page.Title())
	a.hotkeyInfo.Update(page.content)
}

func (a *RootPage) closePage() {
	if len(a.pageStask) == 1 {
		return
	}

	a.pageStask = a.pageStask[:len(a.pageStask)-1]

	a.openPage(a.pageStask[len(a.pageStask)-1])
}

func (a *RootPage) SetError(err error) {
	errorText := s3lib.ErrorText(err)
	a.statusText.SetText("Error: " + errorText)
	a.statusText.SetTextColor(tcell.ColorRed)
}

func (a *RootPage) UpdateContext(c Context) {
	a.profileInfo.UpdateContext(c)
}
