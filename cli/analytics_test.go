package cli

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/bitrise-io/bitrise/models"
	"github.com/bitrise-io/go-utils/pointers"
	stepmanModels "github.com/bitrise-io/stepman/models"
)

func Test_readPayload(t *testing.T) {
	tests := []struct {
		name    string
		r       io.Reader
		want    models.BuildRunResultsModel
		wantErr bool
	}{
		{
			name:    "reading empty payload returns an error (no input provided)",
			r:       strings.NewReader(""),
			want:    models.BuildRunResultsModel{},
			wantErr: true,
		},
		{
			name:    "reading invalid payload returns an error",
			r:       strings.NewReader("invalid json"),
			want:    models.BuildRunResultsModel{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := readPayload(tt.r)
			require.Equal(t, err != nil, tt.wantErr, fmt.Sprintf("expected error: %v, got: %v", tt.wantErr, err == nil))
			require.Equal(t, tt.want, got)
		})
	}
}

const failedBuildPayload = `{
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

var faildBuildBuildRunResult = models.BuildRunResultsModel{
	StepmanUpdates: map[string]int{"https://github.com/bitrise-io/bitrise-steplib.git": 1},
	FailedSteps: []models.StepRunResultsModel{
		models.StepRunResultsModel{
			StepInfo: stepmanModels.StepInfoModel{
				Library:       "https://github.com/bitrise-io/bitrise-steplib.git",
				ID:            "script",
				Version:       "1.1.3",
				LatestVersion: "1.1.3",
				Step: stepmanModels.StepModel{
					Title:         pointers.NewStringPtr("script"),
					SourceCodeURL: pointers.NewStringPtr("https://github.com/bitrise-io/steps-script"),
					SupportURL:    pointers.NewStringPtr("https://github.com/bitrise-io/steps-script/issues"),
				},
			},
			Status:   1,
			Idx:      0,
			RunTime:  time.Duration(2027588963),
			ErrorStr: "exit status 1",
			ExitCode: 1,
		},
	},
}

func Test_readPayloadFromStdin(t *testing.T) {
	if os.Getenv("READ_STDIN") == "1" {
		payload, err := readPayload(os.Stdin)
		require.NoError(t, err)
		require.Equal(t, payload, faildBuildBuildRunResult)
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=Test_readPayloadFromStdin")
	cmd.Env = append(os.Environ(), "READ_STDIN=1")
	cmd.Stdin = strings.NewReader(failedBuildPayload)

	err := cmd.Run()
	require.NoError(t, err)
}
