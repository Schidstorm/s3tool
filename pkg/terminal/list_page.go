package terminal

import (
	"strings"

	"github.com/rivo/tview"
)

type Row struct {
	Header  bool
	Columns []string
}

type ListPage[TItem any] struct {
	*tview.Flex

	tviewTable *tview.Table
	table      *Table[TItem]
}

func NewListPage[TItem any]() *ListPage[TItem] {
	table := tview.NewTable().SetSelectable(true, false)
	table.SetFixed(1, 0)

	flex := tview.NewFlex()
	flex.SetDirection(tview.FlexRow)
	flex.AddItem(table, 0, 1, true)

	box := &ListPage[TItem]{
		tviewTable: table,
		Flex:       flex,
		table:      NewTable[TItem](),
	}

	box.update()
	return box
}

func (b *ListPage[TItem]) SetSelectedFunc(f func(item TItem)) {
	if f == nil {
		b.tviewTable.SetSelectedFunc(nil)
		return
	}

	b.tviewTable.SetSelectedFunc(func(row, column int) {
		item, ok := b.table.GetRowItem(row - 1)
		if !ok {
			return
		}
		f(item)
	})
}

func (b *ListPage[TItem]) SetSearch(search string) {
	b.table.SetFilter(search)
	b.update()
}

func (b *ListPage[TItem]) update() {
	b.tviewTable.Clear()

	columns := b.table.Columns()
	for colIndex, col := range columns {
		cell := tview.NewTableCell(col)
		cell.SetAlign(tview.AlignLeft)
		cell.SetExpansion(1)
		cell.SetStyle(DefaultTheme.TableHeader)
		cell.SetSelectable(false)
		b.tviewTable.SetCell(0, colIndex, cell)
	}

	for rowIndex, row := range b.table.Rows() {
		for columnIndex, item := range row {
			cell := tview.NewTableCell(item)
			cell.SetAlign(tview.AlignLeft)
			cell.SetExpansion(1)
			cell.SetStyle(DefaultTheme.TableCell)
			cell.SetSelectable(true)
			cell.SelectedStyle = DefaultTheme.TableSelected
			b.tviewTable.SetCell(rowIndex+1, columnIndex, cell)
		}
	}
}

func matchAnyItems(term string, items []string) bool {
	if term == "" {
		return true
	}
	for _, item := range items {
		if containsIgnoreCase(item, term) {
			return true
		}
	}
	return false
}

func containsIgnoreCase(s, substr string) bool {
	if len(s) < len(substr) {
		return false
	}
	return strings.Contains(strings.ToLower(s), strings.ToLower(substr))
}

func (b *ListPage[TItem]) Add(row TItem) {
	b.table.Add(row)
	b.update()
}

func (b *ListPage[TItem]) AddAll(rows []TItem) {
	for _, row := range rows {
		b.table.Add(row)
	}
	b.update()
}

func (b *ListPage[TItem]) ClearRows() {
	b.table.Clear()
	b.update()
}

func (b *ListPage[TItem]) GetSelectedRow() (TItem, bool) {
	row, _ := b.tviewTable.GetSelection()
	return b.table.GetRowItem(row - 1)
}

func (b *ListPage[TItem]) AddColumn(name string, filler func(item TItem) string) {
	b.table.AddColumn(name, filler)
	b.update()
}
