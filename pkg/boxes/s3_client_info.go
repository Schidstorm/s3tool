package boxes

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/rivo/tview"
)

type S3ClientInfo struct {
	*tview.Table
}

func NewS3ClientInfo() *S3ClientInfo {
	table := tview.NewTable().SetBorders(false)
	table.SetSelectable(false, false)
	table.Box.SetBorder(true)
	table.Box.SetTitle("S3 Client Info")
	table.Box.SetTitleAlign(tview.AlignLeft)

	info := &S3ClientInfo{
		Table: table,
	}

	return info
}

func (info *S3ClientInfo) Update(client *s3.Client) {
	info.Clear()

	if client == nil {
		info.SetCell(0, 0, tview.NewTableCell("No S3 client available").SetTextColor(tview.Styles.PrimaryTextColor))
		return
	}

	infos := s3ClientToInfos(client)
	for name, value := range infos {
		row, col := infoNameToLocation(name)
		if row >= 0 && col >= 0 {
			info.SetCell(row, col, tview.NewTableCell(value).SetTextColor(tview.Styles.PrimaryTextColor))
		}
	}
}

func s3ClientToInfos(client *s3.Client) map[string]string {
	return map[string]string{
		"Region":   client.Options().Region,
		"Endpoint": aws.ToString(client.Options().BaseEndpoint),
	}
}

func infoNameToLocation(name string) (int, int) {
	switch name {
	case "Region":
		return 0, 1
	case "Endpoint":
		return 1, 1
	default:
		return -1, -1
	}
}
