package main

import (
	"context"
	"math/rand"
	"strconv"
	"time"

	"github.com/schidstorm/s3tool/pkg/s3lib"
)

var memoryBucketLargeObjectCount = 100

func newScreenClient() s3lib.Client {
	factory := s3lib.NewMemoryClientFactory()
	factory.WithBucket("large-bucket", "us-east-1", time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC))
	rng := randomNameGenerator(context.Background())
	for range memoryBucketLargeObjectCount {
		factory.WithObject("large-bucket", "directory/"+<-rng, 0, time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC), "dummy-etag", "STANDARD", nil)
	}
	for range memoryBucketLargeObjectCount {
		factory.WithObject("large-bucket", <-rng, 0, time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC), "dummy-etag", "STANDARD", nil)
	}

	return factory.Build()
}

func randomNameGenerator(ctx context.Context) chan string {
	result := make(chan string)
	rng := rand.New(rand.NewSource(0))

	go func() {
		defer close(result)
		for i := 0; ; i++ {
			select {
			case <-ctx.Done():
				return
			default:
			}
			n := rng.Int63()
			nHex := strconv.FormatInt(n, 16)
			name := "object-" + strconv.FormatInt(int64(i), 10) + "-" + nHex
			result <- name
		}
	}()

	return result
}
