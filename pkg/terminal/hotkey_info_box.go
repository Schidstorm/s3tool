package terminal

import "github.com/rivo/tview"

type HotkeyInfoBox struct {
	*tview.Table
}

func NewHotkeyInfoBox() *HotkeyInfoBox {
	table := tview.NewTable().SetBorders(false)
	table.SetSelectable(false, false)
	table.Box.SetBorder(true)
	table.Box.SetTitle("Hotkeys")
	table.Box.SetTitleAlign(tview.AlignLeft)

	info := &HotkeyInfoBox{
		Table: table,
	}

	return info
}

func (info *HotkeyInfoBox) Update(pageContent PageContent) {
	info.Clear()

	if pageContent == nil {
		return
	}

	hotkeys := pageContent.Hotkeys()
	row := 0
	for key, title := range hotkeys {
		titleCell := tview.NewTableCell(title)
		titleCell.SetTextColor(tview.Styles.PrimaryTextColor)
		titleCell.SetExpansion(5)

		keyCell := tview.NewTableCell(key)
		keyCell.SetTextColor(tview.Styles.PrimaryTextColor)
		keyCell.SetExpansion(1)

		info.SetCell(row, 0, keyCell)
		info.SetCell(row, 1, titleCell)
		row++
	}
}
