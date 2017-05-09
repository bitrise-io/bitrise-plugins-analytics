package configs

import (
	"errors"
	"path"

	"gopkg.in/yaml.v2"

	"github.com/bitrise-io/go-utils/fileutil"
	"github.com/bitrise-io/go-utils/pathutil"
)

//=======================================
// Variables
//=======================================

var (
	// DataDir ...
	DataDir string
	// IsCIMode ...
	IsCIMode bool
)

//=======================================
// Models
//=======================================

// ConfigModel ...
type ConfigModel struct {
	IsAnalyticsDisabled bool `yaml:"is_analytics_disabled"`
}

// NewConfigFromBytes ...
func NewConfigFromBytes(bytes []byte) (ConfigModel, error) {
	var config ConfigModel
	if err := yaml.Unmarshal(bytes, &config); err != nil {
		return ConfigModel{}, err
	}

	return config, nil
}

//=======================================
// Main
//=======================================

func getConfigFilePath() (string, error) {
	if DataDir == "" {
		return "", errors.New("could not determin plugin data dir, run plugin through bitrise (bitrise :analytics COMMAND_TO_RUN)")
	}
	return path.Join(DataDir, "config.yml"), nil
}

// ReadConfig ...
func ReadConfig() (ConfigModel, error) {
	config := ConfigModel{}

	configPth, err := getConfigFilePath()
	if err != nil {
		return config, err
	}

	if exist, err := pathutil.IsPathExists(configPth); err != nil {
		return ConfigModel{}, err
	} else if exist {
		bytes, err := fileutil.ReadBytesFromFile(configPth)
		if err != nil {
			return ConfigModel{}, err
		}

		if config, err = NewConfigFromBytes(bytes); err != nil {
			return ConfigModel{}, err
		}
	}

	return config, nil
}

func saveConfig(config ConfigModel) error {
	bytes, err := yaml.Marshal(config)
	if err != nil {
		return err
	}

	configPth, err := getConfigFilePath()
	if err != nil {
		return err
	}

	return fileutil.WriteBytesToFileWithPermission(configPth, bytes, 0)
}

// SetAnalytics ...
func SetAnalytics(enabled bool) error {
	config, err := ReadConfig()
	if err != nil {
		return err
	}

	config.IsAnalyticsDisabled = !enabled

	return saveConfig(config)
}
