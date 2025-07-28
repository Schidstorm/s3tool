package cli

import (
	"flag"
	"os"
	"path"
	"strings"
)

var Config = S3ToolCliConfig{}

type S3ToolCliConfig struct {
	ContextsYamlPath string
}

func Parse(args []string) error {
	var cfg S3ToolCliConfig
	flag.StringVar(&cfg.ContextsYamlPath, "contexts-yaml", "~/.s3tool/contexts.yaml", "Path to the contexts YAML file")

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
	if strings.HasPrefix(result.ContextsYamlPath, "~") {
		home, err := os.UserHomeDir()
		if err == nil {
			result.ContextsYamlPath = strings.TrimPrefix(result.ContextsYamlPath, "~")
			result.ContextsYamlPath = path.Join(home, result.ContextsYamlPath)
		}
	}

	return result
}
