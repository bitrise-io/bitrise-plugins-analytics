package cli

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/bitrise-io/bitrise-plugins-analytics/analytics"
	"github.com/bitrise-io/bitrise/models"
	"github.com/bitrise-io/bitrise/plugins"
	log "github.com/bitrise-io/go-utils/log"
	"github.com/pkg/errors"
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
	return EnvPayloadSource{os.Getenv(plugins.PluginInputPayloadKey)}
}

// StdinPayloadSource ....
type StdinPayloadSource struct {
	reader io.Reader
}

// Payload ....
func (s StdinPayloadSource) Payload() (models.BuildRunResultsModel, error) {
	b, err := read(s.reader)
	if err != nil {
		return models.BuildRunResultsModel{}, err
	}
	if len(b) == 0 {
		return models.BuildRunResultsModel{}, errNoInput
	}

	var buildRunResults models.BuildRunResultsModel
	if err := json.Unmarshal(b, &buildRunResults); err != nil {
		return models.BuildRunResultsModel{}, fmt.Errorf("failed to parse plugin input (%s): %s", string(b), err)
	}
	return buildRunResults, nil
}

// EnvPayloadSource ...
type EnvPayloadSource struct {
	envValue string
}

// Payload ...
func (s EnvPayloadSource) Payload() (models.BuildRunResultsModel, error) {
	var payload models.BuildRunResultsModel
	if err := json.Unmarshal([]byte(s.envValue), &payload); err != nil {
		return models.BuildRunResultsModel{}, err
	}
	return payload, nil
}

func read(r io.Reader) ([]byte, error) {
	var buff []byte
	for {
		chunk := make([]byte, 100)
		n, err := r.Read(chunk)
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		if n == 0 {
			break
		}
		buff = append(buff, chunk[:n]...)
	}
	return buff, nil
}

var errNoInput = errors.New("nothing to read")

func readPayload(r io.Reader) (models.BuildRunResultsModel, error) {
	b, err := read(r)
	if err != nil {
		return models.BuildRunResultsModel{}, err
	}
	if len(b) == 0 {
		return models.BuildRunResultsModel{}, errNoInput
	}

	var buildRunResults models.BuildRunResultsModel
	if err := json.Unmarshal(b, &buildRunResults); err != nil {
		return models.BuildRunResultsModel{}, fmt.Errorf("failed to parse plugin input (%s): %s", string(b), err)
	}
	return buildRunResults, nil
}

func sendAnalytics(source PayloadSource) error {
	payload, err := source.Payload()
	if err != nil {
		return fmt.Errorf("failed to read payload: %s", err)
	}

	log.Infof("")
	log.Infof("Submitting anonymized usage informations...")
	log.Infof("For more information visit:")
	log.Infof("https://github.com/bitrise-io/bitrise-plugins-analytics/blob/master/README.md")

	return analytics.SendAnonymizedAnalytics(payload)
}
