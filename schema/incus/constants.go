package incus

import "github.com/meigma/go-simplestreams/schema/linuxcontainers"

const (
	// ContentIDImages is the product-file content ID used by Incus image streams.
	ContentIDImages = linuxcontainers.ContentIDImages

	// DataTypeImageDownloads is the datatype used by Incus image download streams.
	DataTypeImageDownloads = linuxcontainers.DataTypeImageDownloads
)
