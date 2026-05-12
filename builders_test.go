package simplestreams_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	simplestreams "github.com/meigma/go-simplestreams"
)

func TestProductTreeBuildersSetNamesAndMapKeys(t *testing.T) {
	productFile := simplestreams.NewProductFile("images")
	productFile.DataType = "image-downloads"
	product := productFile.SetProduct("ubuntu:24.04:amd64", &simplestreams.Product{Name: "wrong"})
	version := product.SetVersion("20260511", &simplestreams.Version{Name: "wrong"})
	item := version.SetItem("root", &simplestreams.Item{
		Name:     "wrong",
		FileType: "root.tar.xz",
		Path:     "images/root.tar.xz",
	})

	assert.Equal(t, simplestreams.ProductsFormat, productFile.Format)
	assert.Same(t, product, productFile.Products["ubuntu:24.04:amd64"])
	assert.Same(t, version, product.Versions["20260511"])
	assert.Same(t, item, version.Items["root"])
	assert.Equal(t, "ubuntu:24.04:amd64", product.Name)
	assert.Equal(t, "20260511", version.Name)
	assert.Equal(t, "root", item.Name)

	items := productFile.Items()
	require.Len(t, items, 1)
	assert.Equal(t, "images/ubuntu:24.04:amd64/20260511/root", items[0].Ref.String())

	data, err := simplestreams.MarshalJSONDocument(productFile)
	require.NoError(t, err)
	assert.JSONEq(t, `{
		"format": "products:1.0",
		"content_id": "images",
		"datatype": "image-downloads",
		"products": {
			"ubuntu:24.04:amd64": {
				"versions": {
					"20260511": {
						"items": {
							"root": {
								"ftype": "root.tar.xz",
								"path": "images/root.tar.xz"
							}
						}
					}
				}
			}
		}
	}`, string(data))
}

func TestMetadataHelpersSetAndGetUnknownFields(t *testing.T) {
	productFile := simplestreams.NewProductFile("images")
	product := productFile.SetProduct("ubuntu", nil)
	version := product.SetVersion("20260511", nil)
	item := version.SetItem("root", nil)

	productFile.SetMetadata("license", "example")
	product.SetMetadata("arch", "amd64")
	version.SetMetadata("label", "current")
	item.SetMetadata("custom", true)

	value, ok := productFile.MetadataValue("license")
	require.True(t, ok)
	assert.Equal(t, "example", value)
	value, ok = product.MetadataValue("arch")
	require.True(t, ok)
	assert.Equal(t, "amd64", value)
	value, ok = version.MetadataValue("label")
	require.True(t, ok)
	assert.Equal(t, "current", value)
	value, ok = item.MetadataValue("custom")
	require.True(t, ok)
	assert.Equal(t, true, value)

	var nilItem *simplestreams.Item
	_, ok = nilItem.MetadataValue("custom")
	assert.False(t, ok)

	data, err := simplestreams.MarshalJSONDocument(productFile)
	require.NoError(t, err)
	assert.JSONEq(t, `{
		"format": "products:1.0",
		"content_id": "images",
		"license": "example",
		"products": {
			"ubuntu": {
				"arch": "amd64",
				"versions": {
					"20260511": {
						"label": "current",
						"items": {
							"root": {
								"custom": true
							}
						}
					}
				}
			}
		}
	}`, string(data))
}
