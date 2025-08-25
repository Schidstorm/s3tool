package s3lib

import "context"

type MemoryConnector struct {
	name string
}

func (c *MemoryConnector) Name() string {
	return c.name
}

func (c *MemoryConnector) Type() string {
	return "memory"
}

func (c *MemoryConnector) CreateClient(ctx context.Context) (Client, error) {
	return NewMemoryClient(), nil
}
