package s3lib

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type AwsConnector struct {
	name string
}

func (c *AwsConnector) Name() string {
	return c.name
}

func (c *AwsConnector) Type() string {
	return "aws"
}

func (c *AwsConnector) CreateClient(ctx context.Context) (Client, error) {
	opts := []func(*config.LoadOptions) error{
		config.WithSharedConfigProfile(c.name),
	}

	sdkConfig, err := config.LoadDefaultConfig(ctx, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS SDK config for profile %s: %w", c.name, err)
	}

	return NewSdkClient(s3.NewFromConfig(sdkConfig)), nil
}
