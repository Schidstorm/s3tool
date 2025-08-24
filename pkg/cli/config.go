package cli

import (
	"flag"
	"os"
	"strings"
)

var Config = DefaultConfig()

type S3ToolCliConfig struct {
	ProfilesDirectory string
	Loaders           S3ToolCliConfigLoader `yaml:"loaders"`
}

type S3ToolCliConfigLoader struct {
	Aws    bool `yaml:"aws"`
	S3Tool bool `yaml:"s3tool"`
	Memory bool `yaml:"memory,omitempty"`
}

func DefaultConfig() *S3ToolCliConfig {
	return &S3ToolCliConfig{
		ProfilesDirectory: "~/.s3tool",
		Loaders: S3ToolCliConfigLoader{
			Aws:    true,
			S3Tool: true,
			Memory: false,
		},
	}
}

func Parse(args []string) error {
	var cfg S3ToolCliConfig
	flag.StringVar(&cfg.ProfilesDirectory, "profiles", Config.ProfilesDirectory, "Path to a directory containing profile yaml files")
	flag.BoolVar(&cfg.Loaders.Aws, "loaders.aws", Config.Loaders.Aws, "Enable AWS loader")
	flag.BoolVar(&cfg.Loaders.S3Tool, "loaders.s3tool", Config.Loaders.S3Tool, "Enable S3Tool loader")
	flag.BoolVar(&cfg.Loaders.Memory, "loaders.memory", Config.Loaders.Memory, "Enable Memory loader (for testing purposes)")

	err := flag.CommandLine.Parse(args)
	if err != nil {
		return err
	}

	cfg = cleanup(cfg)
	Config = &cfg

	return nil
}

func cleanup(cfg S3ToolCliConfig) S3ToolCliConfig {
	result := cfg
	if strings.HasPrefix(result.ProfilesDirectory, "~") {
		home, err := os.UserHomeDir()
		if err == nil {
			result.ProfilesDirectory = strings.TrimPrefix(result.ProfilesDirectory, "~")
			result.ProfilesDirectory = home + result.ProfilesDirectory
		}
	}

	return result
}
