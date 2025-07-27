package boxes

import (
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Row struct {
	Header  bool
	Columns []string
}

type ListPage struct {
	*tview.Flex

	table      *tview.Table
	searchTerm string

	rows []Row
}

func NewListPage() *ListPage {
	table := tview.NewTable().SetSelectable(true, false)

	flex := tview.NewFlex()
	flex.SetDirection(tview.FlexRow)
	flex.AddItem(table, 0, 1, true)

	box := &ListPage{
		table: table,
		Flex:  flex,
	}

	box.update()
	return box
}

func (b *ListPage) SetSelectedFunc(f func(columns []string)) {
	if f == nil {
		b.table.SetSelectedFunc(nil)
		return
	}

	b.table.SetSelectedFunc(func(row, column int) {
		var columns []string
		if row >= 0 && row < len(b.rows) {
			for column := 0; column < b.table.GetColumnCount(); column++ {
				columns = append(columns, b.table.GetCell(row, column).Text)
			}
		}
		f(columns)
	})
}

func (b *ListPage) SetSearch(search string) {
	b.searchTerm = search
	b.update()
}

func (b *ListPage) update() {
	b.table.Clear()

	rowIndex := 0
	for _, row := range b.rows {
		if !row.Header && !matchAnyItems(b.searchTerm, row.Columns) {
			continue
		}

		for columnIndex, item := range row.Columns {
			bucketNameCell := tview.NewTableCell(item)
			if row.Header {
				bucketNameCell.SetTextColor(tcell.ColorYellow)
				bucketNameCell.SetAlign(tview.AlignCenter)
			} else {
				bucketNameCell.SetTextColor(tcell.ColorWhite)
			}
			bucketNameCell.SetAlign(tview.AlignLeft)
			bucketNameCell.SetSelectable(!row.Header)
			bucketNameCell.SetExpansion(1)
			b.table.SetCell(rowIndex, columnIndex, bucketNameCell)
		}

		rowIndex++
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

func (b *ListPage) AddRow(row Row) {
	b.rows = append(b.rows, row)
	b.update()
}

func (b *ListPage) AddRows(rows []Row) {
	b.rows = append(b.rows, rows...)
	b.update()
}

func (b *ListPage) AddRowAt(index int, row Row) {
	if index < 0 || index > len(b.rows) {
		return
	}
	b.rows = append(b.rows[:index], append([]Row{row}, b.rows[index:]...)...)
	b.update()
}

func (b *ListPage) AddRowsAt(index int, rows []Row) {
	if index < 0 || index > len(b.rows) {
		return
	}
	b.rows = append(b.rows[:index], append(rows, b.rows[index:]...)...)
	b.update()
}

func (b *ListPage) RemoveRow(index int) {
	if index < 0 || index >= len(b.rows) {
		return
	}
	b.rows = append(b.rows[:index], b.rows[index+1:]...)
	b.update()
}

func (b *ListPage) ClearRows() {
	b.rows = nil
	b.update()
}

func (b *ListPage) GetSelectedRow() (int, Row) {
	row, _ := b.table.GetSelection()
	if row < 0 || row >= len(b.rows) {
		return -1, Row{}
	}
	return row, b.rows[row]
}
