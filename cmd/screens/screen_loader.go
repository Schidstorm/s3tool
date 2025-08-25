package main

import "github.com/schidstorm/s3tool/internal/s3lib"

type ScreenLoader struct {
}

func (l *ScreenLoader) Load() ([]s3lib.Connector, error) {
	return []s3lib.Connector{
		&ScreenConnector{
			name: "aws",
		},
	}, nil
}
