package analytics

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/bitrise-core/bitrise-plugins-analytics/configs"
)

//=======================================
// Consts
//=======================================

const analyticsURL = "https://bitrise-stats.herokuapp.com/save"

//=======================================
// Models
//=======================================

// AnonymizedUsageModel ...
type AnonymizedUsageModel struct {
	ID      string  `json:"step"`
	Version string  `json:"version"`
	RunTime float64 `json:"duration"`
	Error   bool    `json:"error"`
	CI      bool    `json:"ci"`
}

// AnonymizedUsageGroupModel ...
type AnonymizedUsageGroupModel struct {
	Steps []AnonymizedUsageModel `json:"steps"`
}

// BuildRunResultsModel ...
type BuildRunResultsModel struct {
	SuccessSteps         []StepRunResultsModel `json:"SuccessSteps"`
	FailedSteps          []StepRunResultsModel `json:"FailedSteps"`
	FailedSkippableSteps []StepRunResultsModel `json:"FailedSkippableSteps"`
	SkippedSteps         []StepRunResultsModel `json:"SkippedSteps"`
}

// StepRunResultsModel ...
type StepRunResultsModel struct {
	StepInfo struct {
		ID      string `json:"step_id"`
		Version string `json:"step_version"`
	}
	RunTime time.Duration `json:"RunTime"`
	Status  int           `json:"Status"`
	Idx     int           `json:"Idx"`
}

//=======================================
// Main
//=======================================

// SendAnonymizedAnalytics ...
func SendAnonymizedAnalytics(buildRunResults BuildRunResultsModel) error {
	stepRunResults := buildRunResults.SuccessSteps
	stepRunResults = append(stepRunResults, buildRunResults.FailedSteps...)
	stepRunResults = append(stepRunResults, buildRunResults.FailedSkippableSteps...)
	stepRunResults = append(stepRunResults, buildRunResults.SkippedSteps...)

	anonymizedUsageGroup := AnonymizedUsageGroupModel{}
	for _, stepRunResult := range stepRunResults {
		anonymizedUsageData := AnonymizedUsageModel{
			ID:      stepRunResult.StepInfo.ID,
			Version: stepRunResult.StepInfo.Version,
			RunTime: stepRunResult.RunTime.Seconds(),
			Error:   stepRunResult.Status != 0,
			CI:      configs.IsCIMode,
		}

		anonymizedUsageGroup.Steps = append(anonymizedUsageGroup.Steps, anonymizedUsageData)
	}

	data, err := json.MarshalIndent(anonymizedUsageGroup, "", "\t")
	if err != nil {
		return err
	}

	data, err = json.Marshal(anonymizedUsageGroup)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", analyticsURL, bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("failed to create request with usage data (%s), error: %s", string(data), err)
	}

	req.Header.Set("Content-Type", "application/json")

	timeout := time.Duration(10 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to perform request with usage data (%s), error: %s", string(data), err)
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Errorf("failed to close response body, error: %#v", err)
		}
	}()

	if resp.StatusCode < 200 || resp.StatusCode > 210 {
		return fmt.Errorf("sending analytics data (%s), failed with status code: %d", string(data), resp.StatusCode)
	}

	return nil
}
