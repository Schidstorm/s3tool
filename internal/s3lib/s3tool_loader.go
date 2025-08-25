package s3lib

import (
	"os"
	"strings"

	"github.com/schidstorm/s3tool/internal/cli"
	"gopkg.in/yaml.v3"
)

type S3ToolLoader struct{}

func (l *S3ToolLoader) Load() ([]Connector, error) {
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

		var parameters S3ToolConnectorParameters

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
