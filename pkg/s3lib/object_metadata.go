package s3lib

import (
	"time"
)

type ObjectMetadata struct {
	Region       string
	LastModified *time.Time
	Size         *int64
	Type         *string
	Key          string
	Bucket       string
	Owner        *string
	Tags         map[string]string
	LegalHold    string
	ETag         *string
	Metadata     map[string]string
}
