package analytics

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	analyticsModels "github.com/bitrise-io/bitrise-step-analytics/models"
	models "github.com/bitrise-io/bitrise/models"
	"github.com/bitrise-io/bitrise/plugins"
	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-io/go-utils/pointers"
)

//=======================================
// Consts
//=======================================

const (
	stackIDEnvKey    = "BITRISEIO_STACK_ID"
	appSlugEnvKey    = "BITRISE_APP_SLUG"
	buildSlugEnvKey  = "BITRISE_BUILD_SLUG"
	workflowName     = "BITRISE_TRIGGERED_WORKFLOW_TITLE"
	repoSlug         = "BITRISEIO_GIT_REPOSITORY_SLUG"
	analyticsBaseURL = "https://bitrise-step-analytics.herokuapp.com"
)

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
		stepAnalytics []analyticsModels.StepAnalytics
	)

	stepInputWhitelist := map[string]map[string]bool{
		"xcode-test": {
			"simulator_device":         true,
			"simulator_os_version":     true,
			"single_build":             true,
			"should_build_before_test": true,
		},
		"ios-auto-provision-appstoreconnect": {
			"connection":             true,
			"min_profile_days_valid": true,
		},
		"ios-auto-provision": {
			"min_profile_days_valid": true,
		},
	}

	for _, stepResult := range buildRunResults.OrderedResults() {
		filteredStepInputs := map[string]string{}
		for key, value := range stepResult.StepInputs {
			if stepInputWhitelist[stepResult.StepInfo.ID][key] {
				filteredStepInputs[key] = value
			}
		}

		stepAnalytics, runtime = append(stepAnalytics, analyticsModels.StepAnalytics{
			StepID:      stepResult.StepInfo.ID,
			StepTitle:   pointers.StringWithDefault(stepResult.StepInfo.Step.Title, ""),
			StepVersion: stepResult.StepInfo.Version,
			StepSource:  pointers.StringWithDefault(stepResult.StepInfo.Step.SourceCodeURL, ""),
			StepInputs:  filteredStepInputs,
			Status:      stepStatus(stepResult.Status),
			Runtime:     stepResult.RunTime,
			StartTime:   stepResult.StartTime,
		}), runtime+stepResult.RunTime
	}

	var body bytes.Buffer
	if err := json.NewEncoder(&body).Encode(analyticsModels.BuildAnalytics{
		Runtime:       runtime,
		StartTime:     buildRunResults.StartTime,
		Platform:      buildRunResults.ProjectType,
		StackID:       os.Getenv(stackIDEnvKey),
		AppSlug:       os.Getenv(appSlugEnvKey),
		BuildSlug:     os.Getenv(buildSlugEnvKey),
		Status:        buildStatus(buildRunResults.IsBuildFailed()),
		CLIVersion:    os.Getenv(plugins.PluginConfigBitriseVersionKey),
		StepAnalytics: stepAnalytics,
		RepositoryID:  os.Getenv(repoSlug),
		WorkflowName:  os.Getenv(workflowName),
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
