package simplestreams_test

import (
	"crypto/sha256"
	"crypto/sha512"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	simplestreams "github.com/meigma/go-simplestreams"
)

func TestArtifactRefPaths(t *testing.T) {
	size := int64(3)
	artifact := simplestreams.ArtifactRef{
		Path:    "images/root.tar.xz",
		Mirrors: []simplestreams.MirrorPrefix{"https://mirror.example/base/", "https://other.example"},
		Size:    &size,
	}

	assert.Equal(t, []string{
		"images/root.tar.xz",
		"https://mirror.example/base/images/root.tar.xz",
		"https://other.example/images/root.tar.xz",
	}, artifact.Paths())
}

func TestBestChecksumPrefersStrongestAvailable(t *testing.T) {
	checksum, ok := simplestreams.BestChecksum(map[string]string{
		"md5":    "weak",
		"sha256": "better",
		"sha512": "best",
	})

	require.True(t, ok)
	assert.Equal(t, "sha512", checksum.Algorithm)
	assert.Equal(t, "best", checksum.Value)
}

func TestVerifyReader(t *testing.T) {
	body := "abc"
	sha256sum := fmt.Sprintf("%x", sha256.Sum256([]byte(body)))
	sha512sum := fmt.Sprintf("%x", sha512.Sum512([]byte(body)))
	size := int64(len(body))
	mismatchSize := int64(99)

	tests := []struct {
		name      string
		body      string
		checksums map[string]string
		size      *int64
		wantErr   string
	}{
		{name: "matches size and checksum", body: body, checksums: map[string]string{"sha256": sha256sum}, size: &size},
		{
			name:      "prefers sha512",
			body:      body,
			checksums: map[string]string{"sha256": "wrong", "sha512": sha512sum},
			size:      &size,
		},
		{
			name:      "rejects size mismatch",
			body:      body,
			checksums: map[string]string{"sha256": sha256sum},
			size:      &mismatchSize,
			wantErr:   "size mismatch",
		},
		{
			name:      "rejects checksum mismatch",
			body:      body,
			checksums: map[string]string{"sha256": "wrong"},
			size:      &size,
			wantErr:   "sha256 mismatch",
		},
		{name: "size only", body: body, size: &size},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := simplestreams.VerifyReader(strings.NewReader(tt.body), tt.checksums, tt.size)
			if tt.wantErr != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.wantErr)
				return
			}
			require.NoError(t, err)
		})
	}
}
