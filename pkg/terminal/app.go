package terminal

import (
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/rivo/tview"
)

var activeApp *App

type App struct {
	*tview.Application
	root *RootPage
}

func NewApp(page PageContent) *App {
	root := NewRootPage()
	app := &App{
		root:        root,
		Application: tview.NewApplication(),
	}

	activeApp = app
	if page == nil {
		page = NewProfilePage()
	}

	root.OpenPage(page)

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

func (a *App) SetError(err error) {
	a.root.SetError(err)
}

func (a *App) SetS3Client(client *s3.Client, bucketName string) {
	a.root.SetS3Client(client, bucketName)
}
