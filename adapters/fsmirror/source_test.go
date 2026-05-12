package fsmirror_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	simplestreams "github.com/meigma/go-simplestreams"
	"github.com/meigma/go-simplestreams/adapters/fsmirror"
)

func TestSourceOpenReadsFromRoot(t *testing.T) {
	root := t.TempDir()
	err := os.MkdirAll(filepath.Join(root, "streams", "v1"), 0o755)
	require.NoError(t, err)
	err = os.WriteFile(filepath.Join(root, "streams", "v1", "index.json"), []byte("index"), 0o644)
	require.NoError(t, err)

	source, err := fsmirror.New(root)
	require.NoError(t, err)
	reader, err := source.Open(context.Background(), "streams/v1/index.json")
	require.NoError(t, err)
	defer reader.Close()

	body := make([]byte, 5)
	_, err = reader.Read(body)
	require.NoError(t, err)
	assert.Equal(t, "index", string(body))
}

func TestSourceOpenRejectsUnsafePath(t *testing.T) {
	source, err := fsmirror.New(t.TempDir())
	require.NoError(t, err)

	_, err = source.Open(context.Background(), "../secret")
	require.Error(t, err)
}

func TestSourceOpenReportsMissingPath(t *testing.T) {
	source, err := fsmirror.New(t.TempDir())
	require.NoError(t, err)

	_, err = source.Open(context.Background(), "streams/v1/index.sjson")
	require.Error(t, err)
	assert.ErrorIs(t, err, simplestreams.ErrNotFound)
}
