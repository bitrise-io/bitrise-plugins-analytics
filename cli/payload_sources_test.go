package cli

import (
	"fmt"
	"io"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/bitrise-io/bitrise/models"
	"github.com/bitrise-io/go-utils/pointers"
	stepmanModels "github.com/bitrise-io/stepman/models"
	"github.com/stretchr/testify/require"
)

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

func TestEnvPayloadSource(t *testing.T) {
	tests := []struct {
		name    string
		e       string
		want    models.BuildRunResultsModel
		wantErr bool
	}{
		{
			name:    "reading empty payload returns an error (no input provided)",
			e:       "",
			want:    models.BuildRunResultsModel{},
			wantErr: true,
		},
		{
			name:    "reading invalid payload returns an error",
			e:       "invalid json",
			want:    models.BuildRunResultsModel{},
			wantErr: true,
		},
		{
			name:    "parses valid payload",
			e:       failedBuildPayload,
			want:    faildBuildBuildRunResult,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := EnvPayloadSource{tt.e}.Payload()
			require.Equal(t, err != nil, tt.wantErr, fmt.Sprintf("expected error: %v, got: %v", tt.wantErr, err == nil))
			require.Equal(t, tt.want, got)
		})
	}
}

func TestStdinPayloadSource(t *testing.T) {
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
			got, err := StdinPayloadSource{tt.r}.Payload()
			require.Equal(t, err != nil, tt.wantErr, fmt.Sprintf("expected error: %v, got: %v", tt.wantErr, err == nil))
			require.Equal(t, tt.want, got)
		})
	}
}

func TestStdinPayloadSourceReadsFromPipe(t *testing.T) {
	r, w, err := os.Pipe()
	require.NoError(t, err)

	n, err := w.Write([]byte(failedBuildPayload))
	require.NoError(t, err)
	require.Equal(t, len([]byte(failedBuildPayload)), n)
	require.NoError(t, w.Close())

	payload, err := StdinPayloadSource{r}.Payload()
	require.NoError(t, err)
	require.Equal(t, payload, faildBuildBuildRunResult)
}
