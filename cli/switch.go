package cli

import (
	"fmt"
	"os"

	"github.com/bitrise-io/bitrise-plugins-analytics/configs"
	log "github.com/bitrise-io/go-utils/log"
	"github.com/urfave/cli"
)

type state bool

func (s state) String() string {
	if s {
		return "on"
	}
	return "off"
}

func createSwitchCommand(s state) cli.Command {
	return cli.Command{
		Name:  s.String(),
		Usage: fmt.Sprintf("Turn sending anonimized usage information %s.", s),
		Action: func(c *cli.Context) {
			log.Infof("")
			log.Infof("Turning analytics %s...", s)

			if err := configs.SetAnalytics(bool(s)); err != nil {
				log.Errorf("Failed to turn %s analytics, error: %s", s, err)
				os.Exit(1)
			}
		},
	}
}
