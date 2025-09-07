package terminal

import "github.com/gdamore/tcell/v2"

type Theme struct {
	HotkeyKey          tcell.Style
	HotkeyTitle        tcell.Style
	TableHeader        tcell.Style
	TableSelected      tcell.Style
	TableCell          tcell.Style
	PageBorder         tcell.Style
	ModalBorder        tcell.Style
	ProfileKey         tcell.Style
	ProfileValue       tcell.Style
	PageTitlePrimary   tcell.Color
	ModalTitlePrimary  tcell.Color
	PageTitleSecondary tcell.Color
	ErrorMessage       tcell.Style
}

var DefaultTheme = Theme{
	HotkeyKey:          tcell.StyleDefault.Foreground(tcell.ColorNavy),
	HotkeyTitle:        tcell.StyleDefault.Foreground(tcell.ColorGray),
	TableHeader:        tcell.StyleDefault.Foreground(tcell.ColorWhite).Bold(true),
	TableSelected:      tcell.StyleDefault.Foreground(tcell.ColorBlack).Background(tcell.ColorLightBlue),
	TableCell:          tcell.StyleDefault.Foreground(tcell.ColorWhite),
	PageBorder:         tcell.StyleDefault.Foreground(tcell.ColorNavy),
	ModalBorder:        tcell.StyleDefault.Foreground(tcell.ColorNavy),
	ProfileKey:         tcell.StyleDefault.Foreground(tcell.ColorOrange),
	ProfileValue:       tcell.StyleDefault.Foreground(tcell.ColorWhite).Bold(true),
	PageTitlePrimary:   tcell.ColorLightBlue,
	ModalTitlePrimary:  tcell.ColorLightBlue,
	PageTitleSecondary: tcell.ColorGray,
	ErrorMessage:       tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorRed),
}
