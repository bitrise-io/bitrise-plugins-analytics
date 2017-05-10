package plugins

const (
	// TypeGeneric ...
	TypeGeneric = "_"
	// TypeInit ...
	TypeInit = "init"
	// TypeRun ....
	TypeRun = "run"
)

// PluginRoute ...
type PluginRoute struct {
	Name                   string `yaml:"name"`
	Source                 string `yaml:"source"`
	Version                string `yaml:"version"`
	CommitHash             string `yaml:"commit_hash"`
	Executable             string `yaml:"executable"`
	TriggerEvent           string `yaml:"trigger"`
	LatestAvailableVersion string `yaml:"latest_available_version"`
}

// PluginRouting ...
type PluginRouting struct {
	RouteMap map[string]PluginRoute `yaml:"route_map"`
}

// ExecutableModel ...
type ExecutableModel struct {
	OSX   string `yaml:"osx,omitempty"`
	Linux string `yaml:"linux,omitempty"`
}

// Plugin ...
type Plugin struct {
	Name         string          `yaml:"name,omitempty"`
	Description  string          `yaml:"description,omitempty"`
	Executable   ExecutableModel `yaml:"executable,omitempty"`
	TriggerEvent string          `yaml:"trigger,omitempty"`
	Requirements []Requirement   `yaml:"requirements,omitempty"`
}

// Requirement ...
type Requirement struct {
	Tool       string `yaml:"tool"`
	MinVersion string `yaml:"min_version"`
	MaxVersion string `yaml:"max_version"`
}