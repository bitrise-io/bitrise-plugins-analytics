package cli

import (
	"fmt"
	"os"

	"github.com/bitrise-io/bitrise-plugins-analytics/configs"
	log "github.com/bitrise-io/go-utils/log"
	"github.com/urfave/cli"
)

type onOff bool

func (o onOff) String() string {
	if o == true {
		return "on"
	}
	return "off"
}

func createSwitchCommand(on onOff) cli.Command {
	return cli.Command{
		Name:  on.String(),
		Usage: fmt.Sprintf("Turn sending anonimized usage information %s.", on),
		Action: func(c *cli.Context) {
			log.Infof("")
			log.Infof("Turning analytics %s...", on)

			if err := configs.SetAnalytics(bool(on)); err != nil {
				log.Errorf("Failed to turn %s analytics, error: %s", on, err)
				os.Exit(1)
			}
		},
	}
}
