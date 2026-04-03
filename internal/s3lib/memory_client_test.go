package s3lib

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
)

func TestMemoryClientCreateBucketAndListBucketsSorted(t *testing.T) {
	client := NewMemoryClient()

	if err := client.CreateBucket(context.Background(), "z-bucket", "eu-central-1"); err != nil {
		t.Fatalf("create bucket z-bucket failed: %v", err)
	}
	if err := client.CreateBucket(context.Background(), "a-bucket", "us-east-1"); err != nil {
		t.Fatalf("create bucket a-bucket failed: %v", err)
	}

	if err := client.CreateBucket(context.Background(), "a-bucket", "us-east-1"); err == nil {
		t.Fatal("expected duplicate bucket error, got nil")
	}

	page, err := client.ListBuckets(context.Background()).NextPage(context.Background())
	if err != nil {
		t.Fatalf("list buckets failed: %v", err)
	}
	if len(page) != 2 {
		t.Fatalf("expected 2 buckets, got %d", len(page))
	}

	if aws.ToString(page[0].Name) != "a-bucket" || aws.ToString(page[1].Name) != "z-bucket" {
		t.Fatalf("expected sorted buckets [a-bucket z-bucket], got [%s %s]", aws.ToString(page[0].Name), aws.ToString(page[1].Name))
	}
}

func TestMemoryClientUploadDownloadAndDeleteObject(t *testing.T) {
	client := NewMemoryClient()
	if err := client.CreateBucket(context.Background(), "files", "us-east-1"); err != nil {
		t.Fatalf("create bucket failed: %v", err)
	}

	srcPath := filepath.Join(t.TempDir(), "src.txt")
	if err := os.WriteFile(srcPath, []byte("hello world"), 0o600); err != nil {
		t.Fatalf("write source file failed: %v", err)
	}

	if err := client.UploadFile(context.Background(), "files", "hello.txt", srcPath); err != nil {
		t.Fatalf("upload failed: %v", err)
	}

	dstPath := filepath.Join(t.TempDir(), "dst.txt")
	if err := client.DownloadFile(context.Background(), "files", "hello.txt", dstPath); err != nil {
		t.Fatalf("download failed: %v", err)
	}

	got, err := os.ReadFile(dstPath)
	if err != nil {
		t.Fatalf("read downloaded file failed: %v", err)
	}
	if string(got) != "hello world" {
		t.Fatalf("expected downloaded content hello world, got %q", string(got))
	}

	if err := client.DeleteObject(context.Background(), "files", "hello.txt"); err != nil {
		t.Fatalf("delete object failed: %v", err)
	}

	if err := client.DeleteObject(context.Background(), "files", "hello.txt"); err == nil {
		t.Fatal("expected object not found after delete, got nil")
	}
}

func TestMemoryClientListObjectsWithPrefixAndDirectories(t *testing.T) {
	now := time.Now().UTC()
	client := NewMemoryClientFactory().
		WithBucket("bucket", "eu-west-1", now).
		WithObject("bucket", "photos/2024/a.jpg", 10, now, "e1", "STANDARD", []byte("a")).
		WithObject("bucket", "photos/2024/b.jpg", 10, now, "e2", "STANDARD", []byte("b")).
		WithObject("bucket", "photos/root.txt", 8, now, "e3", "STANDARD", []byte("root")).
		WithObject("bucket", "videos/clip.mp4", 12, now, "e4", "STANDARD", []byte("clip")).
		Build()

	items, err := client.ListObjects(context.Background(), "bucket", "photos/").NextPage(context.Background())
	if err != nil {
		t.Fatalf("list objects failed: %v", err)
	}
	if len(items) != 2 {
		t.Fatalf("expected 2 items (one dir + one file), got %d", len(items))
	}

	if !items[0].IsDirectory() || aws.ToString(items[0].Object.Key) != "2024/" {
		t.Fatalf("expected first item directory 2024/, got kind=%v key=%s", items[0].Kind, aws.ToString(items[0].Object.Key))
	}
	if !items[1].IsFile() || aws.ToString(items[1].Object.Key) != "photos/root.txt" {
		t.Fatalf("expected second item file photos/root.txt, got kind=%v key=%s", items[1].Kind, aws.ToString(items[1].Object.Key))
	}
}

func TestMemoryClientErrorPaths(t *testing.T) {
	client := NewMemoryClient()

	if err := client.UploadFile(context.Background(), "missing", "k", "/tmp/nope"); err == nil {
		t.Fatal("expected bucket not found on upload")
	}
	if err := client.DownloadFile(context.Background(), "missing", "k", filepath.Join(t.TempDir(), "x")); err == nil {
		t.Fatal("expected bucket not found on download")
	}
	if err := client.DeleteBucket(context.Background(), "missing"); err == nil {
		t.Fatal("expected bucket not found on delete bucket")
	}
	if err := client.DeleteObject(context.Background(), "missing", "k"); err == nil {
		t.Fatal("expected bucket not found on delete object")
	}
	if _, err := client.GetObject(context.Background(), "missing", "k"); err == nil {
		t.Fatal("expected bucket not found on get object")
	}

	p := client.ListObjects(context.Background(), "missing", "")
	if p.HasMorePages() {
		t.Fatal("expected paginator with error to report no pages")
	}
	if _, err := p.NextPage(context.Background()); err == nil {
		t.Fatal("expected list objects paginator error")
	}
}

func TestMemoryClientGetObjectAndConnectionParameters(t *testing.T) {
	now := time.Now().UTC()
	client := NewMemoryClientFactory().
		WithBucket("bucket", "ap-south-1", now).
		WithObject("bucket", "doc.txt", 3, now, "etag-1", "STANDARD", []byte("doc")).
		Build()

	meta, err := client.GetObject(context.Background(), "bucket", "doc.txt")
	if err != nil {
		t.Fatalf("get object failed: %v", err)
	}
	if meta.Key != "doc.txt" || meta.Bucket != "bucket" {
		t.Fatalf("unexpected object metadata key/bucket: %s/%s", meta.Key, meta.Bucket)
	}
	if meta.ETag == nil || aws.ToString(meta.ETag) != "etag-1" {
		t.Fatalf("unexpected etag: %v", meta.ETag)
	}

	params := client.ConnectionParameters("bucket")
	if params.Endpoint == nil || aws.ToString(params.Endpoint) != "memory" {
		t.Fatalf("unexpected endpoint: %v", params.Endpoint)
	}
	if params.Region != nil {
		t.Fatalf("expected nil region for memory client, got %v", aws.ToString(params.Region))
	}
}
