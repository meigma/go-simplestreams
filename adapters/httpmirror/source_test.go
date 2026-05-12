package httpmirror_test

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	simplestreams "github.com/meigma/go-simplestreams"
	"github.com/meigma/go-simplestreams/adapters/httpmirror"
)

func TestSourceOpenReadsFromBaseURL(t *testing.T) {
	var gotUserAgent string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotUserAgent = r.UserAgent()
		assert.Equal(t, "/base/streams/v1/index.json", r.URL.Path)
		_, _ = w.Write([]byte("index"))
	}))
	defer server.Close()

	source, err := httpmirror.New(server.URL+"/base", httpmirror.WithUserAgent("go-simplestreams-test"))
	require.NoError(t, err)
	reader, err := source.Open(context.Background(), "streams/v1/index.json")
	require.NoError(t, err)
	defer reader.Close()

	body, err := io.ReadAll(reader)
	require.NoError(t, err)
	assert.Equal(t, "index", string(body))
	assert.Equal(t, "go-simplestreams-test", gotUserAgent)
}

func TestSourceOpenRejectsNonOKStatus(t *testing.T) {
	server := httptest.NewServer(http.NotFoundHandler())
	defer server.Close()

	source, err := httpmirror.New(server.URL)
	require.NoError(t, err)
	_, err = source.Open(context.Background(), "streams/v1/index.json")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "404")
	assert.ErrorIs(t, err, simplestreams.ErrNotFound)
}

func TestSourceOpenRejectsServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		http.Error(w, "unavailable", http.StatusServiceUnavailable)
	}))
	defer server.Close()

	source, err := httpmirror.New(server.URL)
	require.NoError(t, err)
	_, err = source.Open(context.Background(), "streams/v1/index.json")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "503")
	assert.NotErrorIs(t, err, simplestreams.ErrNotFound)
}

func TestNewRejectsNonHTTPURL(t *testing.T) {
	_, err := httpmirror.New("file:///tmp/mirror")
	require.Error(t, err)
}
