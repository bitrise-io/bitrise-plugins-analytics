package cli

import "testing"
import "github.com/stretchr/testify/require"

func TestEnsureFormatVersion(t *testing.T) {
	t.Log("format version equals")
	{
		warn, err := ensureFormatVersion("3", "3")
		require.NoError(t, err)
		require.Equal(t, "", warn)

		warn, err = ensureFormatVersion("1.4.0", "1.4.0")
		require.NoError(t, err)
		require.Equal(t, "", warn)
	}

	t.Log("bitrise version < 1.6.0 - does not export format version")
	{
		warn, err := ensureFormatVersion("3", "")
		require.NoError(t, err)
		require.Equal(t, "This analytics plugin version would need bitrise-cli version >= 1.6.0 to submit analytics", warn)
	}

	t.Log("bitrise format version < plugin format version")
	{
		warn, err := ensureFormatVersion("3", "1.4.0")
		require.NoError(t, err)
		require.Equal(t, "Outdated bitrise-cli, used format version is lower then the analytics plugin's format version, please update the bitrise-cli", warn)
	}

	t.Log("bitrise format version > plugin format version")
	{
		warn, err := ensureFormatVersion("3", "4")
		require.NoError(t, err)
		require.Equal(t, "Outdated analytics plugin, used format version is lower then host bitrise-cli's format version, please update the plugin", warn)
	}
}
