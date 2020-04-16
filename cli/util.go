package cli

import (
	"fmt"
	"os"

	"github.com/bitrise-io/bitrise-plugins-analytics/configs"
	ver "github.com/hashicorp/go-version"
	"github.com/pkg/errors"
)

func isAnalyticsEnabled() (bool, error) {
	config, err := configs.ReadConfig()
	if err != nil {
		return false, fmt.Errorf("failed to read analytics configuration: %s", err)
	}
	return !config.IsAnalyticsDisabled, nil
}

// minBitriseCLIVersion points to the version of Bitrise CLI introduceses Bitrise plugins.
const minBitriseCLIVersion = "1.6.0"

func checkFormatVersion(pluginFormatVersionStr, hostBitriseFormatVersionStr string) (string, error) {
	if hostBitriseFormatVersionStr == "" {
		return fmt.Sprintf("This analytics plugin version would need bitrise-cli version >= %s to submit analytics", minBitriseCLIVersion), nil
	}

	hostBitriseFormatVersion, err := ver.NewVersion(hostBitriseFormatVersionStr)
	if err != nil {
		return "", errors.Wrapf(err, "failed to parse bitrise format version (%s)", hostBitriseFormatVersionStr)
	}

	pluginFormatVersion, err := ver.NewVersion(pluginFormatVersionStr)
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

// HasStat ...
type HasStat interface {
	Stat() (os.FileInfo, error)
}

func hasContent(f HasStat) (bool, error) {
	fi, err := f.Stat()
	if err != nil {
		return false, err
	}

	if fi.Mode()&os.ModeNamedPipe != 0 {
		return true, nil
	}

	return fi.Size() > 0, nil
}
