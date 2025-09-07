package terminal

import "github.com/gdamore/tcell/v2"

type Theme struct {
	// reworked colors
	PrimaryColor   tcell.Color
	SecondaryColor tcell.Color
	KeyColor       tcell.Color
	LabelColor     tcell.Color
	ErrorColor     tcell.Color
	InfoColor      tcell.Color
	BorderColor    tcell.Color
	HighlightStyle tcell.Style
}

var DefaultTheme = Theme{
	PrimaryColor:   tcell.ColorWhite,
	SecondaryColor: tcell.ColorGray,
	KeyColor:       tcell.ColorNavy,
	LabelColor:     tcell.ColorOrange,
	ErrorColor:     tcell.ColorRed,
	InfoColor:      tcell.ColorGreen,
	BorderColor:    tcell.ColorNavy,
	HighlightStyle: tcell.StyleDefault.Foreground(tcell.ColorBlack).Background(tcell.ColorLightBlue),
}

var DefaultStyle = tcell.StyleDefault.Foreground(DefaultTheme.PrimaryColor).Background(tcell.ColorBlack)
