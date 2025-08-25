package s3lib

import "context"

type Connector interface {
	Name() string
	Type() string
	CreateClient(ctx context.Context) (Client, error)
}
