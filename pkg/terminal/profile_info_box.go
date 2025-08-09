package terminal

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/rivo/tview"
)

type ProfileInfoBox struct {
	*tview.Table
}

func NewProfileInfoBox() *ProfileInfoBox {
	table := tview.NewTable().SetBorders(false)
	table.SetSelectable(false, false)
	table.Box.SetBorder(true)
	table.Box.SetTitle("Profile Info")
	table.Box.SetTitleAlign(tview.AlignLeft)

	info := &ProfileInfoBox{
		Table: table,
	}

	return info
}

func (info *ProfileInfoBox) Update(client *s3.Client, bucketName string) {
	info.Clear()

	if client == nil {
		info.SetCell(0, 0, tview.NewTableCell("No S3 client available").SetTextColor(tview.Styles.PrimaryTextColor))
		return
	}

	infoItems := s3ClientToInfos(client, bucketName)
	for row, infoItem := range infoItems {
		info.SetCell(row, 0, tview.NewTableCell(infoItem.Title).SetTextColor(tview.Styles.SecondaryTextColor))
		info.SetCell(row, 1, tview.NewTableCell(infoItem.Info).SetTextColor(tview.Styles.PrimaryTextColor))
	}
}

type s3ClientInfoItem struct {
	Title string
	Info  string
}

func s3ClientToInfos(client *s3.Client, bucketName string) []s3ClientInfoItem {
	items := []s3ClientInfoItem{
		{Title: "Region", Info: client.Options().Region},
	}

	ep, err := client.Options().EndpointResolverV2.ResolveEndpoint(context.Background(), s3.EndpointParameters{
		Bucket:         aws.String(bucketName),
		Region:         aws.String(client.Options().Region),
		UseFIPS:        aws.Bool(client.Options().EndpointOptions.UseFIPSEndpoint == aws.FIPSEndpointStateEnabled),
		UseDualStack:   aws.Bool(client.Options().EndpointOptions.UseDualStackEndpoint == aws.DualStackEndpointStateEnabled),
		Endpoint:       client.Options().BaseEndpoint,
		ForcePathStyle: aws.Bool(client.Options().UsePathStyle),
		Accelerate:     aws.Bool(client.Options().UseAccelerate),
	})
	if err == nil {
		items = append(items, s3ClientInfoItem{Title: "Endpoint", Info: ep.URI.String()})
	}

	return items
}
