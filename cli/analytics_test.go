package cli

import (
	"io"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/bitrise-io/bitrise/models"
	"github.com/bitrise-io/go-utils/pathutil"
	"github.com/bitrise-io/go-utils/pointers"
	stepmanModels "github.com/bitrise-io/stepman/models"
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

func Test_readPayload(t *testing.T) {
	tests := []struct {
		name    string
		r       io.Reader
		want    models.BuildRunResultsModel
		wantErr bool
	}{
		{
			name:    "reading an empty string returns an error (no input provided)",
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
		{
			name:    "can read payload",
			r:       strings.NewReader(failedBuildPayload),
			want:    faildBuildBuildRunResult,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := readPayload(tt.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("readPayload() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("readPayload() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_readPayloadFromEmptyFileFails(t *testing.T) {
	tmp, err := pathutil.NormalizedOSTempDirPath("_readPayload_")
	if err != nil {
		t.Errorf("failed to create tmp dir: %s", err)
		return
	}
	pth := filepath.Join(tmp, "file")
	f, err := os.Create(pth)
	if err != nil {
		t.Errorf("failed to open %s: %s", pth, err)
		return
	}
	payload, err := readPayload(f)
	if err == nil {
		t.Errorf("reading empty payload should fail")
	}
	if !reflect.DeepEqual(payload, models.BuildRunResultsModel{}) {
		t.Errorf("readPayload() = %v, want %v", payload, models.BuildRunResultsModel{})
	}
}
