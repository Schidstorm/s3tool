package s3lib

import (
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

func TestMemoryClientFactoryBuildsBucketAndObject(t *testing.T) {
	now := time.Now().UTC()
	client := NewMemoryClientFactory().
		WithBucket("photos", "eu-central-1", now).
		WithObject("photos", "a.jpg", 5, now, "etag-1", "STANDARD", []byte("img")).
		Build()

	bucket, ok := client.buckets["photos"]
	if !ok {
		t.Fatal("expected photos bucket to exist")
	}
	if bucket.region != "eu-central-1" {
		t.Fatalf("expected region eu-central-1, got %s", bucket.region)
	}
	if len(bucket.objects) != 1 {
		t.Fatalf("expected one object, got %d", len(bucket.objects))
	}
	if bucket.objects[0].key != "a.jpg" {
		t.Fatalf("expected object key a.jpg, got %s", bucket.objects[0].key)
	}
}

func TestMemoryClientFactoryWithObjectWithoutBucketNoop(t *testing.T) {
	now := time.Now().UTC()
	client := NewMemoryClientFactory().
		WithObject("missing", "a.jpg", 5, now, "etag-1", "STANDARD", []byte("img")).
		Build()

	if len(client.buckets) != 0 {
		t.Fatalf("expected no buckets, got %d", len(client.buckets))
	}
}

func TestObjectConstructorsAndKindChecks(t *testing.T) {
	file := NewObjectFile(types.Object{Key: aws.String("file.txt")})
	dir := NewObjectDirectory("folder/")

	if !file.IsFile() || file.IsDirectory() {
		t.Fatal("expected file object kind checks to match")
	}
	if !dir.IsDirectory() || dir.IsFile() {
		t.Fatal("expected directory object kind checks to match")
	}
	if aws.ToString(dir.Object.Key) != "folder/" {
		t.Fatalf("expected folder/ directory key, got %s", aws.ToString(dir.Object.Key))
	}
}
