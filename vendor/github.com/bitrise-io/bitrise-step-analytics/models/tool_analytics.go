package models

type ToolAnalytics struct {
	BuildSlug string `json:"build_slug"`

	ToolUsage []ToolUsage `json:"tool_usage"`
}
