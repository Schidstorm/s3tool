package cli

import (
	"os"
	"strings"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	if cfg.ProfilesDirectory != "~/.s3tool" {
		t.Fatalf("expected default profiles dir ~/.s3tool, got %q", cfg.ProfilesDirectory)
	}
	if !cfg.Loaders.Aws || !cfg.Loaders.S3Tool || cfg.Loaders.Memory {
		t.Fatalf("unexpected default loaders: %#v", cfg.Loaders)
	}
}

func TestCleanupExpandsHome(t *testing.T) {
	cfg := cleanup(S3ToolCliConfig{ProfilesDirectory: "~/profiles"})
	home, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("user home dir failed: %v", err)
	}

	expected := home + "/profiles"
	if cfg.ProfilesDirectory != expected {
		t.Fatalf("expected %q, got %q", expected, cfg.ProfilesDirectory)
	}
}

func TestCleanupKeepsAbsolutePath(t *testing.T) {
	original := "/tmp/s3tool-profiles"
	cfg := cleanup(S3ToolCliConfig{ProfilesDirectory: original})
	if cfg.ProfilesDirectory != original {
		t.Fatalf("expected unchanged path %q, got %q", original, cfg.ProfilesDirectory)
	}
}

func TestRootCmdFlags(t *testing.T) {
	cfg := S3ToolCliConfig{}
	cmd := rootCmd(&cfg)

	profilesFlag := cmd.Flags().Lookup("profiles")
	if profilesFlag == nil {
		t.Fatal("profiles flag not registered")
	}

	memoryFlag := cmd.Flags().Lookup("loaders.memory")
	if memoryFlag == nil {
		t.Fatal("loaders.memory flag not registered")
	}
	if !memoryFlag.Hidden {
		t.Fatal("loaders.memory flag should be hidden")
	}
}

func TestParseUpdatesGlobalConfig(t *testing.T) {
	original := Config
	Config = DefaultConfig()
	t.Cleanup(func() {
		Config = original
	})

	runApp, err := ParseAndShouldRun([]string{"--profiles", "~/custom", "--loaders.aws=false", "--loaders.s3tool=false"})
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	if !runApp {
		t.Fatal("expected Parse to run app for root command")
	}

	if Config.Loaders.Aws {
		t.Fatal("expected aws loader disabled")
	}
	if Config.Loaders.S3Tool {
		t.Fatal("expected s3tool loader disabled")
	}
	if !strings.Contains(Config.ProfilesDirectory, "/custom") {
		t.Fatalf("expected profiles dir to contain /custom, got %q", Config.ProfilesDirectory)
	}
}

func TestParseInvalidFlagReturnsError(t *testing.T) {
	original := Config
	t.Cleanup(func() {
		Config = original
	})

	_, err := ParseAndShouldRun([]string{"--does-not-exist"})
	if err == nil {
		t.Fatal("expected parse error for unknown flag")
	}
}

func TestParseCompletionReturnsNoRun(t *testing.T) {
	original := Config
	t.Cleanup(func() {
		Config = original
	})

	runApp, err := ParseAndShouldRun([]string{"completion", "--shell", "bash"})
	if err != nil {
		t.Fatalf("expected no error for completion command, got %v", err)
	}
	if runApp {
		t.Fatal("expected runApp false for completion subcommand")
	}
}
