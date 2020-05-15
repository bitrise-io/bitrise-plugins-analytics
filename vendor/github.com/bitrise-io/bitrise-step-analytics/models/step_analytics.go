package models

import (
	"time"
)

// StepAnalytics ...
type StepAnalytics struct {
	// connect step analytics to Bitrise builds
	AppSlug   string `json:"-"`
	BuildSlug string `json:"-" track:"build_slug"`

	StepID      string            `json:"step_id" track:"step_id"`
	StepTitle   string            `json:"step_title" track:"step_title"`
	StepVersion string            `json:"step_version" track:"step_version"`
	StepSource  string            `json:"step_source" track:"step_source"`
	StepInputs  map[string]string `json:"step_inputs" track:"step_inputs"`

	Status    string        `json:"status" track:"status"`
	StartTime time.Time     `json:"start_time" track:"start_time"`
	Runtime   time.Duration `json:"run_time" track:"run_time"`
}

// Event ...
func (a StepAnalytics) Event() string {
	return "step_finished"
}

// Model ...
func (a StepAnalytics) Model() interface{} {
	return a
}

// UserID ...
func (a StepAnalytics) UserID() string {
	return a.AppSlug
}
