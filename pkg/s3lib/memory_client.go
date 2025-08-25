package s3lib

import (
	"context"
	"errors"
	"os"
	"slices"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

type MemoryClient struct {
	buckets map[string]*MemoryBucket
}

func NewMemoryClient() *MemoryClient {
	client := &MemoryClient{
		buckets: make(map[string]*MemoryBucket),
	}
	return client
}

func (c *MemoryClient) ConnectionParameters(bucket string) ConnectionParameters {
	return ConnectionParameters{
		Endpoint: aws.String("memory"),
		Region:   nil,
	}
}

func (c *MemoryClient) ListBuckets(ctx context.Context) Paginator[types.Bucket] {
	var bucketList []types.Bucket
	for name := range c.buckets {
		bucketList = append(bucketList, types.Bucket{
			Name:         aws.String(name),
			CreationDate: aws.Time(c.buckets[name].creationDate),
			BucketRegion: aws.String(c.buckets[name].region),
		})
	}
	slices.SortFunc(bucketList, func(a, b types.Bucket) int {
		return strings.Compare(aws.ToString(a.Name), aws.ToString(b.Name))
	})
	return &memoryPaginator[types.Bucket]{items: bucketList}
}

func (c *MemoryClient) ListObjects(ctx context.Context, bucket string, prefix string) Paginator[Object] {
	if memBucket, exists := c.buckets[bucket]; exists {
		var objects []Object
		var directories []string
		for _, obj := range memBucket.objects {
			if after, ok := strings.CutPrefix(obj.key, prefix); ok {
				if strings.Contains(after, "/") {
					// It's a directory
					dir := strings.SplitN(after, "/", 2)[0] + "/"
					directories = append(directories, dir)
				} else {
					// It's a file
					objects = append(objects, NewObjectFile(types.Object{
						Key:          aws.String(obj.key),
						Size:         aws.Int64(obj.size),
						LastModified: aws.Time(obj.lastModified),
						ETag:         aws.String(obj.etag),
						StorageClass: types.ObjectStorageClass(obj.storageClass),
					}))
				}
			}
		}

		slices.Sort(directories)
		uniqueDirs := slices.Compact(directories)
		result := make([]Object, 0, len(uniqueDirs)+len(objects))
		for _, dir := range uniqueDirs {
			result = append(result, NewObjectDirectory(dir))
		}
		result = append(result, objects...)

		return &memoryPaginator[Object]{items: result}
	}
	return &memoryPaginator[Object]{
		err: errors.New("bucket not found"),
	}
}

func (c *MemoryClient) CreateBucket(ctx context.Context, bucket, region string) error {
	if _, exists := c.buckets[bucket]; exists {
		return errors.New("bucket already exists")
	}
	c.buckets[bucket] = &MemoryBucket{}
	return nil
}

func (c *MemoryClient) UploadFile(ctx context.Context, bucket, key, filePath string) error {
	if memBucket, exists := c.buckets[bucket]; exists {
		data, err := os.ReadFile(filePath)
		if err != nil {
			return err
		}

		memBucket.objects = append(memBucket.objects, MemoryObject{
			key:          key,
			data:         data,
			size:         int64(len(data)), // Simulated size
			lastModified: time.Now(),
			etag:         "dummy-etag",
			storageClass: "STANDARD",
		})

		return nil
	}
	return errors.New("bucket not found")
}

func (c *MemoryClient) DownloadFile(ctx context.Context, bucket, key, filePath string) error {
	if memBucket, exists := c.buckets[bucket]; exists {
		for _, obj := range memBucket.objects {
			if obj.key == key {
				if err := os.WriteFile(filePath, obj.data, os.ModePerm); err != nil {
					return err
				}
				return nil
			}
		}
		return errors.New("object not found")
	}
	return errors.New("bucket not found")
}

func (c *MemoryClient) GetObject(ctx context.Context, bucket, key string) (ObjectMetadata, error) {
	if memBucket, exists := c.buckets[bucket]; exists {
		for _, obj := range memBucket.objects {
			if obj.key == key {
				return ObjectMetadata{
					Region:       memBucket.region,
					LastModified: &obj.lastModified,
					Size:         &obj.size,
					Type:         aws.String("image/png"),
					Key:          obj.key,
					Bucket:       bucket,
					Owner:        aws.String("memory-user"),
					Tags:         map[string]string{"Environment": "Test"},
					LegalHold:    "OFF",
					ETag:         aws.String(obj.etag),
				}, nil
			}
		}
		return ObjectMetadata{}, errors.New("object not found")
	}
	return ObjectMetadata{}, errors.New("bucket not found")
}

func (c *MemoryClient) DeleteBucket(ctx context.Context, bucket string) error {
	if _, exists := c.buckets[bucket]; exists {
		delete(c.buckets, bucket)
		return nil
	}
	return errors.New("bucket not found")
}

func (c *MemoryClient) DeleteObject(ctx context.Context, bucket, key string) error {
	if memBucket, exists := c.buckets[bucket]; exists {
		for i, obj := range memBucket.objects {
			if obj.key == key {
				memBucket.objects = append(memBucket.objects[:i], memBucket.objects[i+1:]...)
				return nil
			}
		}
		return errors.New("object not found")
	}
	return errors.New("bucket not found")
}
