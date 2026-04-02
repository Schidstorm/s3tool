package s3lib

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/schidstorm/s3tool/internal/cli"
)

func TestS3ToolLoaderLoadProfiles(t *testing.T) {
	tmpDir := t.TempDir()

	originalConfig := cli.Config
	cli.Config = &cli.S3ToolCliConfig{ProfilesDirectory: tmpDir}
	t.Cleanup(func() {
		cli.Config = originalConfig
	})

	mustWriteFile(t, filepath.Join(tmpDir, "minio.yaml"), "access_key_id: test\nsecret_access_key: secret\nregion: us-east-1\n")
	mustWriteFile(t, filepath.Join(tmpDir, "alt.yml"), "access_key_id: alt\nsecret_access_key: secret\nregion: us-west-1\n")
	mustWriteFile(t, filepath.Join(tmpDir, "ignored.txt"), "ignored")

	mustMkdir(t, filepath.Join(tmpDir, "nested"))

	loader := &S3ToolLoader{}
	profiles, err := loader.Load()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(profiles) != 2 {
		t.Fatalf("expected 2 profiles, got %d", len(profiles))
	}

	found := map[string]bool{}
	for _, profile := range profiles {
		found[profile.Name()] = true
		if profile.Type() != "s3tool" {
			t.Fatalf("expected connector type s3tool, got %s", profile.Type())
		}
	}

	if !found["minio"] || !found["alt"] {
		t.Fatalf("expected profiles minio and alt, got %#v", found)
	}
}

func TestS3ToolLoaderLoadInvalidYAML(t *testing.T) {
	tmpDir := t.TempDir()

	originalConfig := cli.Config
	cli.Config = &cli.S3ToolCliConfig{ProfilesDirectory: tmpDir}
	t.Cleanup(func() {
		cli.Config = originalConfig
	})

	mustWriteFile(t, filepath.Join(tmpDir, "broken.yaml"), "access_key_id: [broken")

	loader := &S3ToolLoader{}
	_, err := loader.Load()
	if err == nil {
		t.Fatal("expected yaml parse error, got nil")
	}
}

func TestS3ToolLoaderLoadMissingDirectory(t *testing.T) {
	originalConfig := cli.Config
	cli.Config = &cli.S3ToolCliConfig{ProfilesDirectory: filepath.Join(t.TempDir(), "does-not-exist")}
	t.Cleanup(func() {
		cli.Config = originalConfig
	})

	loader := &S3ToolLoader{}
	_, err := loader.Load()
	if err == nil {
		t.Fatal("expected read dir error, got nil")
	}
}

func mustWriteFile(t *testing.T, path string, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("failed to write %s: %v", path, err)
	}
}

func mustMkdir(t *testing.T, path string) {
	t.Helper()
	if err := os.MkdirAll(path, 0o755); err != nil {
		t.Fatalf("failed to create directory %s: %v", path, err)
	}
}
