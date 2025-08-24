package terminal

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestListPage(t *testing.T) {
	page := NewListPage[int]()
	page.AddColumn("Number", func(item int) string { return fmt.Sprintf("%d", item) })
	page.AddColumn("Square", func(item int) string { return fmt.Sprintf("%d", item*item) })

	for i := 1; i <= 5; i++ {
		page.Add(i)
	}

	rows := getTableRows(page.tviewTable)
	assert.Equal(t, 6, len(rows)) // 1 header + 5 items

	assert.EqualValues(t, [][]string{
		{"Number", "Square"},
		{"1", "1"},
		{"2", "4"},
		{"3", "9"},
		{"4", "16"},
		{"5", "25"},
	}, rows)

	page.SetSearch("2")
	rows = getTableRows(page.tviewTable)
	assert.Equal(t, 3, len(rows))

	assert.EqualValues(t, [][]string{
		{"Number", "Square"},
		{"2", "4"},
		{"5", "25"},
	}, rows)

	page.SetSearch("1337")
	rows = getTableRows(page.tviewTable)
	assert.Equal(t, 1, len(rows))

	assert.EqualValues(t, [][]string{
		{"Number", "Square"},
	}, rows)
}
