package terminal

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/rivo/tview"
	"github.com/schidstorm/s3tool/pkg/s3lib"
	"github.com/stretchr/testify/assert"
)

func TestProfileInfoBoxPathStyle(t *testing.T) {
	client := s3.NewFromConfig(aws.Config{}, func(o *s3.Options) {
		o.Region = "us-west-2"
		o.UsePathStyle = true
		o.BaseEndpoint = aws.String("https://s3.us-west-2.amazonaws.com:654")
	})

	box := NewProfileInfoBox()
	box.Update(s3lib.NewSdkClient(client), "test-bucket")
	rows := getTableRows(box.table)
	assert.Equal(t, 2, len(rows))
	assert.EqualValues(t, [][]string{
		{"Endpoint", "https://s3.us-west-2.amazonaws.com:654/test-bucket"},
		{"Region", "us-west-2"},
	}, rows)
}

func TestProfileInfoBoxNoPathStyle(t *testing.T) {
	client := s3.NewFromConfig(aws.Config{}, func(o *s3.Options) {
		o.Region = "us-west-2"
		o.UsePathStyle = false
		o.BaseEndpoint = aws.String("https://asasd:654")
	})

	box := NewProfileInfoBox()
	box.Update(s3lib.NewSdkClient(client), "test-bucket")
	rows := getTableRows(box.table)
	assert.Equal(t, 2, len(rows))
	assert.EqualValues(t, [][]string{
		{"Endpoint", "https://test-bucket.asasd:654"},
		{"Region", "us-west-2"},
	}, rows)
}

func TestProfileInfoBoxNoPathStyleNoBucket(t *testing.T) {
	client := s3.NewFromConfig(aws.Config{}, func(o *s3.Options) {
		o.Region = "us-west-2"
		o.UsePathStyle = false
		o.BaseEndpoint = aws.String("https://asasd:654/")
	})

	box := NewProfileInfoBox()
	box.Update(s3lib.NewSdkClient(client), "")
	rows := getTableRows(box.table)
	assert.Equal(t, 2, len(rows))
	assert.EqualValues(t, [][]string{
		{"Endpoint", "https://asasd:654/"},
		{"Region", "us-west-2"},
	}, rows)
}

func TestProfileInfoBoxNoClient(t *testing.T) {
	box := NewProfileInfoBox()
	box.Update(nil, "")
	rows := getTableRows(box.table)
	assert.Equal(t, 1, len(rows))
	assert.Equal(t, "No S3 client available", rows[0][0])
}

func getTableRows(t *tview.Table) [][]string {
	rows := [][]string{}
	for row := 0; row < t.GetRowCount(); row++ {
		cells := []string{}
		for col := 0; col < t.GetColumnCount(); col++ {
			cell := t.GetCell(row, col)
			if cell != nil {
				cells = append(cells, cell.Text)
			} else {
				cells = append(cells, "")
			}
		}
		rows = append(rows, cells)
	}
	return rows
}
