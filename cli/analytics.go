package cli

import (
	"encoding/json"
	"fmt"
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

func sendAnalytics() {
	config, err := configs.ReadConfig()
	if err != nil {
		log.Errorf("Failed to read analytics configuration, error: %s", err)
		os.Exit(1)
	}

	if config.IsAnalyticsDisabled {
		return
	}

	hostBitriseFormatVersionStr := os.Getenv(plugins.PluginInputFormatVersionKey)
	pluginFormatVersionStr := models.Version
	if warn, err := ensureFormatVersion(pluginFormatVersionStr, hostBitriseFormatVersionStr); err != nil {
		log.Errorf(err.Error())
		os.Exit(1)
	} else if warn != "" {
		log.Warnf(warn)
	}

	payload := os.Getenv(plugins.PluginInputPayloadKey)
	var buildRunResults models.BuildRunResultsModel
	if err := json.Unmarshal([]byte(payload), &buildRunResults); err != nil {
		log.Errorf("Failed to parse plugin input (%s), error: %s", payload, err)
		os.Exit(1)
	}

	log.Infof("")
	log.Infof("Submitting anonymized usage informations...")
	log.Infof("For more information visit:")
	log.Infof("https://github.com/bitrise-io/bitrise-plugins-analytics/blob/master/README.md")

	if err := analytics.SendAnonymizedAnalytics(buildRunResults); err != nil {
		log.Errorf("Failed to send analytics, error: %s", err)
		os.Exit(1)
	}
}
