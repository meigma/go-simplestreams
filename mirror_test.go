package simplestreams_test

import (
	"context"
	"fmt"
	"io"
	"sort"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	simplestreams "github.com/meigma/go-simplestreams"
)

type mapSource struct {
	files map[string]string
	errs  map[string]error
	opens map[string]int
}

func (source *mapSource) Open(_ context.Context, path simplestreams.RelativePath) (io.ReadCloser, error) {
	source.opens[path.String()]++
	if err := source.errs[path.String()]; err != nil {
		return nil, err
	}
	content, ok := source.files[path.String()]
	if !ok {
		return nil, fmt.Errorf("%w: %s", simplestreams.ErrNotFound, path)
	}
	return io.NopCloser(strings.NewReader(content)), nil
}

func TestMirrorLoadsIndexAndProductFilesLazily(t *testing.T) {
	source := &mapSource{
		files: map[string]string{
			"streams/v1/index.json": `{
				"format": "index:1.0",
				"index": {
					"com.example:download": {
						"format": "products:1.0",
						"datatype": "image-downloads",
						"path": "streams/v1/download.json",
						"products": ["com.example:product"]
					}
				}
			}`,
			"streams/v1/download.json": `{
				"format": "products:1.0",
				"content_id": "com.example:download",
				"datatype": "image-downloads",
				"products": {
					"com.example:product": {
						"versions": {
							"20260511": {
								"items": {
									"disk": {"ftype": "disk.img", "path": "images/disk.img"}
								}
							}
						}
					}
				}
			}`,
		},
		opens: map[string]int{},
	}

	mirror, err := simplestreams.NewMirror(source)
	require.NoError(t, err)

	index, err := mirror.Index(context.Background())
	require.NoError(t, err)
	require.Len(t, index.Entries, 1)
	assert.Equal(t, 1, source.opens["streams/v1/index.json"])
	assert.Zero(t, source.opens["streams/v1/download.json"])

	entry := index.Entries["com.example:download"]
	productFile, err := entry.ProductFile(context.Background())
	require.NoError(t, err)
	assert.Equal(t, "com.example:download", productFile.ContentID)
	assert.Equal(t, 1, source.opens["streams/v1/download.json"])

	again, err := entry.ProductFile(context.Background())
	require.NoError(t, err)
	assert.Same(t, productFile, again)
	assert.Equal(t, 1, source.opens["streams/v1/download.json"])
}

func TestIndexEntryLoadsNestedIndexLazily(t *testing.T) {
	source := &mapSource{
		files: map[string]string{
			"streams/v1/index.json": `{
				"format": "index:1.0",
				"index": {
					"nested": {
						"format": "index:1.0",
						"path": "streams/v1/nested.json"
					}
				}
			}`,
			"streams/v1/nested.json": `{
				"format": "index:1.0",
				"index": {
					"images": {
						"format": "products:1.0",
						"path": "streams/v1/images.json"
					}
				}
			}`,
			"streams/v1/images.json": `{
				"format": "products:1.0",
				"content_id": "images",
				"products": {}
			}`,
		},
		opens: map[string]int{},
	}
	mirror, err := simplestreams.NewMirror(source)
	require.NoError(t, err)

	index, err := mirror.Index(context.Background())
	require.NoError(t, err)
	assert.Zero(t, source.opens["streams/v1/nested.json"])

	nestedEntry := index.Entries["nested"]
	_, err = nestedEntry.ProductFile(context.Background())
	require.Error(t, err)
	assert.Contains(t, err.Error(), "nested index")
	assert.Zero(t, source.opens["streams/v1/nested.json"])

	nested, err := nestedEntry.Index(context.Background())
	require.NoError(t, err)
	assert.Contains(t, nested.Entries, "images")
	assert.Equal(t, 1, source.opens["streams/v1/nested.json"])

	again, err := nestedEntry.Index(context.Background())
	require.NoError(t, err)
	assert.Same(t, nested, again)
	assert.Equal(t, 1, source.opens["streams/v1/nested.json"])

	productFile, err := nested.Entries["images"].ProductFile(context.Background())
	require.NoError(t, err)
	assert.Equal(t, "images", productFile.ContentID)
}

