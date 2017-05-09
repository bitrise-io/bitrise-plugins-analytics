package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"path"

	"github.com/bitrise-core/bitrise-plugins-analytics/analytics"
	"github.com/bitrise-core/bitrise-plugins-analytics/configs"
	"github.com/bitrise-core/bitrise-plugins-analytics/version"
	models "github.com/bitrise-io/bitrise/models"
	log "github.com/bitrise-io/go-utils/log"
	"github.com/codegangsta/cli"
	ver "github.com/hashicorp/go-version"
)

//=======================================
// Variables
//=======================================

const (
	pluginInputPayloadKey       = "BITRISE_PLUGIN_INPUT_PAYLOAD"
	pluginInputPluginModeKey    = "BITRISE_PLUGIN_INPUT_PLUGIN_MODE"
	pluginInputDataDirKey       = "BITRISE_PLUGIN_INPUT_DATA_DIR"
	pluginInputFormatVersionKey = "BITRISE_PLUGIN_INPUT_FORMAT_VERSION"
	cIModeKey                   = "CI"
)

const (
	triggerMode PluginMode = "trigger"
)

// PluginMode ...
type PluginMode string

var commands = []cli.Command{
	cli.Command{
		Name:   "on",
		Usage:  "Turn sending anonimized usage information on.",
		Action: analyticsON,
	},
	cli.Command{
		Name:   "off",
		Usage:  "Turn sending anonimized usage information off.",
		Action: analyticsOFF,
	},
}

//=======================================
// Functions
//=======================================

func printVersion(c *cli.Context) {
	fmt.Fprintf(c.App.Writer, "%v\n", c.App.Version)
}

func before(c *cli.Context) error {
	configs.DataDir = os.Getenv(pluginInputDataDirKey)
	configs.IsCIMode = (os.Getenv(cIModeKey) == "true")

	return nil
}

func action(c *cli.Context) {
	pluginMode := os.Getenv(pluginInputPluginModeKey)
	if pluginMode == string(triggerMode) {
		// ensure plugin's format version matches to host bitrise-cli's format version
		hostBitriseFormatVersionStr := os.Getenv(pluginInputFormatVersionKey)
		pluginBitriseFormatVersionStr := models.Version

		if hostBitriseFormatVersionStr == "" {
			log.Warnf("This analytics plugin version would need bitrise-cli version >= 1.6.0 to submit analytics")
			return
		}

		hostBitriseFormatVersion, err := ver.NewVersion(hostBitriseFormatVersionStr)
		if err != nil {
			log.Errorf("Failed to parse bitrise format version (%s), error: %s", hostBitriseFormatVersionStr, err)
			os.Exit(1)
		}

		pluginBitriseFormatVersion, err := ver.NewVersion(pluginBitriseFormatVersionStr)
		if err != nil {
			log.Errorf("Failed to parse analytics plugin format version (%s), error: %s", pluginBitriseFormatVersionStr, err)
			os.Exit(1)
		}

		if pluginBitriseFormatVersion.LessThan(hostBitriseFormatVersion) {
			log.Warnf("Outdated analytics plugin, used format version is lower then host bitrise-cli's format version, please update the plugin")
			return
		} else if pluginBitriseFormatVersion.GreaterThan(hostBitriseFormatVersion) {
			log.Warnf("Outdated bitrise-cli, used format version is lower then the analytics plugin's format version, please update the bitrise-cli")
			return
		}
		// ---

		config, err := configs.ReadConfig()
		if err != nil {
			log.Errorf("Failed to read analytics configuration, error: %s", err)
			os.Exit(1)
		}

		if config.IsAnalyticsDisabled {
			return
		}

		payload := os.Getenv(pluginInputPayloadKey)
		var buildRunResults models.BuildRunResultsModel
		if err := json.Unmarshal([]byte(payload), &buildRunResults); err != nil {
			log.Errorf("Failed to parse plugin input (%s), error: %s", payload, err)
			os.Exit(1)
		}

		log.Infof("")
		log.Infof("Submitting anonymized usage informations...")
		log.Infof("For more information visit:")
		log.Infof("https://github.com/bitrise-core/bitrise-plugins-analytics/blob/master/README.md")

		if err := analytics.SendAnonymizedAnalytics(buildRunResults); err != nil {
			log.Errorf("Failed to send analytics, error: %s", err)
			os.Exit(1)
		}
	} else {
		if err := cli.ShowAppHelp(c); err != nil {
			log.Errorf("Failed to show help, error: %s", err)
			os.Exit(1)
		}
	}
}

func analyticsON(c *cli.Context) {
	log.Infof("")
	log.Infof("\x1b[34;1mTurning analytics on...\x1b[0m")

	if err := configs.SetAnalytics(true); err != nil {
		log.Errorf("Failed to turn on analytics, error: %s", err)
		os.Exit(1)
	}
}

func analyticsOFF(c *cli.Context) {
	log.Infof("")
	log.Infof("\x1b[34;1mTurning analytics off...\x1b[0m")

	if err := configs.SetAnalytics(false); err != nil {
		log.Errorf("Failed to turn off analytics, error: %s", err)
		os.Exit(1)
	}
}

//=======================================
// Main
//=======================================

// Run ...
func Run() {
	// Parse cl
	cli.VersionPrinter = printVersion

	app := cli.NewApp()

	app.Name = path.Base(os.Args[0])
	app.Usage = "Bitrise Analytics plugin"
	app.Version = version.VERSION

	app.Author = ""
	app.Email = ""

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "loglevel, l",
			Usage:  "Log level (options: debug, info, warn, error, fatal, panic).",
			EnvVar: "LOGLEVEL",
		},
	}
	app.Before = before
	app.Commands = commands
	app.Action = action

	if err := app.Run(os.Args); err != nil {
		log.Errorf("Finished with Error: %s", err)
		os.Exit(1)
	}
}
