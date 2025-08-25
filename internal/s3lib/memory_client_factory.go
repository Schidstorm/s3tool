package s3lib

import "time"

type MemoryClientFactory struct {
	client *MemoryClient
}

func NewMemoryClientFactory() *MemoryClientFactory {
	return &MemoryClientFactory{
		client: NewMemoryClient(),
	}
}

func (f *MemoryClientFactory) WithBucket(bucket string, region string, creationDate time.Time) *MemoryClientFactory {
	f.client.buckets[bucket] = &MemoryBucket{
		region:       region,
		creationDate: creationDate,
	}
	return f
}

func (f *MemoryClientFactory) WithObject(bucket, key string, size int64, lastModified time.Time, etag, storageClass string, data []byte) *MemoryClientFactory {
	if memBucket, exists := f.client.buckets[bucket]; exists {
		memBucket.objects = append(memBucket.objects, MemoryObject{
			key:          key,
			size:         size,
			lastModified: lastModified,
			etag:         etag,
			storageClass: storageClass,
			data:         data,
		})
	}
	return f
}

func (f *MemoryClientFactory) Build() *MemoryClient {
	return f.client
}