func TestIndexPreservesUnknownMetadata(t *testing.T) {
	source := &mapSource{
		files: map[string]string{
			"streams/v1/index.json": `{
				"format": "index:1.0",
				"updated": "Mon, 11 May 2026 00:00:00 +0000",
				"source": "fixture",
				"index": {
					"images": {
						"format": "products:1.0",
						"path": "streams/v1/images.json",
						"region": "global"
					}
				}
			}`,
		},
		opens: map[string]int{},
	}
	mirror, err := simplestreams.NewMirror(source)
	require.NoError(t, err)

	index, err := mirror.Index(context.Background())
	require.NoError(t, err)
	assert.Equal(t, "fixture", index.Metadata["source"])
	assert.Equal(t, "global", index.Entries["images"].Metadata["region"])

	data, err := simplestreams.MarshalJSONDocument(index)
	require.NoError(t, err)
	assert.JSONEq(t, `{
		"format": "index:1.0",
		"updated": "Mon, 11 May 2026 00:00:00 +0000",
		"source": "fixture",
		"index": {
			"images": {
				"format": "products:1.0",
				"path": "streams/v1/images.json",
				"region": "global"
			}
		}
	}`, string(data))
}

func TestProductFileRejectsContentIDMismatchByDefault(t *testing.T) {
	source := &mapSource{
		files: map[string]string{
			"streams/v1/index.json": `{
				"format": "index:1.0",
				"index": {"expected": {"path": "streams/v1/products.json"}}
			}`,
			"streams/v1/products.json": `{
				"format": "products:1.0",
				"content_id": "actual",
				"products": {}
			}`,
		},
		opens: map[string]int{},
	}

	mirror, err := simplestreams.NewMirror(source)
	require.NoError(t, err)
	index, err := mirror.Index(context.Background())
	require.NoError(t, err)

	_, err = index.Entries["expected"].ProductFile(context.Background())
	require.Error(t, err)
	assert.Contains(t, err.Error(), "does not match")
}

func TestProductFileCanAllowContentIDMismatch(t *testing.T) {
	source := &mapSource{
		files: map[string]string{
			"streams/v1/index.json": `{
				"format": "index:1.0",
				"index": {"expected": {"path": "streams/v1/products.json"}}
			}`,
			"streams/v1/products.json": `{
				"format": "products:1.0",
				"content_id": "actual",
				"products": {}
			}`,
		},
		opens: map[string]int{},
	}

	mirror, err := simplestreams.NewMirror(source, simplestreams.WithAllowContentIDMismatch(true))
	require.NoError(t, err)
	index, err := mirror.Index(context.Background())
	require.NoError(t, err)

	productFile, err := index.Entries["expected"].ProductFile(context.Background())
	require.NoError(t, err)
	assert.Equal(t, "actual", productFile.ContentID)
}

func TestProductFileItemsPreserveUnknownMetadataAndApplyInheritance(t *testing.T) {
	productFile := loadProductFile(t, `{
		"format": "products:1.0",
		"content_id": "images",
		"datatype": "image-downloads",
		"license": "example",
		"products": {
			"ubuntu:24.04:amd64": {
				"arch": "amd64",
				"release": "24.04",
				"versions": {
					"20260510": {
						"label": "old",
						"items": {
							"disk": {"ftype": "disk.img", "path": "images/old.img", "sha256": "old"}
						}
					},
					"20260511": {
						"label": "current",
						"items": {
							"root": {
								"ftype": "root.tar.xz",
								"path": "images/root.tar.xz",
								"sha256": "abc",
								"size": 3,
								"custom": true
							}
						}
					}
				}
			}
		}
	}`)

	items := productFile.Items()
	require.Len(t, items, 2)
	latest, ok := simplestreams.LatestByVersion(
		simplestreams.FilterItems(items, simplestreams.MatchFileType("root.tar.xz")),
	)
	require.True(t, ok)

	assert.Equal(t, "images", latest.Ref.ContentID)
	assert.Equal(t, "ubuntu:24.04:amd64", latest.Ref.ProductName)
	assert.Equal(t, "20260511", latest.Ref.VersionName)
	assert.Equal(t, "root", latest.Ref.ItemName)
	assert.Equal(t, "example", latest.Metadata["license"])
	assert.Equal(t, "amd64", latest.Metadata["arch"])
	assert.Equal(t, "24.04", latest.Metadata["release"])
	assert.Equal(t, "current", latest.Metadata["label"])
	assert.Equal(t, "root.tar.xz", latest.Metadata["ftype"])
	assert.Equal(t, "images/root.tar.xz", latest.Metadata["path"])
	assert.Equal(t, int64(3), latest.Metadata["size"])
	assert.Equal(t, true, latest.Metadata["custom"])
}

