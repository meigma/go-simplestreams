package linuxcontainers

import "github.com/meigma/go-simplestreams/schema"

@go(linuxcontainers)

// ProductFile is the shared Linux Containers image server profile.
#ProductFile: {
	format!:     schema.#ProductsFormat
	content_id!: "images"          @go(ContentID)
	datatype!:   "image-downloads" @go(DataType)
	updated?:    schema.#Timestamp
	products!: [string]: #Product
}

// Product describes shared Incus/LXD image product metadata.
#Product: {
	aliases!:       string
	arch!:          string
	distro?:        string
	os!:            string @go(OS)
	release!:       string
	release_title!: string @go(ReleaseTitle)
	variant!:       string
	requirements!:  #Requirements
	versions!: [string]: #Version
}

// Requirements allows the broader current Linux Containers requirement vocabulary.
#Requirements: {
	cgroup?:           "v1" | "v2"
	secureboot?:       #StringBool @go(SecureBoot)
	cdrom_agent?:      #StringBool @go(CDROMAgent)
	cdrom_cloud_init?: #StringBool @go(CDROMCloudInit)
}

// Version describes one image build.
#Version: {
	items!: [string]: #Item
}

// Item describes shared Linux Containers image artifacts.
#Item: {
	ftype!:                          #FileType            @go(FileType)
	path!:                           schema.#RelativePath @go(,type="github.com/meigma/go-simplestreams/schema".RelativePath)
	size!:                           schema.#Size
	sha256!:                         schema.#Checksum @go(SHA256)
	delta_base?:                     string           @go(DeltaBase)
	combined_sha256?:                schema.#Checksum @go(CombinedSHA256)
	combined_rootxz_sha256?:         schema.#Checksum @go(CombinedRootXZSHA256)
	combined_squashfs_sha256?:       schema.#Checksum @go(CombinedSquashfsSHA256)
	"combined_disk-kvm-img_sha256"?: schema.#Checksum @go(CombinedDiskKVMImageSHA256)
}

// FileType is the artifact vocabulary shared by current Incus and LXD streams.
#FileType: #IncusFileType | #LXDFileType

// IncusFileType is the current Incus/Linux Containers artifact vocabulary.
#IncusFileType: "incus.tar.xz" | "lxd.tar.xz" | "root.tar.xz" | "squashfs" | "disk-kvm.img" | "squashfs.vcdiff"

// LXDFileType is the current Canonical LXD artifact vocabulary.
#LXDFileType: "lxd.tar.xz" | "squashfs" | "disk-kvm.img" | "disk-kvm.img.vcdiff" | "squashfs.vcdiff"

// StringBool is a boolean encoded as a JSON string, matching current image streams.
#StringBool: "true" | "false"
