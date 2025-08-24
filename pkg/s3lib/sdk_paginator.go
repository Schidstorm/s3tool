package s3lib

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

type ListBucketsPaginator struct {
	*s3.ListBucketsPaginator
}

func (p ListBucketsPaginator) NextPage(ctx context.Context, optFns ...func(*s3.Options)) ([]types.Bucket, error) {
	output, err := p.ListBucketsPaginator.NextPage(ctx, optFns...)
	if err != nil {
		return nil, err
	}

	return output.Buckets, nil
}

func (p ListBucketsPaginator) HasMorePages() bool {
	return p.ListBucketsPaginator.HasMorePages()
}

type ListObjectsPaginator struct {
	*s3.ListObjectsV2Paginator
}

func (p ListObjectsPaginator) NextPage(ctx context.Context, optFns ...func(*s3.Options)) ([]Object, error) {
	output, err := p.ListObjectsV2Paginator.NextPage(ctx, optFns...)
	if err != nil {
		return nil, err
	}

	result := make([]Object, 0, len(output.Contents)+len(output.CommonPrefixes))
	for _, prefix := range output.CommonPrefixes {
		result = append(result, NewObjectDirectory(aws.ToString(prefix.Prefix)))
	}
	for _, obj := range output.Contents {
		result = append(result, NewObjectFile(obj))
	}

	return result, nil
}

func (p ListObjectsPaginator) HasMorePages() bool {
	return p.ListObjectsV2Paginator.HasMorePages()
}
