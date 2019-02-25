package cli

import (
	"path/filepath"
	"regexp"
	"testing"

	"github.com/bitrise-io/bitrise/models"
	"github.com/bitrise-io/go-utils/command/git"
	"github.com/bitrise-io/go-utils/fileutil"
	"github.com/bitrise-io/go-utils/pathutil"
)

func Test_ensureFormatVersion(t *testing.T) {

	tests := []struct {
		name                        string
		pluginFormatVersionStr      string
		hostBitriseFormatVersionStr string
		wantWarn                    string
		wantErr                     bool
	}{
		{
			name:                        "Simple test",
			pluginFormatVersionStr:      "3",
			hostBitriseFormatVersionStr: "3",
			wantWarn:                    "",
			wantErr:                     false,
		},
		{
			name:                        "Semver test",
			pluginFormatVersionStr:      "1.4.0",
			hostBitriseFormatVersionStr: "1.4.0",
			wantWarn:                    "",
			wantErr:                     false,
		},
		{
			name:                        "Missing host Bitrise CLI format version",
			pluginFormatVersionStr:      "3",
			hostBitriseFormatVersionStr: "",
			wantWarn:                    "This analytics plugin version would need bitrise-cli version >= 1.6.0 to submit analytics",
			wantErr:                     false,
		},
		{
			name:                        "Outdated Host Bitrise format version",
			pluginFormatVersionStr:      "3",
			hostBitriseFormatVersionStr: "1.4.0",
			wantWarn:                    "Outdated bitrise-cli, used format version is lower then the analytics plugin's format version, please update the bitrise-cli",
			wantErr:                     false,
		},
		{
			name:                        "Outdated Plugin Format version",
			pluginFormatVersionStr:      "3",
			hostBitriseFormatVersionStr: "4",
			wantWarn:                    "Outdated analytics plugin, used format version is lower then host bitrise-cli's format version, please update the plugin",
			wantErr:                     false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			warn, err := ensureFormatVersion(tt.pluginFormatVersionStr, tt.hostBitriseFormatVersionStr)
			if (err != nil) != tt.wantErr {
				t.Errorf("ensureFormatVersion() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if warn != tt.wantWarn {
				t.Errorf("ensureFormatVersion() = %v, want %v", warn, tt.wantWarn)
			}
		})
	}
}

func TestFormat(t *testing.T) {
	tmpDir, err := pathutil.NormalizedOSTempDirPath("analytics")
	if err != nil {
		t.Fatalf("Failed to create tmp dir: %s", err)
	}

	gitCmd, err := git.New(tmpDir)
	if err != nil {
		t.Fatalf("Failed to create git cmd: %s", err)
	}

	cloneCmd := gitCmd.Clone("https://github.com/bitrise-io/bitrise.git")
	out, err := cloneCmd.RunAndReturnTrimmedCombinedOutput()
	if err != nil {
		t.Fatalf("Failed clone Bitrise CLI repo: %s", out)
	}

	modelsPth := filepath.Join(tmpDir, "models/models.go")
	content, err := fileutil.ReadStringFromFile(modelsPth)
	if err != nil {
		t.Fatalf("Failed to read %s: %s", modelsPth, err)
	}

	pattern := `Version = "(.*)"`
	re := regexp.MustCompile(pattern)
	match := re.FindStringSubmatch(content)
	if len(match) < 2 {
		t.Fatalf("Failed to find model version in %s using regex: %s", modelsPth, pattern)
	}

	if match[1] != models.Version {
		t.Errorf("Latest Bitrise Model's version: %s, embedded Model's version: %s", match[1], models.Version)
		t.Fatal("Update go dependencies to fetch latest bitrise/models package")
	}
}
