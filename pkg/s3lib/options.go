package s3lib

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/go-ini/ini"
	"github.com/schidstorm/s3tool/pkg/cli"
	"gopkg.in/yaml.v3"
)

type ContextType string

const (
	ContextTypeAws    ContextType = "aws"
	ContextTypeS3Tool ContextType = "s3tool"
)

type Context struct {
	Type ContextType
	Name string
}

type S3ToolContexts struct {
	Contexts map[string]S3ToolContext `yaml:"contexts"`
	Default  string                   `yaml:"default,omitempty"`
}

type S3ToolContext struct {
	config.EnvConfig
	UsePathStyle bool `yaml:"use_path_style"`
}

type configProvider struct {
	LoadSdkConfig func(*config.LoadOptions) error
	LoadS3Options func(*s3.Options)
}

func ListContexts() []Context {
	var contexts []Context
	awsProfiles := listAwsProfiles()

	for _, profile := range awsProfiles {
		contexts = append(contexts, Context{
			Type: ContextTypeAws,
			Name: profile,
		})
	}

	s3ToolContexts := listS3ToolContexts()
	for _, context := range s3ToolContexts {
		contexts = append(contexts, Context{
			Type: ContextTypeS3Tool,
			Name: context,
		})
	}

	return contexts
}

func listS3ToolContexts() []string {
	contexts, err := loadS3ToolContexts()
	if err != nil {
		return nil // no contexts or error loading contexts
	}

	var contextList []string
	for name := range contexts.Contexts {
		contextList = append(contextList, name)
	}
	return contextList
}

func loadS3ToolContexts() (*S3ToolContexts, error) {
	fileContent, err := os.ReadFile(cli.Config.ContextsYamlPath)
	if err != nil {
		return nil, err // error reading contexts file
	}

	var contexts S3ToolContexts

	err = yaml.Unmarshal(fileContent, &contexts)
	if err != nil {
		return nil, err // error parsing contexts file
	}

	return &contexts, nil
}

func listAwsProfiles() []string {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil // no home == no profiles
	}

	configFile := filepath.Join(home, ".aws", "config")
	credsFile := filepath.Join(home, ".aws", "credentials")

	profiles := make(map[string]bool)

	loadProfiles := func(path string, isConfig bool) {
		cfg, err := ini.Load(path)
		if err != nil {
			return
		}

		for _, section := range cfg.Sections() {
			name := section.Name()
			if name == ini.DefaultSection {
				profiles["default"] = true
			} else if isConfig && strings.HasPrefix(name, "profile ") {
				profiles[strings.TrimPrefix(name, "profile ")] = true
			} else {
				profiles[name] = true
			}
		}
	}

	loadProfiles(configFile, true)
	loadProfiles(credsFile, false)

	var profileList []string
	for profile := range profiles {
		profileList = append(profileList, profile)
	}
	return profileList
}

func LoadContext(typeStr, name string) (*s3.Client, error) {
	var configProvider configProvider
	switch ContextType(typeStr) {
	case ContextTypeAws:
		configProvider = loadAwsOptions(name)
	case ContextTypeS3Tool:
		configProvider = loadS3ToolOptions(name)
	default:
		return nil, fmt.Errorf("unknown context type: %s", typeStr)
	}

	opts := []func(*config.LoadOptions) error{}
	if configProvider.LoadSdkConfig != nil {
		opts = append(opts, configProvider.LoadSdkConfig)
	}

	sdkConfig, err := config.LoadDefaultConfig(context.Background(), opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to load config for context %s: %w", name, err)
	}

	var s3Options []func(*s3.Options)
	if configProvider.LoadS3Options != nil {
		s3Options = append(s3Options, configProvider.LoadS3Options)
	}

	return s3.NewFromConfig(sdkConfig, s3Options...), nil
}

func loadAwsOptions(profile string) configProvider {
	var provider configProvider
	provider.LoadSdkConfig = func(o *config.LoadOptions) error {
		return config.WithSharedConfigProfile(profile)(o)
	}

	return provider
}

func loadS3ToolOptions(name string) configProvider {
	var provider configProvider
	provider.LoadSdkConfig = func(o *config.LoadOptions) error {
		contexts, err := loadS3ToolContexts()
		if err != nil {
			return err // error loading S3Tool contexts
		}

		if ctx, exists := contexts.Contexts[name]; !exists {
			return fmt.Errorf("context %s not found", name) // context not found
		} else {
			o.Region = ctx.Region
			o.BaseEndpoint = ctx.BaseEndpoint
			o.Credentials = aws.CredentialsProviderFunc(func(_ context.Context) (aws.Credentials, error) {
				return aws.Credentials{
					AccessKeyID:     ctx.Credentials.AccessKeyID,
					SecretAccessKey: ctx.Credentials.SecretAccessKey,
					SessionToken:    ctx.Credentials.SessionToken,
				}, nil
			})

		}

		return nil
	}

	provider.LoadS3Options = func(cfg *s3.Options) {
		contexts, err := loadS3ToolContexts()
		if err != nil {
			return
		}

		if context, exists := contexts.Contexts[name]; !exists {
			return
		} else {
			cfg.UsePathStyle = context.UsePathStyle
		}
	}

	return provider
}
