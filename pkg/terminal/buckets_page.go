package terminal

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/schidstorm/s3tool/pkg/s3lib"
)

type BucketsPage struct {
	*ListPage

	client s3lib.Client
}

func NewBucketsPage(client s3lib.Client) *BucketsPage {
	listPage := NewListPage()

	box := &BucketsPage{
		ListPage: listPage,
		client:   client,
	}

	listPage.SetSelectedFunc(func(columns []string) {
		if len(columns) < 1 {
			return
		}

		bucketName := columns[0]
		activeApp.SetS3Client(client, bucketName)
		activeApp.OpenPage(AttachClose{
			PageContent: NewObjectsPage(box.client, bucketName, ""),
			Closer: CloseFunc(func() {
				activeApp.SetS3Client(client, "")
			}),
		})
	})

	box.load()

	return box
}

func (b *BucketsPage) Title() string {
	return "Buckets"
}

func (b *BucketsPage) Hotkeys() map[tcell.EventKey]Hotkey {
	return map[tcell.EventKey]Hotkey{
		EventKey(tcell.KeyRune, 'n', 0): Hotkey{
			Title:   "New Bucket",
			Handler: func(event *tcell.EventKey) *tcell.EventKey { b.newBucketForm(); return nil },
		},
	}
}

func (b *BucketsPage) load() {
	b.ListPage.ClearRows()
	b.ListPage.AddRow(Row{
		Header:  true,
		Columns: []string{"Bucket Name", "Region", "Created At"},
	})

	paginator := b.client.ListBuckets(context.Background())
	var buckets []types.Bucket
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(context.Background())
		if err != nil {
			activeApp.SetError(err)
			return
		}
		buckets = append(buckets, page...)
	}

	rows := make([]Row, len(buckets))
	for i, bucket := range buckets {
		rows[i] = Row{
			Header: false,
			Columns: []string{
				aws.ToString(bucket.Name),
				aws.ToString(bucket.BucketRegion),
				humanizeTime(bucket.CreationDate),
			},
		}
	}
	b.ListPage.AddRows(rows)
}

func (b *BucketsPage) newBucketForm() {
	modalName := "newBucket"
	form := tview.NewForm()
	form.AddInputField("Name", "", 20, nil, func(text string) {})
	form.AddInputField("Region", "", 20, nil, func(text string) {})
	form.AddButton("Create", func() {
		b.createBucket(form)
		activeApp.CloseModal(modalName)
	})
	form.AddButton("Cancel", func() {
		activeApp.CloseModal(modalName)
	})

	activeApp.Modal(form, modalName, 40, 10)

}

func (b *BucketsPage) createBucket(form *tview.Form) {
	name := form.GetFormItemByLabel("Name").(*tview.InputField).GetText()
	region := form.GetFormItemByLabel("Region").(*tview.InputField).GetText()

	if name == "" {
		activeApp.SetError(errors.New("bucket name cannot be empty"))
		return
	}

	err := b.client.CreateBucket(context.Background(), name, region)
	if err != nil {
		activeApp.SetError(err)
		return
	}
	b.load()

}
