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
	"github.com/bitrise-io/go-utils/log"
)

//=======================================
// Consts
//=======================================

const (
	stackIDEnvKey    = "BITRISE_STACK_ID"
	appSlugEnvKey    = "BITRISE_APP_SLUG"
	buildSlugEnvKey  = "BITRISE_BUILD_SLUG"
	analyticsBaseURL = "https://bitrise-step-analytics.herokuapp.com"
)

//=======================================
// Models
//=======================================

// BuildAnalytics ...
type BuildAnalytics struct {
	Status        string          `json:"status"`
	StackID       string          `json:"stack_id"`
	AppSlug       string          `json:"app_slug"`
	Runtime       time.Duration   `json:"run_time"`
	Platform      string          `json:"platform"`
	BuildSlug     string          `json:"build_slug"`
	StartTime     time.Time       `json:"start_time"`
	CLIVersion    string          `json:"cli_version"`
	StepAnalytics []StepAnalytics `json:"step_analytics"`
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

func buildStatus(buildFailed bool) string {
	return map[bool]string{false: "successful", true: "failed"}[buildFailed]
}

func stepStatus(i int) string {
	if status, ok := map[int]string{
		models.StepRunStatusCodeFailed:           "failed",
		models.StepRunStatusCodeSuccess:          "success",
		models.StepRunStatusCodeSkipped:          "skipped",
		models.StepRunStatusCodeFailedSkippable:  "failed_skippable",
		models.StepRunStatusCodeSkippedWithRunIf: "skipped_with_runif",
	}[i]; ok {
		return status
	}
	return "unknown"
}

// SendAnonymizedAnalytics ...
func SendAnonymizedAnalytics(buildRunResults models.BuildRunResultsModel) error {
	var (
		runtime       time.Duration
		stepAnalytics []StepAnalytics
	)
	for _, stepResult := range buildRunResults.OrderedResults() {
		stepAnalytics, runtime = append(stepAnalytics, StepAnalytics{
			StepID:    stepResult.StepInfo.ID,
			Status:    stepStatus(stepResult.Status),
			Runtime:   stepResult.RunTime,
			StartTime: stepResult.StartTime,
		}), runtime+stepResult.RunTime
	}

	var body bytes.Buffer
	if err := json.NewEncoder(&body).Encode(BuildAnalytics{
		Runtime:    runtime,
		StartTime:  buildRunResults.StartTime,
		Platform:   buildRunResults.ProjectType,
		StackID:    os.Getenv(stackIDEnvKey),
		AppSlug:    os.Getenv(appSlugEnvKey),
		BuildSlug:  os.Getenv(buildSlugEnvKey),
		Status:     buildStatus(buildRunResults.IsBuildFailed()),
		CLIVersion: os.Getenv(plugins.PluginInputBitriseVersionKey),
	}); err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, analyticsBaseURL+"/metrics", &body)
	if err != nil {
		return fmt.Errorf("failed to create request with usage data (%s), error: %s", body.String(), err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := http.Client{
		Timeout: time.Duration(10 * time.Second),
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
