package integration

import (
	"fmt"
	"testing"

	"github.com/bitrise-io/bitrise-plugins-analytics/version"
	"github.com/bitrise-io/go-utils/command"
	"github.com/bitrise-io/go-utils/pathutil"
	"github.com/stretchr/testify/require"
)

var helpStr = fmt.Sprintf(`NAME:
   bitrise-plugins-analytics - Bitrise Analytics plugin

USAGE:
   bitrise-plugins-analytics [global options] command [command options] [arguments...]

VERSION:
   %s

COMMANDS:
     on       Turn sending anonimized usage information on.
     off      Turn sending anonimized usage information off.
     help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --loglevel value, -l value  Log level (options: debug, info, warn, error, fatal, panic). [$LOGLEVEL]
   --help, -h                  show help
   --version, -v               print the version`, version.VERSION)

func Test_HelpTest(t *testing.T) {
	t.Log("help command")
	{
		tmpDir, err := pathutil.NormalizedOSTempDirPath("")
		require.NoError(t, err)

		cmd := command.New(binPth, "help")
		cmd.SetDir(tmpDir)
		out, err := cmd.RunAndReturnTrimmedCombinedOutput()
		require.NoError(t, err, out)
		require.Equal(t, helpStr, out)
	}

	t.Log("help short command")
	{
		tmpDir, err := pathutil.NormalizedOSTempDirPath("")
		require.NoError(t, err)

		cmd := command.New(binPth, "h")
		cmd.SetDir(tmpDir)
		out, err := cmd.RunAndReturnTrimmedCombinedOutput()
		require.NoError(t, err, out)
		require.Equal(t, helpStr, out)
	}

	t.Log("help flag")
	{
		tmpDir, err := pathutil.NormalizedOSTempDirPath("")
		require.NoError(t, err)

		cmd := command.New(binPth, "--help")
		cmd.SetDir(tmpDir)
		out, err := cmd.RunAndReturnTrimmedCombinedOutput()
		require.NoError(t, err, out)
		require.Equal(t, helpStr, out)
	}

	t.Log("help short flag")
	{
		tmpDir, err := pathutil.NormalizedOSTempDirPath("")
		require.NoError(t, err)

		cmd := command.New(binPth, "-h")
		cmd.SetDir(tmpDir)
		out, err := cmd.RunAndReturnTrimmedCombinedOutput()
		require.NoError(t, err, out)
		require.Equal(t, helpStr, out)
	}
}
