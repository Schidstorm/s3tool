package s3lib

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

type Client interface {
	ConnectionParameters(bucket string) ConnectionParameters
	ListBuckets(ctx context.Context) Paginator[types.Bucket]
	ListObjects(ctx context.Context, bucket string, prefix string) Paginator[Object]
	CreateBucket(ctx context.Context, bucket, region string) error
	UploadFile(ctx context.Context, bucket, key, filePath string) error
	DownloadFile(ctx context.Context, bucket, key, filePath string) error
	GetObject(ctx context.Context, bucket, key string) (ObjectMetadata, error)
	DeleteBucket(ctx context.Context, bucket string) error
	DeleteObject(ctx context.Context, bucket, key string) error
}
