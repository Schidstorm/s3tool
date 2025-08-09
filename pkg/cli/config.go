package cli

import (
	"flag"
	"os"
	"strings"
)

var Config = S3ToolCliConfig{}

type S3ToolCliConfig struct {
	ProfilesDirectory string
}

func Parse(args []string) error {
	var cfg S3ToolCliConfig
	flag.StringVar(&cfg.ProfilesDirectory, "profiles", "~/.s3tool", "Path to a directory containing profile yaml files")

	err := flag.CommandLine.Parse(args)
	if err != nil {
		return err
	}

	cfg = cleanup(cfg)
	Config = cfg

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
