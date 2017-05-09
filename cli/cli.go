package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"path"

	log "github.com/Sirupsen/logrus"
	"github.com/bitrise-core/bitrise-plugins-analytics/analytics"
	"github.com/bitrise-core/bitrise-plugins-analytics/configs"
	"github.com/bitrise-core/bitrise-plugins-analytics/version"
	"github.com/codegangsta/cli"
)

//=======================================
// Variables
//=======================================

const (
	pluginInputPayloadKey        = "BITRISE_PLUGIN_INPUT_PAYLOAD"
	pluginInputBitriseVersionKey = "BITRISE_PLUGIN_INPUT_BITRISE_VERSION"
	pluginInputTriggerEventKey   = "BITRISE_PLUGIN_INPUT_TRIGGER"
	pluginInputPluginModeKey     = "BITRISE_PLUGIN_INPUT_PLUGIN_MODE"
	pluginInputDataDirKey        = "BITRISE_PLUGIN_INPUT_DATA_DIR"
	cIModeKey                    = "CI"
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
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:   true,
		ForceColors:     true,
		TimestampFormat: "15:04:05",
	})

	// Log level
	// If log level defined - use it
	logLevelStr := c.String("loglevel")
	if logLevelStr == "" {
		logLevelStr = "info"
	}

	level, err := log.ParseLevel(logLevelStr)
	if err != nil {
		return err
	}
	log.SetLevel(level)

	bitriseVersion := os.Getenv(pluginInputBitriseVersionKey)

	log.Debug("")
	log.Debugf("pluginInputBitriseVersion: %s", bitriseVersion)

	triggerEvent := os.Getenv(pluginInputTriggerEventKey)

	log.Debug("")
	log.Debugf("pluginInputTriggerEvent: %s", triggerEvent)

	configs.DataDir = os.Getenv(pluginInputDataDirKey)

	log.Debug("")
	log.Debugf("pluginInputDataDir: %s", configs.DataDir)

	configs.IsCIMode = (os.Getenv(cIModeKey) == "true")

	log.Debug("")
	log.Debugf("pluginInputCIModeKey: %v", configs.IsCIMode)

	return nil
}

func action(c *cli.Context) {
	pluginMode := os.Getenv(pluginInputPluginModeKey)

	log.Debug("")
	log.Debugf("pluginInputPluginMode: %s", pluginMode)

	if pluginMode == string(triggerMode) {
		config, err := configs.ReadConfig()
		if err != nil {
			log.Fatalf("Failed to read analytics configuration, error: %s", err)
		}

		if config.IsAnalyticsDisabled {
			log.Debug("")
			log.Debug("Analytics disabled")
			return
		}

		payload := os.Getenv(pluginInputPayloadKey)

		log.Debug("")
		log.Debugf("pluginInputPayload: %s", payload)

		var buildRunResults analytics.BuildRunResultsModel
		if err := json.Unmarshal([]byte(payload), &buildRunResults); err != nil {
			log.Fatalf("Failed to parse plugin input (%s), error: %s", payload, err)
		}

		log.Infof("")
		log.Infof("Submitting anonymized usage information...")
		log.Infof("For more information visit:")
		log.Infof("https://github.com/bitrise-core/bitrise-plugins-analytics/blob/master/README.md")

		if err := analytics.SendAnonymizedAnalytics(buildRunResults); err != nil {
			log.Fatalf("Failed to send analytics, error: %s", err)
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
		log.Fatalf("Failed to turn on analytics, error: %s", err)
	}
}

func analyticsOFF(c *cli.Context) {
	log.Infof("")
	log.Infof("\x1b[34;1mTurning analytics off...\x1b[0m")

	if err := configs.SetAnalytics(false); err != nil {
		log.Fatalf("Failed to turn off analytics, error: %s", err)
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
		log.Fatal("Finished with Error:", err)
	}
}
