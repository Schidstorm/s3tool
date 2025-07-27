package boxes

import "github.com/rivo/tview"

var activeApp *App

type App struct {
	*tview.Application
	root *RootPage
}

func NewApp() *App {
	root := NewRootPage()
	app := &App{
		root:        root,
		Application: tview.NewApplication(),
	}

	activeApp = app

	return app
}

func (a *App) Run() error {
	return a.Application.SetRoot(a.root, true).Run()
}

func (a *App) Modal(p tview.Primitive, name string, width, height int) {
	a.root.Modal(p, name, width, height)
}

func (a *App) CloseModal(name string) {
	a.root.CloseModal(name)
}

func (a *App) OpenPage(page PageContent) {
	a.root.OpenPage(page)
}

func (a *App) closePage() {
	a.root.closePage()
}
