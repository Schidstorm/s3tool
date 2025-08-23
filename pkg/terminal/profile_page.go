package terminal

import (
	"context"
	"errors"

	"github.com/gdamore/tcell/v2"
	"github.com/schidstorm/s3tool/pkg/s3lib"
)

type ProfilePage struct {
	*ListPage

	loaders []s3lib.ConnectorLoader
}

func NewProfilePage(loaders []s3lib.ConnectorLoader) *ProfilePage {
	page := &ProfilePage{
		ListPage: NewListPage(),
		loaders:  loaders,
	}

	page.ListPage.SetSelectedFunc(func(columns []string) {
		if len(columns) < 2 {
			return
		}

		typeStr := columns[0]
		name := columns[1]

		profiles, err := loadConnectors(loaders)
		if err != nil {
			activeApp.SetError(err)
			return
		}

		var profile s3lib.Connector
		for _, p := range profiles {
			if p.Type() == typeStr && p.Name() == name {
				profile = p
				break
			}
		}
		if profile == nil {
			activeApp.SetError(errors.New("profile not found"))
			return
		}

		client, err := profile.CreateClient(context.Background())
		if err != nil {
			activeApp.SetError(err)
			return
		}

		activeApp.SetS3Client(client, "")
		activeApp.OpenPage(AttachClose{
			PageContent: NewBucketsPage(client),
			Closer: CloseFunc(func() {
				activeApp.SetS3Client(nil, "")
			}),
		})
	})

	page.load()

	return page
}

func loadConnectors(loaders []s3lib.ConnectorLoader) ([]s3lib.Connector, error) {
	var profiles []s3lib.Connector

	for _, loader := range loaders {
		loadedProfiles, err := loader.Load()
		if err != nil {
			return nil, err
		}
		profiles = append(profiles, loadedProfiles...)
	}

	return profiles, nil
}

func (b *ProfilePage) Title() string {
	return "Profiles"
}

func (b *ProfilePage) Hotkeys() map[tcell.EventKey]Hotkey {
	return map[tcell.EventKey]Hotkey{}
}

func (b *ProfilePage) load() {
	b.ListPage.ClearRows()
	b.ListPage.AddRow(Row{
		Header:  true,
		Columns: []string{"Type", "Name"},
	})

	profiles, err := loadConnectors(b.loaders)
	if err != nil {
		activeApp.SetError(err)
		return
	}

	for _, p := range profiles {
		b.ListPage.AddRow(Row{
			Columns: []string{p.Type(), p.Name()},
		})
	}

}
