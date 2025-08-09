package profile

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type Loader interface {
	LoadProfiles() ([]Connector, error)
}

type Connector interface {
	Name() string
	Type() string
	CreateClient(ctx context.Context) (*s3.Client, error)
}
