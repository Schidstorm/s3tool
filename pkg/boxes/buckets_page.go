package boxes

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/schidstorm/s3tool/pkg/s3lib"
)

type BucketsPage struct {
	*tview.Flex

	list       *ListPage
	statusText *tview.TextView

	client *s3.Client
}

func NewBucketsBox(client *s3.Client) *BucketsPage {
	listPage := NewListPage()
	statusText := tview.NewTextView().SetTextAlign(tview.AlignCenter)

	flex := tview.NewFlex()
	flex.SetDirection(tview.FlexRow)
	flex.AddItem(listPage, 0, 1, true)
	flex.AddItem(statusText, 1, 1, false)

	box := &BucketsPage{
		Flex:       flex,
		list:       listPage,
		statusText: statusText,
		client:     client,
	}

	listPage.SetSelectedFunc(func(columns []string) {
		activeApp.OpenPage(NewObjectsPage(box.client, columns[0]))
	})

	listPage.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyRune:
			if event.Rune() == 'n' {
				box.newBucketForm()
				return nil
			}
		}
		return event
	})

	box.load()

	return box
}

func (b *BucketsPage) Title() string {
	return "Buckets"
}

func (b *BucketsPage) Hotkeys() map[string]string {
	return map[string]string{
		"n": "New Bucket",
	}
}

func (b *BucketsPage) SetSearch(search string) {
	b.list.SetSearch(search)
}

func (b *BucketsPage) load() {
	b.list.ClearRows()
	b.list.AddRow(Row{
		Header:  true,
		Columns: []string{"Bucket Name", "Region", "Created At"},
	})

	buckets, err := s3lib.ListBuckets(b.client, context.Background())
	if err != nil {
		b.setError(err)
		return
	}

	rows := make([]Row, len(buckets))
	for i, bucket := range buckets {
		rows[i] = Row{
			Header:  false,
			Columns: []string{aws.ToString(bucket.Name), aws.ToString(bucket.BucketRegion), aws.ToTime(bucket.CreationDate).Format("2006-01-02 15:04:05")},
		}
	}
	b.list.AddRows(rows)
}

func (b *BucketsPage) setError(err error) {
	b.statusText.SetText("Error: " + err.Error())
	b.statusText.SetTextColor(tcell.ColorRed)
}

func (b *BucketsPage) newBucketForm() {
	modalName := "newBucket"
	form := tview.NewForm()
	form.AddInputField("Name", "", 20, nil, func(text string) {})
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

	if name == "" {
		b.setError(errors.New("bucket name cannot be empty"))
		return
	}

	err := s3lib.CreateBucket(b.client, context.Background(), name, "")
	if err != nil {
		b.setError(err)
		return
	}
	b.load()

}
