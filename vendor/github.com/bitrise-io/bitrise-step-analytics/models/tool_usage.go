package models

// ToolUsage ...
type ToolUsage struct {
	BuildSlug    string `json:"-" track:"build_slug"`
	Name         string `json:"name" track:"name"`
	Version      string `json:"version" track:"version"`
	FreshInstall bool   `json:"fresh_install" track:"fresh_install"`
}

// Event ...
func (a ToolUsage) Event() string {
	return "tool_used"
}

// Model ...
func (a ToolUsage) Model() interface{} {
	return a
}

// UserID ...
func (a ToolUsage) UserID() string {
	return a.BuildSlug
}
