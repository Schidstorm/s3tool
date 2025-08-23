package terminal

import "github.com/gdamore/tcell/v2"

type Theme struct {
	Hotkey             tcell.Style
	HotkeyLabel        tcell.Style
	TableHeader        tcell.Style
	TableSelected      tcell.Style
	TableCell          tcell.Style
	PageBorder         tcell.Style
	ProfileKey         tcell.Style
	ProfileValue       tcell.Style
	PageTitlePrimary   tcell.Style
	PageTitleSecondary tcell.Style
}

var DefaultTheme = Theme{
	Hotkey:             tcell.StyleDefault.Foreground(tcell.ColorNavy),
	HotkeyLabel:        tcell.StyleDefault.Foreground(tcell.ColorGray),
	TableHeader:        tcell.StyleDefault.Foreground(tcell.ColorWhite),
	TableSelected:      tcell.StyleDefault.Foreground(tcell.ColorBlack).Background(tcell.ColorLightBlue),
	TableCell:          tcell.StyleDefault.Foreground(tcell.ColorWhite),
	PageBorder:         tcell.StyleDefault.Foreground(tcell.ColorNavy),
	ProfileKey:         tcell.StyleDefault.Foreground(tcell.ColorOrange),
	ProfileValue:       tcell.StyleDefault.Foreground(tcell.ColorWhite).Bold(true),
	PageTitlePrimary:   tcell.StyleDefault.Foreground(tcell.ColorLightBlue),
	PageTitleSecondary: tcell.StyleDefault.Foreground(tcell.ColorGray),
}
