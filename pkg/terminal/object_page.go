package terminal

import (
	"context"
	"path"

	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type ObjectPage struct {
	*tview.Flex
	table *tview.Table

	context    Context
	searchTerm string
}

func NewObjectPage(context Context) *ObjectPage {
	table := tview.NewTable().SetSelectable(true, false)
	table.SetFixed(1, 0)

	flex := tview.NewFlex()
	flex.SetDirection(tview.FlexRow)
	flex.AddItem(table, 0, 1, true)

	page := &ObjectPage{
		table:   table,
		Flex:    flex,
		context: context,
	}

	page.load()

	table.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyDelete:
		case tcell.KeyRune:
			if event.Rune() == 'v' {
				err := viewObject(context)
				if err != nil {
					context.SetError(err)
				}
				return nil
			}
			if event.Rune() == 'e' {
				err := editObject(context)
				if err != nil {
					context.SetError(err)
				}
				return nil
			}
		}
		return event
	})

	page.load()

	return page
}

func (b *ObjectPage) SetSearch(search string) {
	b.searchTerm = search
	b.load()
}

func (b *ObjectPage) Context() Context {
	return b.context
}

func (b *ObjectPage) Title() string {
	title := "Objects in " + b.context.Bucket()
	if b.context.ObjectKey() != "" {
		title += " - " + b.context.ObjectKey()
	}

	return title
}

func (b *ObjectPage) Hotkeys() map[tcell.EventKey]Hotkey {
	return map[tcell.EventKey]Hotkey{
		EventKey(tcell.KeyRune, 'v', 0): {
			Title: "View Object",
			Handler: func(event *tcell.EventKey) *tcell.EventKey {
				err := viewObject(b.context)
				if err != nil {
					b.context.SetError(err)
				}
				return nil
			},
		},
		EventKey(tcell.KeyRune, 'e', 0): {
			Title: "Edit Object",
			Handler: func(event *tcell.EventKey) *tcell.EventKey {
				err := editObject(b.context)
				if err != nil {
					b.context.SetError(err)
				}
				return nil
			},
		},
	}
}

type item struct {
	title string
	value []string
}

func (b *ObjectPage) load() {
	b.table.Clear()

	obj, err := b.context.S3Client().GetObject(context.Background(), b.context.Bucket(), b.context.ObjectKey())
	if err != nil {
		b.context.SetError(err)
		return
	}

	var items []item
	addItem := func(title string, value *string) {
		if value != nil {
			items = append(items, item{title: title, value: []string{*value}})
		}
	}
	addItem("Bucket", &obj.Bucket)
	name := path.Base(obj.Key)
	addItem("Name", &name)
	addItem("Region", &obj.Region)
	addItem("Owner", obj.Owner)
	addItem("Type", obj.Type)
	size := humanizeSize(obj.Size)
	addItem("Size", &size)
	addItem("ETag", obj.ETag)
	addItem("LegalHold", &obj.LegalHold)
	lastModified := humanizeTime(obj.LastModified)
	addItem("LastModified", &lastModified)
	for k, v := range obj.Metadata {
		addItem(k, &v)
	}
	if len(obj.Tags) > 0 {
		var tags []string
		for k, v := range obj.Tags {
			tags = append(tags, k+"="+v)
		}
		items = append(items, item{title: "Tags", value: tags})
	}

	for row, it := range items {
		if !matchAnyItems(b.searchTerm, append([]string{it.title}, it.value...)) {
			continue
		}

		b.table.SetCell(row, 0, tview.NewTableCell(it.title).
			SetTextColor(tcell.ColorYellow).
			SetAlign(tview.AlignLeft).
			SetSelectable(false))
		b.table.SetCell(row, 1, tview.NewTableCell(strings.Join(it.value, ", ")).
			SetTextColor(tcell.ColorWhite).
			SetAlign(tview.AlignLeft).
			SetSelectable(false))
	}
}
