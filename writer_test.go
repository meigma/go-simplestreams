package simplestreams_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	simplestreams "github.com/meigma/go-simplestreams"
)

func TestBuildIndexRejectsDuplicateContentIDs(t *testing.T) {
	_, err := simplestreams.BuildIndex([]simplestreams.BuildIndexEntry{
		{ContentID: "images", Path: "streams/v1/images.json"},
		{ContentID: "images", Path: "streams/v1/other.json"},
	}, "")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "duplicate")
}

func TestBuildIndexCreatesEntries(t *testing.T) {
	index, err := simplestreams.BuildIndex([]simplestreams.BuildIndexEntry{
		{
			ContentID: "images",
			Path:      "streams/v1/images.json",
			Format:    "products:1.0",
			DataType:  "image-downloads",
			Products:  []string{"ubuntu"},
		},
	}, "Mon, 11 May 2026 00:00:00 +0000")

	require.NoError(t, err)
	entry := index.Entries["images"]
	require.NotNil(t, entry)
	assert.Equal(t, "streams/v1/images.json", entry.Path.String())
	assert.Equal(t, []string{"ubuntu"}, entry.Products)
}

func TestMarshalJSONDocumentRendersBuiltIndex(t *testing.T) {
	index, err := simplestreams.BuildIndex([]simplestreams.BuildIndexEntry{
		{
			ContentID: "images",
			Path:      "streams/v1/images.json",
			Format:    "products:1.0",
			DataType:  "image-downloads",
			Updated:   "Mon, 11 May 2026 00:00:00 +0000",
			Products:  []string{"ubuntu"},
		},
	}, "Mon, 11 May 2026 00:00:00 +0000")
	require.NoError(t, err)

	data, err := simplestreams.MarshalJSONDocument(index)
	require.NoError(t, err)

	assert.JSONEq(t, `{
		"format": "index:1.0",
		"updated": "Mon, 11 May 2026 00:00:00 +0000",
		"index": {
			"images": {
				"format": "products:1.0",
				"datatype": "image-downloads",
				"path": "streams/v1/images.json",
				"updated": "Mon, 11 May 2026 00:00:00 +0000",
				"products": ["ubuntu"]
			}
		}
	}`, string(data))
	assert.Contains(t, string(data), "\n")
}

func TestMarshalJSONDocumentRejectsUnsupportedTypes(t *testing.T) {
	_, err := simplestreams.MarshalJSONDocument(struct{}{})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported JSON document type")
}

func TestCheckDuplicateItemRefs(t *testing.T) {
	ref := simplestreams.ItemRef{ContentID: "images", ProductName: "p", VersionName: "v", ItemName: "i"}

	err := simplestreams.CheckDuplicateItemRefs([]simplestreams.ItemRef{ref, ref})
	require.Error(t, err)
	assert.Contains(t, err.Error(), ref.String())
}
