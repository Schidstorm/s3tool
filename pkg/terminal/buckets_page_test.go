package terminal

import (
	"testing"
	"time"

	"github.com/schidstorm/s3tool/pkg/s3lib"
	"github.com/stretchr/testify/assert"
)

func TestBucketsPage(t *testing.T) {
	europeLocale := time.FixedZone("Europe/Berin", int((1 * time.Hour).Seconds()))
	time.Local = europeLocale

	client := s3lib.NewMemoryClientFactory().
		WithBucket("test-bucket", "eu-central-1", time.Date(2023, 10, 1, 11, 0, 0, 0, time.UTC)).
		WithBucket("test-bucket-2", "eu-central-2", time.Date(2022, 10, 2, 12, 0, 0, 0, europeLocale)).
		Build()
	page := NewBucketsPage(NewContext().WithClient(client).WithErrorFunc(func(err error) {
		t.Error(err)
	}))
	rows := getTableRows(page.ListPage.tviewTable)

	assert.EqualValues(t, [][]string{
		{"Bucket Name", "Region", "Created At"},
		{"test-bucket", "eu-central-1", "2023-10-01 12:00:00"},
		{"test-bucket-2", "eu-central-2", "2022-10-02 12:00:00"},
	}, rows)
}
