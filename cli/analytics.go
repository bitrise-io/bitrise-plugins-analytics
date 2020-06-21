package cli

import (
	"fmt"
	"os"

	"github.com/bitrise-io/bitrise-plugins-analytics/analytics"
	"github.com/bitrise-io/bitrise-plugins-analytics/configs"
	"github.com/bitrise-io/bitrise/models"
	log "github.com/bitrise-io/go-utils/log"
)

// PayloadSource ...
type PayloadSource interface {
	Payload() (models.BuildRunResultsModel, error)
}

// SourceType ...
type SourceType int

// SourceTypes ...
const (
	StdinSource SourceType = iota
	EnvSource
)

// PayloadSourceFactory ....
func PayloadSourceFactory(t SourceType) PayloadSource {
	if t == StdinSource {
		return StdinPayloadSource{os.Stdin}
	}
	return EnvPayloadSource{os.Getenv(configs.PluginConfigPayloadKey)}
}

func sendAnalytics(source PayloadSource) error {
	payload, err := source.Payload()
	if err != nil {
		return fmt.Errorf("failed to read payload: %s", err)
	}

	log.Infof("")
	log.Infof("Submitting anonymized usage information...")
	log.Infof("For more information visit:")
	log.Infof("https://github.com/bitrise-io/bitrise-plugins-analytics/blob/master/README.md")

	return analytics.SendAnonymizedAnalytics(payload)
}
