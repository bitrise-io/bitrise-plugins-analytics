package analytics

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	models "github.com/bitrise-io/bitrise/models"
	"github.com/bitrise-io/bitrise/plugins"
	"github.com/lunny/log"
)

//=======================================
// Consts
//=======================================

const analyticsBaseURL = "https://bitrise-step-analytics.herokuapp.com"

//=======================================
// Models
//=======================================

// BuildAnalytics ...
type BuildAnalytics struct {
	AppSlug       string          `json:"app_slug"`
	BuildSlug     string          `json:"build_slug"`
	Status        string          `json:"status"`
	StartTime     time.Time       `json:"start_time"`
	Runtime       time.Duration   `json:"run_time"`
	StepAnalytics []StepAnalytics `json:"step_analytics"`
	CLIVersion    string          `json:"cli_version"`
	Platform      string          `json:"platform"`
	// StackID    string        `json:"stack_id"` // not supported
}

// StepAnalytics ...
type StepAnalytics struct {
	StepID    string        `json:"step_id"`
	Status    string        `json:"status"`
	Runtime   time.Duration `json:"run_time"`
	StartTime time.Time     `json:"start_time"`
}

//=======================================
// Main
//=======================================

// SendAnonymizedAnalytics ...
func SendAnonymizedAnalytics(buildRunResults models.BuildRunResultsModel) error {
	var (
		body          bytes.Buffer
		runtime       time.Duration
		stepAnalytics []StepAnalytics
		buildResults  = map[bool]string{
			true:  "failed",
			false: "success",
		}
		stepResults = func(i int) string {
			switch i {
			case models.StepRunStatusCodeSuccess:
				return "success"
			case models.StepRunStatusCodeFailed:
				return "failed"
			case models.StepRunStatusCodeFailedSkippable:
				return "failed_skippable"
			case models.StepRunStatusCodeSkipped:
				return "skipped"
			case models.StepRunStatusCodeSkippedWithRunIf:
				return "skipped_with_runif"
			default:
				return "unknown"
			}
		}
	)

	for _, stepResult := range buildRunResults.OrderedResults() {
		runtime += stepResult.RunTime
		stepAnalytics = append(stepAnalytics, StepAnalytics{
			StepID:    stepResult.StepInfo.ID,
			Status:    stepResults(stepResult.Status),
			Runtime:   stepResult.RunTime,
			StartTime: stepResult.StartTime,
		})
	}

	if err := json.NewEncoder(&body).Encode(BuildAnalytics{
		CLIVersion: os.Getenv(plugins.PluginInputBitriseVersionKey),
		Platform:   buildRunResults.ProjectType,
		AppSlug:    os.Getenv("BITRISE_APP_SLUG"),
		BuildSlug:  os.Getenv("BITRISE_BUILD_SLUG"),
		StartTime:  buildRunResults.StartTime,
		Status:     buildResults[buildRunResults.IsBuildFailed()],
		Runtime:    runtime,
	}); err != nil {
		return err
	}

	req, err := http.NewRequest("POST", analyticsBaseURL+"/metrics", &body)
	if err != nil {
		return fmt.Errorf("failed to create request with usage data (%s), error: %s", body.String(), err)
	}

	req.Header.Set("Content-Type", "application/json")

	timeout := time.Duration(10 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to perform request with usage data (%s), error: %s", body.String(), err)
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Errorf("failed to close response body, error: %#v", err)
		}
	}()

	if resp.StatusCode < 200 || resp.StatusCode > 210 {
		return fmt.Errorf("sending analytics data (%s), failed with status code: %d", body.String(), resp.StatusCode)
	}

	return nil
}
