package terminal

import (
	"context"
	"errors"

	"github.com/schidstorm/s3tool/pkg/profile"
)

type ProfilePage struct {
	*ListPage
}

func NewProfilePage() *ProfilePage {
	page := &ProfilePage{
		ListPage: NewListPage(),
	}

	page.ListPage.SetSelectedFunc(func(columns []string) {
		typeStr := columns[0]
		name := columns[1]

		profiles := profile.List()
		var profile profile.Connector
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

		activeApp.OpenPage(NewBucketsBox(client))
	})

	page.load()

	return page
}

func (b *ProfilePage) Title() string {
	return "Profiles"
}

func (b *ProfilePage) Hotkeys() map[string]string {
	return map[string]string{}
}

func (b *ProfilePage) load() {
	b.ListPage.ClearRows()
	b.ListPage.AddRow(Row{
		Header:  true,
		Columns: []string{"Type", "Name"},
	})

	profiles := profile.List()
	for _, p := range profiles {
		b.ListPage.AddRow(Row{
			Columns: []string{p.Type(), p.Name()},
		})
	}

}
