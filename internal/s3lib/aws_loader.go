package s3lib

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-ini/ini"
)

type AwsLoader struct{}

func (l *AwsLoader) Load() ([]Connector, error) {
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
		profileList = append(profileList, &AwsConnector{name: profile})
	}
	return profileList, nil
}
