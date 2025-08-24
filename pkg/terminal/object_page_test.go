package terminal

import (
	"testing"
	"time"

	"github.com/schidstorm/s3tool/pkg/s3lib"
	"github.com/stretchr/testify/assert"
)

func TestObjectPage(t *testing.T) {
	europeLocale := time.FixedZone("Europe/Berin", int((1 * time.Hour).Seconds()))
	time.Local = europeLocale

	client := s3lib.NewMemoryClientFactory().
		WithBucket("test-bucket", "eu-central-1", time.Date(2023, 10, 1, 11, 0, 0, 0, time.UTC)).
		WithObject("test-bucket", "file1.txt", 1024, time.Date(2023, 10, 1, 11, 0, 0, 0, time.UTC), "etag1", "STANDARD", []byte("content1")).
		WithObject("test-bucket", "file2.txt", 2049, time.Date(2022, 10, 2, 12, 0, 0, 0, europeLocale), "etag2", "STANDARD", []byte("content2")).
		Build()
	page := NewObjectPage(NewContext().WithClient(client).WithErrorFunc(func(err error) {
		t.Error(err)
	}).WithBucket("test-bucket").WithObjectKey("file1.txt"))
	rows := getTableRows(page.Table)

	assert.EqualValues(t, [][]string{
		{"Bucket", "test-bucket"},
		{"Name", "file1.txt"},
		{"Region", "eu-central-1"},
		{"Owner", "memory-user"},
		{"Type", "image/png"},
		{"Size", "1 KiB"},
		{"ETag", "etag1"},
		{"LegalHold", "OFF"},
		{"LastModified", "2023-10-01 12:00:00"},
		{"Tags", "Environment=Test"},
	}, rows)
}
