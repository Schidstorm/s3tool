package boxes

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/schidstorm/s3tool/pkg/s3lib"
)

type ContextPage struct {
	*tview.Flex

	list       *ListPage
	statusText *tview.TextView
}

func NewContextPage() *ContextPage {
	listPage := NewListPage()
	statusText := tview.NewTextView().SetTextAlign(tview.AlignCenter)

	flex := tview.NewFlex()
	flex.SetDirection(tview.FlexRow)
	flex.AddItem(listPage, 0, 1, true)
	flex.AddItem(statusText, 1, 1, false)

	page := &ContextPage{
		Flex:       flex,
		list:       listPage,
		statusText: statusText,
	}

	listPage.SetSelectedFunc(func(columns []string) {
		typeStr := columns[0]
		name := columns[1]

		client, err := s3lib.LoadContext(typeStr, name)
		if err != nil {
			page.setError(err)
			return
		}

		activeApp.OpenPage(NewBucketsBox(client))
	})

	page.load()

	return page
}

func (b *ContextPage) Title() string {
	return "Contexts"
}

func (b *ContextPage) Hotkeys() map[string]string {
	return map[string]string{}
}

func (b *ContextPage) SetSearch(search string) {
	b.list.SetSearch(search)
}

func (b *ContextPage) load() {
	b.list.ClearRows()
	b.list.AddRow(Row{
		Header:  true,
		Columns: []string{"Type", "Name"},
	})

	contexts := s3lib.ListContexts()
	for _, context := range contexts {
		b.list.AddRow(Row{
			Columns: []string{
				string(context.Type),
				context.Name,
			},
		})
	}
}

func (b *ContextPage) setError(err error) {
	b.statusText.SetText("Error: " + err.Error())
	b.statusText.SetTextColor(tcell.ColorRed)
}
