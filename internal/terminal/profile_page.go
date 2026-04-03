package terminal

import (
	"context"
	"slices"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/schidstorm/s3tool/internal/s3lib"
)

type ProfilePage struct {
	*ListPage[s3lib.Connector]

	loaders []s3lib.ConnectorLoader
	context Context
}

func NewProfilePage(c Context, loaders []s3lib.ConnectorLoader) *ProfilePage {
	page := &ProfilePage{
		ListPage: NewListPage[s3lib.Connector](),
		loaders:  loaders,
		context:  c,
	}

	page.SetMultiSelect(false)
	page.AddColumn("Type", func(item s3lib.Connector) string { return item.Type() })
	page.AddColumn("Name", func(item s3lib.Connector) string { return item.Name() })

	page.SetSelectedFunc(func(connector s3lib.Connector) {
		client, err := connector.CreateClient(context.Background())
		if err != nil {
			c.SetError(err)
			return
		}

		c.OpenPage(NewBucketsPage(c.WithClient(client)))
	})

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

	slices.SortFunc(profiles, func(a, b s3lib.Connector) int {
		if byType := strings.Compare(a.Type(), b.Type()); byType != 0 {
			return byType
		}
		return strings.Compare(a.Name(), b.Name())
	})

	return profiles, nil
}

func (b *ProfilePage) Title() string {
	return "Profiles"
}

func (b *ProfilePage) Context() Context {
	return b.context
}

func (b *ProfilePage) Hotkeys() map[tcell.EventKey]Hotkey {
	return map[tcell.EventKey]Hotkey{}
}

func (b *ProfilePage) Load() error {
	b.ClearRows()

	profiles, err := loadConnectors(b.loaders)
	if err != nil {
		return err
	}

	b.AddAll(profiles)
	return nil
}
