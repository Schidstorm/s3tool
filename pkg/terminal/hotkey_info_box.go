package terminal

import (
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type HotkeyInfoBox struct {
	*tview.Table
}

func NewHotkeyInfoBox() *HotkeyInfoBox {
	table := tview.NewTable().SetBorders(false)
	table.SetSelectable(false, false)
	table.Box.SetBorder(false)
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
	for key, hk := range hotkeys {
		keyCell := tview.NewTableCell(eventKeyToString(key))
		keyCell.SetStyle(DefaultTheme.Hotkey)
		keyCell.SetExpansion(1)

		titleCell := tview.NewTableCell(hk.Title)
		titleCell.SetStyle(DefaultTheme.HotkeyLabel)
		titleCell.SetExpansion(5)

		info.SetCell(row, 0, keyCell)
		info.SetCell(row, 1, titleCell)
		row++
	}
}

func eventKeyToString(key tcell.EventKey) string {
	var keyNameParts []string
	if key.Key() == tcell.KeyRune {
		keyNameParts = []string{string(key.Rune())}
	} else {
		keyNameParts = []string{strings.ToLower(key.Name())}
	}

	if key.Modifiers()&tcell.ModCtrl != 0 {
		keyNameParts = append([]string{"ctrl"}, keyNameParts...)
	}
	if key.Modifiers()&tcell.ModAlt != 0 {
		keyNameParts = append([]string{"alt"}, keyNameParts...)
	}
	if key.Modifiers()&tcell.ModShift != 0 {
		keyNameParts = append([]string{"shift"}, keyNameParts...)
	}
	return "<" + strings.Join(keyNameParts, "+") + ">"
}
