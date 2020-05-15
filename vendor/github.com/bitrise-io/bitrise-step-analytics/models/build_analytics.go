package models

import (
	"time"
)

// BuildAnalytics ...
type BuildAnalytics struct {
	AppSlug      string `json:"app_slug" track:"app_slug"`
	BuildSlug    string `json:"build_slug" track:"build_slug"`
	RepositoryID string `json:"repo_id" track:"repository_id"`

	StackID      string `json:"stack_id" track:"stack_id"`
	Platform     string `json:"platform" track:"platform"`
	CLIVersion   string `json:"cli_version" track:"cli_version"`
	WorkflowName string `json:"workflow_name" track:"workflow_name"`

	Status    string        `json:"status" track:"status"`
	Runtime   time.Duration `json:"run_time" track:"run_time"`
	StartTime time.Time     `json:"start_time" track:"start_time"`

	StepAnalytics []StepAnalytics `json:"step_analytics"`
}

// Event ...
func (a BuildAnalytics) Event() string {
	return "build_finished"
}

// Model ...
func (a BuildAnalytics) Model() interface{} {
	return a
}

// UserID ...
func (a BuildAnalytics) UserID() string {
	return a.AppSlug
}
