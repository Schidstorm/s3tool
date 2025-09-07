package terminal

import (
	"github.com/rivo/tview"
)

type ProfileInfoBox struct {
	*tview.Flex
	table *tview.Table
}

func NewProfileInfoBox() *ProfileInfoBox {
	table := tview.NewTable().SetBorders(false)
	table.SetSelectable(false, false)
	table.Box.SetBorder(false)
	table.Box.SetTitle("Profile Info")
	table.Box.SetTitleAlign(tview.AlignLeft)

	flex := tview.NewFlex()
	flex.SetDirection(tview.FlexColumn)
	flex.AddItem(tview.NewBox(), 1, 0, false) // left padding
	flex.AddItem(table, 0, 1, false)

	info := &ProfileInfoBox{
		Flex:  flex,
		table: table,
	}

	return info
}

func (info *ProfileInfoBox) UpdateContext(c Context) {
	info.table.Clear()

	if c.S3Client() == nil {
		return
	}

	infoItems := s3ClientToInfos(c)
	for row, infoItem := range infoItems {
		info.table.SetCell(row, 0, tview.NewTableCell(infoItem.Title).SetStyle(DefaultStyle.Foreground(DefaultTheme.LabelColor)))
		info.table.SetCell(row, 1,
			tview.NewTableCell(infoItem.Info).
				SetStyle(DefaultStyle.Foreground(DefaultTheme.PrimaryColor)),
		)
	}
}

type s3ClientInfoItem struct {
	Title string
	Info  string
}

func s3ClientToInfos(c Context) []s3ClientInfoItem {
	items := []s3ClientInfoItem{}

	parameters := c.S3Client().ConnectionParameters(c.Bucket())
	if parameters.Endpoint != nil {
		items = append(items, s3ClientInfoItem{
			Title: "Endpoint:",
			Info:  *parameters.Endpoint,
		})
	}
	if parameters.Region != nil {
		items = append(items, s3ClientInfoItem{
			Title: "Region:",
			Info:  *parameters.Region,
		})
	}
	if c.Bucket() != "" {
		items = append(items, s3ClientInfoItem{
			Title: "Bucket:",
			Info:  c.Bucket(),
		})
	}

	return items
}
