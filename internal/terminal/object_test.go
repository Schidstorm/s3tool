package terminal

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/schidstorm/s3tool/internal/s3lib"
)

type objectTestClient struct {
	downloadData []byte
	downloadErr  error
	uploadErr    error
	uploadCount  int
}

func (c *objectTestClient) ConnectionParameters(bucket string) s3lib.ConnectionParameters {
	return s3lib.ConnectionParameters{}
}

func (c *objectTestClient) ListBuckets(ctx context.Context) s3lib.Paginator[types.Bucket] {
	return nil
}

func (c *objectTestClient) ListObjects(ctx context.Context, bucket string, prefix string) s3lib.Paginator[s3lib.Object] {
	return nil
}

func (c *objectTestClient) CreateBucket(ctx context.Context, bucket, region string) error {
	return nil
}

func (c *objectTestClient) UploadFile(ctx context.Context, bucket, key, filePath string) error {
	if c.uploadErr != nil {
		return c.uploadErr
	}
	c.uploadCount++
	return nil
}

func (c *objectTestClient) DownloadFile(ctx context.Context, bucket, key, filePath string) error {
	if c.downloadErr != nil {
		return c.downloadErr
	}
	return os.WriteFile(filePath, c.downloadData, 0o600)
}

func (c *objectTestClient) GetObject(ctx context.Context, bucket, key string) (s3lib.ObjectMetadata, error) {
	return s3lib.ObjectMetadata{}, nil
}

func (c *objectTestClient) DeleteBucket(ctx context.Context, bucket string) error {
	return nil
}

func (c *objectTestClient) DeleteObject(ctx context.Context, bucket, key string) error {
	return nil
}

func testContextWithClient(client s3lib.Client) Context {
	return NewContext().
		WithClient(client).
		WithBucket("bucket").
		WithObjectKey("object.txt").
		WithSuspendAppFunc(func(f func()) bool {
			f()
			return true
		})
}

func TestFileHashSameAndDifferent(t *testing.T) {
	dir := t.TempDir()
	p1 := filepath.Join(dir, "a.txt")
	p2 := filepath.Join(dir, "b.txt")
	p3 := filepath.Join(dir, "c.txt")

	if err := os.WriteFile(p1, []byte("same"), 0o600); err != nil {
		t.Fatalf("write p1 failed: %v", err)
	}
	if err := os.WriteFile(p2, []byte("same"), 0o600); err != nil {
		t.Fatalf("write p2 failed: %v", err)
	}
	if err := os.WriteFile(p3, []byte("different"), 0o600); err != nil {
		t.Fatalf("write p3 failed: %v", err)
	}

	h1, err := fileHash(p1)
	if err != nil {
		t.Fatalf("hash p1 failed: %v", err)
	}
	h2, err := fileHash(p2)
	if err != nil {
		t.Fatalf("hash p2 failed: %v", err)
	}
	h3, err := fileHash(p3)
	if err != nil {
		t.Fatalf("hash p3 failed: %v", err)
	}

	if h1 != h2 {
		t.Fatal("expected equal hash for equal content")
	}
	if h1 == h3 {
		t.Fatal("expected different hash for different content")
	}
}

func TestDownloadFileToTmp(t *testing.T) {
	client := &objectTestClient{downloadData: []byte("payload")}
	ctx := testContextWithClient(client)

	path, err := downloadFileToTmp(ctx)
	if err != nil {
		t.Fatalf("downloadFileToTmp failed: %v", err)
	}

	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read temp file failed: %v", err)
	}
	if string(content) != "payload" {
		t.Fatalf("expected payload, got %q", string(content))
	}
}

func TestEditObjectNoChangeNoUpload(t *testing.T) {
	prevEdit := editCommand
	editCommand = []string{"sh", "-c", "true", "sh"}
	t.Cleanup(func() { editCommand = prevEdit })

	client := &objectTestClient{downloadData: []byte("payload")}
	ctx := testContextWithClient(client)

	if err := editObject(ctx); err != nil {
		t.Fatalf("editObject failed: %v", err)
	}
	if client.uploadCount != 0 {
		t.Fatalf("expected no upload on unchanged file, got %d", client.uploadCount)
	}
}

func TestEditObjectChangedUploads(t *testing.T) {
	prevEdit := editCommand
	editCommand = []string{"sh", "-c", "echo changed >> \"$1\"", "sh"}
	t.Cleanup(func() { editCommand = prevEdit })

	client := &objectTestClient{downloadData: []byte("payload")}
	ctx := testContextWithClient(client)

	if err := editObject(ctx); err != nil {
		t.Fatalf("editObject failed: %v", err)
	}
	if client.uploadCount != 1 {
		t.Fatalf("expected one upload on changed file, got %d", client.uploadCount)
	}
}

func TestEditObjectMissingFileAfterEdit(t *testing.T) {
	prevEdit := editCommand
	editCommand = []string{"sh", "-c", "rm -f \"$1\"", "sh"}
	t.Cleanup(func() { editCommand = prevEdit })

	client := &objectTestClient{downloadData: []byte("payload")}
	ctx := testContextWithClient(client)

	err := editObject(ctx)
	if err == nil {
		t.Fatal("expected error when file is deleted during edit")
	}
}

func TestViewObject(t *testing.T) {
	prevView := viewCommand
	viewCommand = []string{"sh", "-c", "true", "sh"}
	t.Cleanup(func() { viewCommand = prevView })

	client := &objectTestClient{downloadData: []byte("payload")}
	ctx := testContextWithClient(client)

	if err := viewObject(ctx); err != nil {
		t.Fatalf("viewObject failed: %v", err)
	}
}

func TestDownloadFileToTmpDownloadError(t *testing.T) {
	client := &objectTestClient{downloadErr: errors.New("download failed")}
	ctx := testContextWithClient(client)

	if _, err := downloadFileToTmp(ctx); err == nil {
		t.Fatal("expected error from downloadFileToTmp")
	}
}
