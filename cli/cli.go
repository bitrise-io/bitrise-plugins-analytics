package cli

import (
	"fmt"
	"os"
	"path"

	"github.com/bitrise-io/bitrise-plugins-analytics/configs"
	"github.com/bitrise-io/bitrise-plugins-analytics/version"
	bitriseConfigs "github.com/bitrise-io/bitrise/configs"
	"github.com/bitrise-io/bitrise/plugins"
	log "github.com/bitrise-io/go-utils/log"
	"github.com/urfave/cli"
)

var commands = []cli.Command{
	createSwitchCommand(true),
	createSwitchCommand(false),
}

var flags = []cli.Flag{
	cli.StringFlag{
		Name:   "loglevel, l",
		Usage:  "Log level (options: debug, info, warn, error, fatal, panic).",
		EnvVar: "LOGLEVEL",
	},
}

func before(c *cli.Context) error {
	configs.DataDir = os.Getenv(plugins.PluginInputDataDirKey)
	configs.IsCIMode = (os.Getenv(bitriseConfigs.CIModeEnvKey) == "true")
	return nil
}

func printVersion(c *cli.Context) {
	fmt.Println(c.App.Version)
}

func action(c *cli.Context) {
	if os.Getenv(plugins.PluginInputPluginModeKey) == string(plugins.TriggerMode) {
		sendAnalytics()
		return
	}

	if err := cli.ShowAppHelp(c); err != nil {
		log.Errorf("Failed to show help, error: %s", err)
		os.Exit(1)
	}
}

func createApp() *cli.App {
	app := cli.NewApp()

	app.Name = path.Base(os.Args[0])
	app.Usage = "Bitrise Analytics plugin"
	app.Version = version.VERSION

	app.Author = ""
	app.Email = ""

	app.Before = before
	app.Flags = flags
	app.Commands = commands
	app.Action = action

	return app
}

// Run ...
func Run() {
	cli.VersionPrinter = printVersion

	if err := createApp().Run(os.Args); err != nil {
		log.Errorf("Finished with Error: %s", err)
		os.Exit(1)
	}
}
