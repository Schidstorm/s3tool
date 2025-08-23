package terminal

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/schidstorm/s3tool/pkg/s3lib"
)

type SimulatedApp struct {
	*App
	screen tcell.SimulationScreen
}

func NewSimulatedApp(page PageContent, loaders ...s3lib.ConnectorLoader) SimulatedApp {
	root := NewRootPage()
	app := &App{
		root:        root,
		Application: tview.NewApplication(),
	}

	screen := tcell.NewSimulationScreen("")
	app.Application.SetScreen(screen)

	activeApp = app
	if page == nil {
		page = NewProfilePage(loaders)
	}

	root.OpenPage(page)

	return SimulatedApp{
		App:    app,
		screen: screen,
	}
}

func (a SimulatedApp) GetScreen() tcell.SimulationScreen {
	return a.screen
}
