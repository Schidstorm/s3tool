package s3lib

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3ToolConnectorParameters struct {
	AccessKeyID     string `yaml:"access_key_id"`
	SecretAccessKey string `yaml:"secret_access_key"`
	SessionToken    string `yaml:"session_token,omitempty"`
	Region          string `yaml:"region"`
	BaseEndpoint    string `yaml:"base_endpoint,omitempty"`
	UsePathStyle    *bool  `yaml:"use_path_style,omitempty"`
}

type S3ProfileConnector struct {
	name       string
	parameters S3ToolConnectorParameters
}

func (c *S3ProfileConnector) Name() string {
	return c.name
}

func (c *S3ProfileConnector) Type() string {
	return "s3tool"
}

func (c *S3ProfileConnector) CreateClient(ctx context.Context) (Client, error) {
	cfg, err := config.LoadDefaultConfig(ctx,
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

	return &SdkClient{Client: client}, nil
}
