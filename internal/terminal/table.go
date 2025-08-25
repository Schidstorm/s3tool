package terminal

type ColumnFiller[TItem any] func(item TItem) string

type column[TItem any] struct {
	name   string
	filler ColumnFiller[TItem]
}

type Table[TItem any] struct {
	columns      []column[TItem]
	allRows      [][]string
	allItems     []TItem
	filteredRows []int
	filter       string
}

func NewTable[TItem any]() *Table[TItem] {
	return &Table[TItem]{}
}

func (t *Table[TItem]) AddColumn(name string, filler ColumnFiller[TItem]) {
	t.columns = append(t.columns, column[TItem]{
		name:   name,
		filler: filler,
	})
}

func (t *Table[TItem]) Add(item TItem) {
	t.allItems = append(t.allItems, item)

	var row []string
	for _, col := range t.columns {
		row = append(row, col.filler(item))
	}

	t.allRows = append(t.allRows, row)
	if matchAnyItems(t.filter, row) {
		t.filteredRows = append(t.filteredRows, len(t.allRows)-1)
	}
}

func (t *Table[TItem]) Rows() [][]string {
	var rows [][]string
	for _, rowIndex := range t.filteredRows {
		rows = append(rows, t.allRows[rowIndex])
	}
	return rows
}

func (t *Table[TItem]) Columns() []string {
	var cols []string
	for _, col := range t.columns {
		cols = append(cols, col.name)
	}
	return cols
}

func (t *Table[TItem]) GetRowItem(rowIndex int) (TItem, bool) {
	var zero TItem
	if rowIndex < 0 || rowIndex >= len(t.filteredRows) {
		return zero, false
	}

	rowIndex = t.filteredRows[rowIndex]

	return t.allItems[rowIndex], true
}

func (t *Table[TItem]) SetFilter(filter string) {
	t.filter = filter

	t.filteredRows = t.filteredRows[:0]
	for i, row := range t.allRows {
		if matchAnyItems(t.filter, row) {
			t.filteredRows = append(t.filteredRows, i)
		}
	}
}

func (t *Table[TItem]) Clear() {
	t.allRows = t.allRows[:0]
	t.allItems = t.allItems[:0]
	t.filteredRows = t.filteredRows[:0]
}
