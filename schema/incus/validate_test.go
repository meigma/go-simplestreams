package incus_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	simplestreams "github.com/meigma/go-simplestreams"
	"github.com/meigma/go-simplestreams/schema/incus"
)

func TestValidateRuntimeProductFileAcceptsIncusProfile(t *testing.T) {
	productFile := simplestreams.NewProductFile(incus.ContentIDImages)
	productFile.DataType = incus.DataTypeImageDownloads
	product := productFile.SetProduct("debian:bookworm:amd64:default", nil)
	product.SetMetadata("aliases", "debian/bookworm/default,debian/bookworm")
	product.SetMetadata("arch", "amd64")
	product.SetMetadata("os", "debian")
	product.SetMetadata("release", "bookworm")
	product.SetMetadata("release_title", "Debian 12 bookworm")
	product.SetMetadata("variant", "default")
	product.SetMetadata("requirements", map[string]any{"cdrom_agent": "true"})

	version := product.SetVersion("20260511_05:24", nil)
	size := int64(688)
	item := version.SetItem("incus.tar.xz", &simplestreams.Item{
		FileType: "incus.tar.xz",
		Path:     "images/debian/bookworm/amd64/default/20260511_05:24/incus.tar.xz",
		Size:     &size,
		SHA256:   "39e0ee4e85dd5fcc153ab6e3b6575113f679b5f76a4c1886088eb288c9cc71b5",
	})
	item.SetMetadata("combined_disk-kvm-img_sha256", "96d197fea165faf7478ff2e59e710a1f75f45277006be3103ea235bb7456fc74")

	err := incus.ValidateRuntimeProductFile(productFile)
	require.NoError(t, err)
}

func TestValidateRuntimeProductFileRejectsInvalidIncusProfile(t *testing.T) {
	productFile := simplestreams.NewProductFile("other")
	productFile.DataType = incus.DataTypeImageDownloads

	err := incus.ValidateRuntimeProductFile(productFile)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "validate runtime product file")
}
