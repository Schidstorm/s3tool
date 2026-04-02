package terminal

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type BucketsPage struct {
	*ListPage[types.Bucket]

	context Context
}

func NewBucketsPage(context Context) *BucketsPage {
	listPage := NewListPage[types.Bucket]()
	listPage.AddColumn("Bucket Name", func(item types.Bucket) string { return aws.ToString(item.Name) })
	listPage.AddColumn("Region", func(item types.Bucket) string { return aws.ToString(item.BucketRegion) })
	listPage.AddColumn("Created At", func(item types.Bucket) string { return humanizeTime(item.CreationDate) })

	box := &BucketsPage{
		ListPage: listPage,
		context:  context,
	}

	listPage.SetSelectedFunc(func(selected types.Bucket) {
		context.OpenPage(NewObjectsPage(context.WithBucket(aws.ToString(selected.Name))))
	})

	return box
}

func (b *BucketsPage) Title() string {
	return "Buckets"
}

func (b *BucketsPage) Context() Context {
	return b.context
}

func (b *BucketsPage) Hotkeys() map[tcell.EventKey]Hotkey {
	return map[tcell.EventKey]Hotkey{
		EventKey(tcell.KeyRune, 'n', 0): {
			Title:   "New Bucket",
			Handler: func(event *tcell.EventKey) *tcell.EventKey { b.newBucketForm(); return nil },
		},
		EventKey(tcell.KeyRune, 'd', 0): {
			Title: "Delete Bucket",
			Handler: func(event *tcell.EventKey) *tcell.EventKey {
				items := b.table.GetHighlightedItems()
				if len(items) == 0 {
					if obj, ok := b.ListPage.GetSelectedRow(); ok {
						items = []types.Bucket{obj}
					}
				}

				if len(items) == 0 {
					return nil
				}

				bucketKey := "bucket"
				if len(items) != 1 {
					bucketKey = "buckets\n"
				}
				modalMessage := fmt.Sprintf("Are you sure you want to delete the %d %s '%s'?", len(items), bucketKey, strings.Join(func() []string {
					keys := make([]string, len(items))
					for i, item := range items {
						keys[i] = aws.ToString(item.Name)
					}
					return keys
				}(), ", "))

				b.context.Modal(ConfirmModal(modalMessage, func() {
					for _, item := range items {
						b.deleteBucket(item)
					}
				}))

				return nil
			},
		},
	}
}

func (b *BucketsPage) deleteBucket(bucket types.Bucket) {
	err := b.context.S3Client().DeleteBucket(context.Background(), aws.ToString(bucket.Name))
	if err != nil {
		b.context.SetError(err)
	}

	err = b.Load()
	if err != nil {
		b.context.SetError(err)
	}
}

func (b *BucketsPage) Load() error {
	b.ListPage.ClearRows()

	paginator := b.context.S3Client().ListBuckets(context.Background())
	var buckets []types.Bucket
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(context.Background())
		if err != nil {
			return err
		}
		buckets = append(buckets, page...)
	}

	b.ListPage.AddAll(buckets)
	return nil
}

func (b *BucketsPage) newBucketForm() {
	b.context.Modal(func(close func()) tview.Primitive {
		return NewModal().
			SetTitle("New Bucket").
			AddInput().SetLabel("Region").
			AddInput().SetLabel("Name").
			AddButtons([]string{"Create", "Cancel"}).
			SetDoneFunc(func(buttonLabel string, values map[string]string) {
				close()
				if buttonLabel == "Create" {
					b.createBucket(values)
				}
			})
	})
}

func (b *BucketsPage) createBucket(values map[string]string) {
	name := values["Name"]
	region := values["Region"]

	if name == "" {
		b.context.SetError(errors.New("bucket name cannot be empty"))
		return
	}

	err := b.context.S3Client().CreateBucket(context.Background(), name, region)
	if err != nil {
		b.context.SetError(err)
		return
	}

	err = b.Load()
	if err != nil {
		b.context.SetError(err)
		return
	}

}
