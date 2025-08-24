package main

import (
	"context"

	"github.com/schidstorm/s3tool/pkg/s3lib"
)

type ScreenConnector struct {
	name string
}

func (c *ScreenConnector) Name() string {
	return c.name
}

func (c *ScreenConnector) Type() string {
	return "aws"
}

func (c *ScreenConnector) CreateClient(ctx context.Context) (s3lib.Client, error) {
	return newScreenClient(), nil
}
