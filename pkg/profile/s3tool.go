package profile

import (
	"context"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/schidstorm/s3tool/pkg/cli"
	"gopkg.in/yaml.v3"
)

type S3ClientLoader struct{}

func (l *S3ClientLoader) LoadProfiles() ([]Connector, error) {
	files, err := os.ReadDir(cli.Config.ProfilesDirectory)
	if err != nil {
		return nil, err
	}

	var profiles []Connector
	for _, file := range files {
		var profileName = file.Name()

		if file.IsDir() {
			continue
		}

		if strings.HasSuffix(profileName, ".yaml") {
			profileName = strings.TrimSuffix(profileName, ".yaml")
		} else if strings.HasSuffix(profileName, ".yml") {
			profileName = strings.TrimSuffix(profileName, ".yml")
		} else {
			continue
		}

		fileContent, err := os.ReadFile(cli.Config.ProfilesDirectory + "/" + file.Name())
		if err != nil {
			return nil, err
		}

		var parameters ProfileParameters

		err = yaml.Unmarshal(fileContent, &parameters)
		if err != nil {
			return nil, err
		}

		profiles = append(profiles, &S3ProfileConnector{
			name:       profileName,
			parameters: parameters,
		})
	}

	return profiles, nil
}

type S3ProfileConnector struct {
	name       string
	parameters ProfileParameters
}

func (c *S3ProfileConnector) Name() string {
	return c.name
}

func (c *S3ProfileConnector) Type() string {
	return "s3tool"
}

func (c *S3ProfileConnector) CreateClient(ctx context.Context) (*s3.Client, error) {
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(c.parameters.Region),
		config.WithCredentialsProvider(aws.CredentialsProviderFunc(func(ctx context.Context) (aws.Credentials, error) {
			return aws.Credentials{
				AccessKeyID:     c.parameters.AccessKeyID,
				SecretAccessKey: c.parameters.SecretAccessKey,
				SessionToken:    c.parameters.SessionToken,
			}, nil
		})),
		config.WithBaseEndpoint(c.parameters.BaseEndpoint),
		config.WithRegion(c.parameters.Region),
	)
	if err != nil {
		return nil, err
	}

	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		if c.parameters.UsePathStyle != nil {
			o.UsePathStyle = *c.parameters.UsePathStyle
		}
	})

	return client, nil
}
