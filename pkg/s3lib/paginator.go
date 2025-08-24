package s3lib

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type Paginator[T any] interface {
	NextPage(ctx context.Context, optFns ...func(*s3.Options)) ([]T, error)
	HasMorePages() bool
}
