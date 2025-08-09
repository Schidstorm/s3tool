package profile

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/go-ini/ini"
)

type AwsLoader struct{}

func (l *AwsLoader) LoadProfiles() ([]Connector, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get user home directory: %w", err)
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

	var profileList []Connector
	for profile := range profiles {
		profileList = append(profileList, &AwsProfileConnector{name: profile})
	}
	return profileList, nil
}

type AwsProfileConnector struct {
	name string
}

func (c *AwsProfileConnector) Name() string {
	return c.name
}

func (c *AwsProfileConnector) Type() string {
	return "aws"
}

func (c *AwsProfileConnector) CreateClient(ctx context.Context) (*s3.Client, error) {
	opts := []func(*config.LoadOptions) error{
		config.WithSharedConfigProfile(c.name),
	}

	sdkConfig, err := config.LoadDefaultConfig(ctx, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS SDK config for profile %s: %w", c.name, err)
	}

	return s3.NewFromConfig(sdkConfig), nil
}
