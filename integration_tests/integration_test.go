package integration

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/bitrise-io/go-utils/command"
	"github.com/bitrise-io/go-utils/pathutil"
)

var binPth string

func TestMain(m *testing.M) {
	fmt.Println("main run")
	tmpDir, err := pathutil.NormalizedOSTempDirPath("analytics")
	if err != nil {
		fmt.Printf("Failed to create tmp dir: %s", err)
		os.Exit(1)
	}

	dir, err := os.Getwd()
	if err != nil {
		fmt.Printf("Failed to get current dir: %s", err)
		os.Exit(1)
	}

	binPth = filepath.Join(tmpDir, "bitrise-plugins-analytics")
	build := command.New("go", "build", "-o", binPth)
	build.SetDir(filepath.Dir(dir))
	out, err := build.RunAndReturnTrimmedCombinedOutput()
	if err != nil {
		fmt.Printf("Failed to build analytics plugin: %s", out)
		os.Exit(1)
	}

	os.Exit(m.Run())
}
