package terminal

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testTableItem struct {
	Name string
	Age  int
	City string
}

func testTable() *Table[testTableItem] {
	table := NewTable[testTableItem]()
	table.AddColumn("Name", func(i testTableItem) string { return i.Name })
	table.AddColumn("Age", func(i testTableItem) string { return strconv.Itoa(i.Age) })
	table.AddColumn("City", func(i testTableItem) string { return i.City })

	table.Add(testTableItem{Name: "Alice", Age: 30, City: "New York"})
	table.Add(testTableItem{Name: "Bob", Age: 25, City: "Los Angeles"})
	table.Add(testTableItem{Name: "Charlie", Age: 35, City: "Chicago"})
	table.Add(testTableItem{Name: "Diana", Age: 28, City: "New Orleans"})

	return table
}

func TestTable(t *testing.T) {

	table := testTable()
	assert.EqualValues(t, [][]string{
		{"Alice", "30", "New York"},
		{"Bob", "25", "Los Angeles"},
		{"Charlie", "35", "Chicago"},
		{"Diana", "28", "New Orleans"},
	}, table.Rows())

	i, ok := table.GetRowItem(1)
	assert.True(t, ok)
	assert.Equal(t, "Bob", i.Name)
	assert.Equal(t, 25, i.Age)
	assert.Equal(t, "Los Angeles", i.City)
}

func TestTableNoMatch(t *testing.T) {
	table := testTable()
	table.SetFilter("xyz")
	assert.EqualValues(t, 0, len(table.Rows()))
}

func TestTableMatch(t *testing.T) {
	table := testTable()
	table.SetFilter("New")
	assert.EqualValues(t, [][]string{
		{"Alice", "30", "New York"},
		{"Diana", "28", "New Orleans"},
	}, table.Rows())
	item, ok := table.GetRowItem(0)
	assert.True(t, ok)
	assert.Equal(t, testTableItem{Name: "Alice", Age: 30, City: "New York"}, item)
	item, ok = table.GetRowItem(1)
	assert.True(t, ok)
	assert.Equal(t, testTableItem{Name: "Diana", Age: 28, City: "New Orleans"}, item)
	_, ok = table.GetRowItem(2)
	assert.False(t, ok)
}
