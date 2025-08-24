package terminal

import (
	"context"
	"errors"

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

	box.load()

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
		EventKey(tcell.KeyRune, 'n', 0): Hotkey{
			Title:   "New Bucket",
			Handler: func(event *tcell.EventKey) *tcell.EventKey { b.newBucketForm(); return nil },
		},
	}
}

func (b *BucketsPage) load() {
	b.ListPage.ClearRows()

	paginator := b.context.S3Client().ListBuckets(context.Background())
	var buckets []types.Bucket
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(context.Background())
		if err != nil {
			b.context.SetError(err)
			return
		}
		buckets = append(buckets, page...)
	}

	b.ListPage.AddAll(buckets)
}

func (b *BucketsPage) newBucketForm() {
	b.context.Modal(func(close func()) tview.Primitive {
		form := tview.NewForm()
		form.AddInputField("Name", "", 20, nil, func(text string) {})
		form.AddInputField("Region", "", 20, nil, func(text string) {})
		form.AddButton("Create", func() {
			b.createBucket(form)
			close()
		})
		form.AddButton("Cancel", func() {
			close()
		})
		return form
	}, "newBucket", 40, 10)

}

func (b *BucketsPage) createBucket(form *tview.Form) {
	name := form.GetFormItemByLabel("Name").(*tview.InputField).GetText()
	region := form.GetFormItemByLabel("Region").(*tview.InputField).GetText()

	if name == "" {
		b.context.SetError(errors.New("bucket name cannot be empty"))
		return
	}

	err := b.context.S3Client().CreateBucket(context.Background(), name, region)
	if err != nil {
		b.context.SetError(err)
		return
	}
	b.load()

}
