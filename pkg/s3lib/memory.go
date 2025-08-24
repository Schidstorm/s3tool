package s3lib

import (
	"time"
)

type MemoryBucket struct {
	objects      []MemoryObject
	region       string
	creationDate time.Time
}

type MemoryObject struct {
	key          string
	size         int64
	lastModified time.Time
	etag         string
	storageClass string
	data         []byte
}
