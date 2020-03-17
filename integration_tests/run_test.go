package integration

import (
	"strings"
	"testing"

	"github.com/bitrise-io/bitrise-plugins-analytics/configs"
	bitriseConfigs "github.com/bitrise-io/bitrise/configs"
	"github.com/bitrise-io/bitrise/models"
	"github.com/bitrise-io/bitrise/plugins"
	"github.com/bitrise-io/go-utils/command"
	"github.com/bitrise-io/go-utils/pathutil"
	"github.com/stretchr/testify/require"
)

const failedBuildPayload = `{
   "start_time":"2017-05-10T08:28:00.661803042+02:00",
   "stepman_updates":{
      "https://github.com/bitrise-io/bitrise-steplib.git":1
   },
   "success_steps":null,
   "failed_steps":[
      {
         "step_info":{
            "library":"https://github.com/bitrise-io/bitrise-steplib.git",
            "id":"script",
            "version":"1.1.3",
            "latest_version":"1.1.3",
            "info":{

            },
            "step":{
               "title":"script",
               "source_code_url":"https://github.com/bitrise-io/steps-script",
               "support_url":"https://github.com/bitrise-io/steps-script/issues"
            }
         },
         "status":1,
         "idx":0,
         "run_time":2027588963,
         "error_str":"exit status 1",
         "exit_code":1
      }
   ],
   "failed_skippable_steps":null,
   "skipped_steps":null
}`

const successBuildPayload = `{
   "start_time":"2017-05-10T08:29:25.917342266+02:00",
   "stepman_updates":{
      "https://github.com/bitrise-io/bitrise-steplib.git":1
   },
   "success_steps":[
      {
         "step_info":{
            "library":"https://github.com/bitrise-io/bitrise-steplib.git",
            "id":"script",
            "version":"1.1.3",
            "latest_version":"1.1.3",
            "info":{

            },
            "step":{
               "title":"script",
               "source_code_url":"https://github.com/bitrise-io/steps-script",
               "support_url":"https://github.com/bitrise-io/steps-script/issues"
            }
         },
         "status":0,
         "idx":0,
         "run_time":12144946753,
         "error_str":"",
         "exit_code":0
      }
   ],
   "failed_steps":null,
   "failed_skippable_steps":null,
   "skipped_steps":null
}`

func TestStdinPayload(t *testing.T) {
	t.Log("success build")
	{
		tmpDir, err := pathutil.NormalizedOSTempDirPath("")
		require.NoError(t, err)

		envs := []string{
			plugins.PluginConfigDataDirKey + "=" + tmpDir,
			bitriseConfigs.CIModeEnvKey + "=false",

			plugins.PluginConfigPluginModeKey + "=" + string(plugins.TriggerMode),
			plugins.PluginConfigFormatVersionKey + "=" + models.Version,
		}

		cmd := command.New(binPth)
		cmd.SetEnvs(envs...)
		cmd.SetStdin(strings.NewReader(successBuildPayload))
		out, err := cmd.RunAndReturnTrimmedCombinedOutput()
		require.NoError(t, err, out)
	}

	t.Log("failed build")
	{
		tmpDir, err := pathutil.NormalizedOSTempDirPath("")
		require.NoError(t, err)

		envs := []string{
			plugins.PluginConfigDataDirKey + "=" + tmpDir,
			bitriseConfigs.CIModeEnvKey + "=false",

			plugins.PluginConfigPluginModeKey + "=" + string(plugins.TriggerMode),
			plugins.PluginConfigFormatVersionKey + "=" + models.Version,
		}

		cmd := command.New(binPth)
		cmd.SetEnvs(envs...)
		cmd.SetStdin(strings.NewReader(failedBuildPayload))
		out, err := cmd.RunAndReturnTrimmedCombinedOutput()
		require.NoError(t, err, out)
	}
}

func TestEnvPayload(t *testing.T) {
	t.Log("success build")
	{
		tmpDir, err := pathutil.NormalizedOSTempDirPath("")
		require.NoError(t, err)

		envs := []string{
			plugins.PluginConfigDataDirKey + "=" + tmpDir,
			bitriseConfigs.CIModeEnvKey + "=false",

			plugins.PluginConfigPluginModeKey + "=" + string(plugins.TriggerMode),
			plugins.PluginConfigFormatVersionKey + "=" + models.Version,
			configs.PluginConfigPayloadKey + "=" + successBuildPayload,
		}

		cmd := command.New(binPth)
		cmd.SetEnvs(envs...)
		out, err := cmd.RunAndReturnTrimmedCombinedOutput()
		require.NoError(t, err, out)
	}

	t.Log("failed build")
	{
		tmpDir, err := pathutil.NormalizedOSTempDirPath("")
		require.NoError(t, err)

		envs := []string{
			plugins.PluginConfigDataDirKey + "=" + tmpDir,
			bitriseConfigs.CIModeEnvKey + "=false",

			plugins.PluginConfigPluginModeKey + "=" + string(plugins.TriggerMode),
			plugins.PluginConfigFormatVersionKey + "=" + models.Version,
			configs.PluginConfigPayloadKey + "=" + failedBuildPayload,
		}

		cmd := command.New(binPth)
		cmd.SetEnvs(envs...)
		out, err := cmd.RunAndReturnTrimmedCombinedOutput()
		require.NoError(t, err, out)
	}
}
