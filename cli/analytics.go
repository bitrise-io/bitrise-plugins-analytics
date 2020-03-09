package cli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/bitrise-io/bitrise-plugins-analytics/analytics"
	"github.com/bitrise-io/bitrise-plugins-analytics/configs"
	"github.com/bitrise-io/bitrise/models"
	"github.com/bitrise-io/bitrise/plugins"
	log "github.com/bitrise-io/go-utils/log"
	"github.com/hashicorp/go-version"
	"github.com/pkg/errors"
)

// minBitriseCLIVersion points to the version of Bitrise CLI introduceses Bitrise plugins.
const minBitriseCLIVersion = "1.6.0"

func ensureFormatVersion(pluginFormatVersionStr, hostBitriseFormatVersionStr string) (string, error) {
	if hostBitriseFormatVersionStr == "" {
		return fmt.Sprintf("This analytics plugin version would need bitrise-cli version >= %s to submit analytics", minBitriseCLIVersion), nil
	}

	hostBitriseFormatVersion, err := version.NewVersion(hostBitriseFormatVersionStr)
	if err != nil {
		return "", errors.Wrapf(err, "failed to parse bitrise format version (%s)", hostBitriseFormatVersionStr)
	}

	pluginFormatVersion, err := version.NewVersion(pluginFormatVersionStr)
	if err != nil {
		return "", errors.Errorf("failed to parse analytics plugin format version (%s), error: %s", pluginFormatVersionStr, err)
	}

	if pluginFormatVersion.LessThan(hostBitriseFormatVersion) {
		return "Outdated analytics plugin, used format version is lower then host bitrise-cli's format version, please update the plugin", nil
	} else if pluginFormatVersion.GreaterThan(hostBitriseFormatVersion) {
		return "Outdated bitrise-cli, used format version is lower then the analytics plugin's format version, please update the bitrise-cli", nil
	}

	return "", nil
}

func isAnalyticsEnabled() (bool, error) {
	config, err := configs.ReadConfig()
	if err != nil {
		return false, fmt.Errorf("failed to read analytics configuration: %s", err)
	}
	return !config.IsAnalyticsDisabled, nil
}

func ensureBitriseCLIVersion() (string, error) {
	hostBitriseFormatVersionStr := os.Getenv(plugins.PluginInputFormatVersionKey)
	pluginFormatVersionStr := models.Version
	return ensureFormatVersion(pluginFormatVersionStr, hostBitriseFormatVersionStr)
}

func read(r io.Reader) ([]byte, error) {
	var buff []byte
	for {
		chunk := make([]byte, 100)
		_, err := r.Read(chunk)
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		buff = append(buff, chunk...)
	}
	return buff, nil
}

var errNoInput = errors.New("nothing to read")

func readPayload(r io.Reader) (models.BuildRunResultsModel, error) {
	b, err := read(r)
	if err != nil {
		return models.BuildRunResultsModel{}, err
	}
	// Remove any NULL characters from 'b'
	b = bytes.Trim(b, "\x00")
	if len(b) == 0 {
		return models.BuildRunResultsModel{}, errNoInput
	}

	var buildRunResults models.BuildRunResultsModel
	if err := json.Unmarshal(b, &buildRunResults); err != nil {
		return models.BuildRunResultsModel{}, fmt.Errorf("failed to parse plugin input (%s): %s", string(b), err)
	}
	return buildRunResults, nil
}

func sendAnalyticsIfEnabled() {
	if enabled, err := isAnalyticsEnabled(); err != nil {
		failf(err.Error())
	} else if !enabled {
		return
	}

	if warn, err := ensureBitriseCLIVersion(); err != nil {
		failf(err.Error())
	} else if len(warn) > 0 {
		log.Warnf(warn)
	}

	payload, err := readPayload(os.Stdin)
	if err != nil {
		failf(err.Error())
	}

	log.Infof("")
	log.Infof("Submitting anonymized usage informations...")
	log.Infof("For more information visit:")
	log.Infof("https://github.com/bitrise-io/bitrise-plugins-analytics/blob/master/README.md")

	if err := analytics.SendAnonymizedAnalytics(payload); err != nil {
		failf("Failed to send analytics, error: %s", err)
	}
}
