package s3lib

import (
	"os"
	"path/filepath"
	"testing"
)

func TestAwsLoaderLoadProfilesFromConfigAndCredentials(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	awsDir := filepath.Join(home, ".aws")
	mustMkdirAll(t, awsDir)

	config := `[default]
region = eu-central-1

[profile dev]
region = us-east-1

[profile qa]
region = eu-west-1
`
	credentials := `[default]
aws_access_key_id = x
aws_secret_access_key = y

[prod]
aws_access_key_id = x
aws_secret_access_key = y
`

	mustWriteTestFile(t, filepath.Join(awsDir, "config"), config)
	mustWriteTestFile(t, filepath.Join(awsDir, "credentials"), credentials)

	loader := &AwsLoader{}
	profiles, err := loader.Load()
	if err != nil {
		t.Fatalf("load profiles failed: %v", err)
	}
	if len(profiles) != 4 {
		t.Fatalf("expected 4 profiles, got %d", len(profiles))
	}

	names := map[string]bool{}
	for _, p := range profiles {
		names[p.Name()] = true
		if p.Type() != "aws" {
			t.Fatalf("expected aws type, got %s", p.Type())
		}
	}

	for _, expected := range []string{"default", "dev", "qa", "prod"} {
		if !names[expected] {
			t.Fatalf("expected profile %s in %v", expected, names)
		}
	}
}

func TestAwsLoaderLoadNoFilesReturnsEmptyProfiles(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	loader := &AwsLoader{}
	profiles, err := loader.Load()
	if err != nil {
		t.Fatalf("expected no error when aws files are missing, got %v", err)
	}
	if len(profiles) != 0 {
		t.Fatalf("expected zero profiles, got %d", len(profiles))
	}
}

func TestAwsConnectorNameAndType(t *testing.T) {
	connector := &AwsConnector{name: "dev"}
	if connector.Name() != "dev" {
		t.Fatalf("expected name dev, got %s", connector.Name())
	}
	if connector.Type() != "aws" {
		t.Fatalf("expected type aws, got %s", connector.Type())
	}
}

func mustMkdirAll(t *testing.T, dir string) {
	t.Helper()
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatalf("mkdir %s failed: %v", dir, err)
	}
}

func mustWriteTestFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("write file %s failed: %v", path, err)
	}
}
