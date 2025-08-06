package boxes

import (
	"github.com/schidstorm/s3tool/pkg/s3lib"
)

type ContextPage struct {
	*ListPage
}

func NewContextPage() *ContextPage {
	page := &ContextPage{
		ListPage: NewListPage(),
	}

	page.ListPage.SetSelectedFunc(func(columns []string) {
		typeStr := columns[0]
		name := columns[1]

		client, err := s3lib.LoadContext(typeStr, name)
		if err != nil {
			activeApp.SetError(err)
			return
		}

		activeApp.OpenPage(NewBucketsBox(client))
	})

	page.load()

	return page
}

func (b *ContextPage) Title() string {
	return "Contexts"
}

func (b *ContextPage) Hotkeys() map[string]string {
	return map[string]string{}
}

func (b *ContextPage) load() {
	b.ListPage.ClearRows()
	b.ListPage.AddRow(Row{
		Header:  true,
		Columns: []string{"Type", "Name"},
	})

	contexts := s3lib.ListContexts()
	for _, context := range contexts {
		b.ListPage.AddRow(Row{
			Columns: []string{
				string(context.Type),
				context.Name,
			},
		})
	}
}
