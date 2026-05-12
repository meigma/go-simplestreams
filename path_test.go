package simplestreams_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	simplestreams "github.com/meigma/go-simplestreams"
)

func TestParseRelativePath(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		wantErr bool
	}{
		{name: "metadata path", path: "streams/v1/index.json"},
		{name: "artifact path", path: "images/foo/bar.img"},
		{name: "empty", path: "", wantErr: true},
		{name: "absolute", path: "/streams/v1/index.json", wantErr: true},
		{name: "parent segment", path: "streams/../index.json", wantErr: true},
		{name: "leading parent", path: "../index.json", wantErr: true},
		{name: "ellipsis segment", path: "streams/.../index.json", wantErr: true},
		{name: "empty segment", path: "streams//index.json", wantErr: true},
		{name: "dot segment", path: "streams/./index.json", wantErr: true},
		{name: "backslash", path: `streams\v1\index.json`, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := simplestreams.ParseRelativePath(tt.path)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.path, got.String())
		})
	}
}
