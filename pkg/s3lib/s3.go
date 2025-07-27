package s3lib

import "context"

type Bucket struct {
	Name string
}

type S3Lib interface {
	ListBuckets(ctx context.Context) ([]Bucket, error)
}
