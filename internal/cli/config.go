package cli

import (
	"os"
	"strings"

	"github.com/spf13/cobra"
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
	cmd := rootCmd(&cfg)
	runRoot := false
	cmd.Run = func(cmd *cobra.Command, args []string) {
		runRoot = true
	}

	cmd.AddCommand(completionCmd())

	cmd.SetArgs(args)
	if err := cmd.Execute(); err != nil {
		return err
	}

	if !runRoot {
		os.Exit(0)
		return nil
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

func rootCmd(cfg *S3ToolCliConfig) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "s3tool",
		Short: "s3tool is a terminal based S3 client",
	}

	flag := cmd.Flags()
	flag.StringVarP(&cfg.ProfilesDirectory, "profiles", "p", Config.ProfilesDirectory, "Path to a directory containing profile yaml files")
	flag.BoolVar(&cfg.Loaders.Aws, "loaders.aws", Config.Loaders.Aws, "Enable AWS loader")
	flag.BoolVar(&cfg.Loaders.S3Tool, "loaders.s3tool", Config.Loaders.S3Tool, "Enable S3Tool loader")
	flag.BoolVar(&cfg.Loaders.Memory, "loaders.memory", Config.Loaders.Memory, "Enable Memory loader (for testing purposes)")
	flag.MarkHidden("loaders.memory")

	return cmd
}

func completionCmd() *cobra.Command {
	completionCmd := &cobra.Command{
		Use:   "completion",
		Short: "Generate completion script",
		RunE: func(cmd *cobra.Command, args []string) error {
			shell, err := cmd.Flags().GetString("shell")
			if err != nil {
				return err
			}

			var err2 error
			switch strings.ToLower(shell) {
			case "bash":
				err2 = cmd.Root().GenBashCompletion(os.Stdout)
			case "zsh":
				err2 = cmd.Root().GenZshCompletion(os.Stdout)
			case "fish":
				err2 = cmd.Root().GenFishCompletion(os.Stdout, true)
			case "powershell":
				err2 = cmd.Root().GenPowerShellCompletionWithDesc(os.Stdout)
			default:
				return nil
			}

			return err2
		},
	}
	completionFlags := completionCmd.Flags()
	completionFlags.StringP("shell", "s", "bash", "Shell type (bash|zsh|fish|powershell)")
	return completionCmd
}
