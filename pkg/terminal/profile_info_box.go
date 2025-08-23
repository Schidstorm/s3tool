package terminal

import (
	"github.com/rivo/tview"
	"github.com/schidstorm/s3tool/pkg/s3lib"
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

func (info *ProfileInfoBox) Update(client s3lib.Client, bucketName string) {
	info.table.Clear()

	if client == nil {
		return
	}

	infoItems := s3ClientToInfos(client, bucketName)
	for row, infoItem := range infoItems {
		info.table.SetCell(row, 0, tview.NewTableCell(infoItem.Title).SetStyle(DefaultTheme.ProfileKey))
		info.table.SetCell(row, 1,
			tview.NewTableCell(infoItem.Info).
				SetStyle(DefaultTheme.ProfileValue),
		)
	}
}

type s3ClientInfoItem struct {
	Title string
	Info  string
}

func s3ClientToInfos(client s3lib.Client, bucketName string) []s3ClientInfoItem {
	items := []s3ClientInfoItem{}

	parameters := client.ConnectionParameters(bucketName)
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

	return items
}
