package cli

import (
	"io"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/bitrise-io/go-utils/pathutil"

	"github.com/bitrise-io/bitrise/models"
)

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
			r:       strings.NewReader(`{"project_type":"ios"}`),
			want:    models.BuildRunResultsModel{ProjectType: "ios"},
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
	if err != nil {
		t.Errorf("failed to read payload: %s", err)
		return
	}
	if !reflect.DeepEqual(payload, models.BuildRunResultsModel{}) {
		t.Errorf("readPayload() = %v, want %v", payload, models.BuildRunResultsModel{})
	}
}
