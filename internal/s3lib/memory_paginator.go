package s3lib

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type memoryPaginator[T any] struct {
	err   error
	items []T
	read  bool
}

func (p *memoryPaginator[T]) NextPage(ctx context.Context, optFns ...func(*s3.Options)) ([]T, error) {
	if p.err != nil {
		return nil, p.err
	}
	if p.read {
		return nil, nil
	}
	p.read = true
	return p.items, nil
}

func (p *memoryPaginator[T]) HasMorePages() bool {
	if p.err != nil {
		return false
	}
	if p.read {
		return false
	}
	return len(p.items) > 0
}
