package terminal

import (
	"testing"
	"time"

	"github.com/schidstorm/s3tool/internal/s3lib"
	"github.com/stretchr/testify/assert"
)

func TestObjectsPage(t *testing.T) {
	europeLocale := time.FixedZone("Europe/Berin", int((1 * time.Hour).Seconds()))
	time.Local = europeLocale

	client := s3lib.NewMemoryClientFactory().
		WithBucket("test-bucket", "eu-central-1", time.Now()).
		WithObject("test-bucket", "file1.txt", 1024, time.Date(2023, 10, 1, 11, 0, 0, 0, time.UTC), "etag1", "STANDARD", []byte("content1")).
		WithObject("test-bucket", "file2.txt", 2049, time.Date(2022, 10, 2, 12, 0, 0, 0, europeLocale), "etag2", "STANDARD", []byte("content2")).
		Build()
	page := NewObjectsPage(NewContext().WithErrorFunc(func(err error) {
		t.Error(err)
	}).WithClient(client).WithBucket("test-bucket"))
	page.Load()
	rows := getTableRows(page.ListPage.tviewTable)
	assert.Equal(t, 3, len(rows))

	assert.EqualValues(t, [][]string{
		{"Name", "Size", "Last Modified"},
		{"file1.txt", "1 KiB", "2023-10-01 12:00:00"},
		{"file2.txt", "2 KiB", "2022-10-02 12:00:00"},
	}, rows)

	assert.NotNil(t, page)

}

func TestHumanizeSize(t *testing.T) {
	tests := []struct {
		size     *int64
		expected string
	}{
		{nil, ""},
		{new(int64), "0 B"},
		{int64Ptr(1023), "1023 B"},
		{int64Ptr(1024), "1 KiB"},
		{int64Ptr(1048576), "1 MiB"},
	}

	for _, test := range tests {
		result := humanizeSize(test.size)
		if result != test.expected {
			t.Errorf("Expected %s, got %s", test.expected, result)
		}
	}
}

func int64Ptr(i int64) *int64 {
	return &i
}
