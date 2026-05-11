package incus

import (
	"github.com/meigma/go-simplestreams/schema"
	"github.com/meigma/go-simplestreams/schema/linuxcontainers"
)

@go(incus)

// ProductFile is the Incus/Linux Containers image server profile.
#ProductFile: {
	format!:     schema.#ProductsFormat
	content_id!: "images"          @go(ContentID)
	datatype!:   "image-downloads" @go(DataType)
	updated?:    schema.#Timestamp
	products!: [string]: #Product
}

// Product describes Incus/Linux Containers image product metadata.
#Product: {
	aliases!:       string
	arch!:          string
	os!:            string @go(OS)
	release!:       string
	release_title!: string @go(ReleaseTitle)
	variant!:       string
	requirements!:  #Requirements
	versions!: [string]: #Version
}

// Requirements allows the current Incus/Linux Containers requirement vocabulary.
#Requirements: {
	cgroup?:           "v1" | "v2"
	secureboot?:       linuxcontainers.#StringBool @go(SecureBoot)
	cdrom_agent?:      linuxcontainers.#StringBool @go(CDROMAgent)
	cdrom_cloud_init?: linuxcontainers.#StringBool @go(CDROMCloudInit)
}

// Version describes one Incus/Linux Containers image build.
#Version: {
	items!: [string]: #Item
}

// Item describes Incus/Linux Containers image artifacts.
#Item: {
	ftype!:                          linuxcontainers.#IncusFileType @go(FileType)
	path!:                           schema.#RelativePath           @go(,type="github.com/meigma/go-simplestreams/schema".RelativePath)
	size!:                           schema.#Size
	sha256!:                         schema.#Checksum @go(SHA256)
	delta_base?:                     string           @go(DeltaBase)
	combined_sha256?:                schema.#Checksum @go(CombinedSHA256)
	combined_rootxz_sha256?:         schema.#Checksum @go(CombinedRootXZSHA256)
	combined_squashfs_sha256?:       schema.#Checksum @go(CombinedSquashfsSHA256)
	"combined_disk-kvm-img_sha256"?: schema.#Checksum @go(CombinedDiskKVMImageSHA256)
}
