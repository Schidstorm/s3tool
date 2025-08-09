package profile

type ProfileType string

type ProfileParameters struct {
	AccessKeyID     string `yaml:"access_key_id"`
	SecretAccessKey string `yaml:"secret_access_key"`
	SessionToken    string `yaml:"session_token,omitempty"`
	Region          string `yaml:"region"`
	BaseEndpoint    string `yaml:"base_endpoint,omitempty"`
	UsePathStyle    *bool  `yaml:"use_path_style,omitempty"`
}

func List() []Connector {
	var profiles []Connector

	awsLoader := &AwsLoader{}
	awsProfiles, err := awsLoader.LoadProfiles()
	if err == nil {
		profiles = append(profiles, awsProfiles...)
	}

	s3Loader := &S3ClientLoader{}
	s3Profiles, err := s3Loader.LoadProfiles()
	if err == nil {
		profiles = append(profiles, s3Profiles...)
	}

	return profiles
}
