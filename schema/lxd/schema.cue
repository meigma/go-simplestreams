package lxd

import (
	"github.com/meigma/go-simplestreams/schema"
	"github.com/meigma/go-simplestreams/schema/linuxcontainers"
)

@go(lxd)

// ProductFile is the Canonical LXD image server profile.
#ProductFile: {
	format!:     schema.#ProductsFormat
	content_id!: "images"          @go(ContentID)
	datatype!:   "image-downloads" @go(DataType)
	updated?:    schema.#Timestamp
	products!: [string]: #Product
}

// Product describes Canonical LXD image product metadata.
#Product: {
	aliases!:       string
	arch!:          string
	distro!:        string
	os!:            string @go(OS)
	release!:       string
	release_title!: string @go(ReleaseTitle)
	variant!:       string
	requirements!:  #Requirements
	versions!: [string]: #Version
}

// Requirements allows the narrower current Canonical LXD requirement vocabulary.
#Requirements: {
	cgroup?:     "v1" | "v2"
	secureboot?: linuxcontainers.#StringBool @go(SecureBoot)
}

// Version describes one Canonical LXD image build.
#Version: {
	items!: [string]: #Item
}

// Item describes Canonical LXD image artifacts.
#Item: {
	ftype!:                          linuxcontainers.#LXDFileType @go(FileType)
	path!:                           schema.#RelativePath         @go(,type="github.com/meigma/go-simplestreams/schema".RelativePath)
	size!:                           schema.#Size
	sha256!:                         schema.#Checksum @go(SHA256)
	delta_base?:                     string           @go(DeltaBase)
	combined_squashfs_sha256?:       schema.#Checksum @go(CombinedSquashfsSHA256)
	"combined_disk-kvm-img_sha256"?: schema.#Checksum @go(CombinedDiskKVMImageSHA256)
}
