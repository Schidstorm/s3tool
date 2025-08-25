package terminal

import (
	"github.com/rivo/tview"
	"github.com/schidstorm/s3tool/internal/s3lib"
)

// var activeApp *App

type App struct {
	*tview.Application
	root *RootPage
}

func NewApp(page PageContent, loaders ...s3lib.ConnectorLoader) *App {
	root := NewRootPage()
	app := &App{
		root:        root,
		Application: tview.NewApplication(),
	}

	// activeApp = app
	if page == nil {
		page = NewProfilePage(app.CreateContext(), loaders)
	}

	root.OpenPage(page)

	return app
}

func (a *App) CreateContext() Context {
	return NewContext().
		WithOpenPageFunc(a.OpenPage).
		WithErrorFunc(a.SetError).
		WithModalFunc(a.Modal).
		WithSuspendAppFunc(a.Application.Suspend)
}

func (a *App) Run() error {
	return a.Application.SetRoot(a.root, true).Run()
}

func (a *App) Modal(p ModalBuilder, name string, width, height int) {
	a.root.Modal(p, name, width, height)
}

func (a *App) OpenPage(page PageContent) {
	a.root.OpenPage(page)
}

func (a *App) SetError(err error) {
	a.root.SetError(err)
}