func TestMarshalJSONDocumentRendersProductFile(t *testing.T) {
	productFile := loadProductFile(t, `{
		"format": "products:1.0",
		"content_id": "images",
		"datatype": "image-downloads",
		"license": "example",
		"products": {
			"ubuntu": {
				"arch": "amd64",
				"versions": {
					"20260511": {
						"label": "current",
						"items": {
							"root": {
								"ftype": "root.tar.xz",
								"path": "images/root.tar.xz",
								"sha256": "abc",
								"size": 3,
								"mirrors": ["https://mirror.example"],
								"custom": true
							}
						}
					}
				}
			}
		}
	}`)
	product := productFile.Products["ubuntu"]
	version := product.Versions["20260511"]
	item := version.Items["root"]
	productFile.Metadata["format"] = "bad"
	productFile.Metadata["content_id"] = "bad"
	product.Metadata["versions"] = "bad"
	version.Metadata["items"] = "bad"
	item.Metadata["ftype"] = "bad"
	item.Metadata["path"] = "bad"
	item.Metadata["size"] = int64(99)

	data, err := simplestreams.MarshalJSONDocument(productFile)
	require.NoError(t, err)

	assert.JSONEq(t, `{
		"format": "products:1.0",
		"content_id": "images",
		"datatype": "image-downloads",
		"license": "example",
		"products": {
			"ubuntu": {
				"arch": "amd64",
				"versions": {
					"20260511": {
						"label": "current",
						"items": {
							"root": {
								"ftype": "root.tar.xz",
								"path": "images/root.tar.xz",
								"sha256": "abc",
								"size": 3,
								"mirrors": ["https://mirror.example"],
								"custom": true
							}
						}
					}
				}
			}
		}
	}`, string(data))
	assert.NotContains(t, string(data), "ContentID")
	assert.NotContains(t, string(data), "bad")
}

func TestProductFileWalkItemsUsesDeterministicOrder(t *testing.T) {
	productFile := loadProductFile(t, `{
		"format": "products:1.0",
		"content_id": "images",
		"products": {
			"b": {"versions": {"2": {"items": {"b": {}, "a": {}}}}},
			"a": {"versions": {"2": {"items": {"b": {}}}, "1": {"items": {"a": {}}}}}
		}
	}`)

	got := []string{}
	err := productFile.WalkItems(func(item simplestreams.ItemView) error {
		got = append(got, item.Ref.String())
		return nil
	})
	require.NoError(t, err)

	want := append([]string(nil), got...)
	sort.Strings(want)
	assert.Equal(t, want, got)
}

func loadProductFile(t *testing.T, content string) *simplestreams.ProductFile {
	t.Helper()
	source := &mapSource{
		files: map[string]string{
			"streams/v1/index.json": `{
				"format": "index:1.0",
				"index": {"images": {"path": "streams/v1/images.json"}}
			}`,
			"streams/v1/images.json": content,
		},
		opens: map[string]int{},
	}
	mirror, err := simplestreams.NewMirror(source)
	require.NoError(t, err)
	index, err := mirror.Index(context.Background())
	require.NoError(t, err)
	productFile, err := index.Entries["images"].ProductFile(context.Background())
	require.NoError(t, err)
	return productFile
}
