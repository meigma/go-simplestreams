package schema

// LinuxContainersProductFile is the shared Linux Containers image server profile.
#LinuxContainersProductFile: {
	format!:     #ProductsFormat
	content_id!: "images"          @go(ContentID)
	datatype!:   "image-downloads" @go(DataType)
	updated?:    #Timestamp
	products!: [string]: #LinuxContainersProduct
}

// LinuxContainersProduct describes shared Incus/LXD image product metadata.
#LinuxContainersProduct: {
	aliases!:       string
	arch!:          string
	os!:            string @go(OS)
	release!:       string
	release_title!: string @go(ReleaseTitle)
	variant!:       string
	requirements!:  #LinuxContainersRequirements
	versions!: [string]: #LinuxContainersVersion
}

// LinuxContainersRequirements describes image runtime requirements.
#LinuxContainersRequirements: {
	cgroup?:           "v1" | "v2"
	secureboot?:       #StringBool @go(SecureBoot)
	cdrom_agent?:      #StringBool @go(CDROMAgent)
	cdrom_cloud_init?: #StringBool @go(CDROMCloudInit)
}

// LinuxContainersVersion describes one image build.
#LinuxContainersVersion: {
	items!: [string]: #LinuxContainersItem
}

// LinuxContainersItem describes shared Linux Containers image artifacts.
#LinuxContainersItem: {
	ftype!:                          #LinuxContainersFileType @go(FileType)
	path!:                           #RelativePath            @go(,type=RelativePath)
	size!:                           #Size
	sha256!:                         #Checksum @go(SHA256)
	delta_base?:                     string    @go(DeltaBase)
	combined_sha256?:                #Checksum @go(CombinedSHA256)
	combined_rootxz_sha256?:         #Checksum @go(CombinedRootXZSHA256)
	combined_squashfs_sha256?:       #Checksum @go(CombinedSquashfsSHA256)
	"combined_disk-kvm-img_sha256"?: #Checksum @go(CombinedDiskKVMImageSHA256)
}

// LinuxContainersFileType is the artifact vocabulary shared by current Incus and LXD streams.
#LinuxContainersFileType: #IncusFileType | #LXDFileType

// IncusProductFile is the Incus/Linux Containers image server profile.
#IncusProductFile: {
	format!:     #ProductsFormat
	content_id!: "images"          @go(ContentID)
	datatype!:   "image-downloads" @go(DataType)
	updated?:    #Timestamp
	products!: [string]: #IncusProduct
}

// IncusProduct describes Incus/Linux Containers image product metadata.
#IncusProduct: {
	aliases!:       string
	arch!:          string
	os!:            string @go(OS)
	release!:       string
	release_title!: string @go(ReleaseTitle)
	variant!:       string
	requirements!:  #IncusRequirements
	versions!: [string]: #IncusVersion
}

// IncusRequirements allows the broader current Linux Containers requirement vocabulary.
#IncusRequirements: {
	cgroup?:           "v1" | "v2"
	secureboot?:       #StringBool @go(SecureBoot)
	cdrom_agent?:      #StringBool @go(CDROMAgent)
	cdrom_cloud_init?: #StringBool @go(CDROMCloudInit)
}

// IncusVersion describes one Incus/Linux Containers image build.
#IncusVersion: {
	items!: [string]: #IncusItem
}

// IncusItem describes Incus/Linux Containers image artifacts.
#IncusItem: {
	ftype!:                          #IncusFileType @go(FileType)
	path!:                           #RelativePath  @go(,type=RelativePath)
	size!:                           #Size
	sha256!:                         #Checksum @go(SHA256)
	delta_base?:                     string    @go(DeltaBase)
	combined_sha256?:                #Checksum @go(CombinedSHA256)
	combined_rootxz_sha256?:         #Checksum @go(CombinedRootXZSHA256)
	combined_squashfs_sha256?:       #Checksum @go(CombinedSquashfsSHA256)
	"combined_disk-kvm-img_sha256"?: #Checksum @go(CombinedDiskKVMImageSHA256)
}

// IncusFileType is the current Incus/Linux Containers artifact vocabulary.
#IncusFileType: "incus.tar.xz" | "lxd.tar.xz" | "root.tar.xz" | "squashfs" | "disk-kvm.img" | "squashfs.vcdiff"

// LXDProductFile is the Canonical LXD image server profile.
#LXDProductFile: {
	format!:     #ProductsFormat
	content_id!: "images"          @go(ContentID)
	datatype!:   "image-downloads" @go(DataType)
	updated?:    #Timestamp
	products!: [string]: #LXDProduct
}

// LXDProduct describes Canonical LXD image product metadata.
#LXDProduct: {
	aliases!:       string
	arch!:          string
	distro!:        string
	os!:            string @go(OS)
	release!:       string
	release_title!: string @go(ReleaseTitle)
	variant!:       string
	requirements!:  #LXDRequirements
	versions!: [string]: #LXDVersion
}

// LXDRequirements allows the narrower current Canonical LXD requirement vocabulary.
#LXDRequirements: {
	cgroup?:     "v1" | "v2"
	secureboot?: #StringBool @go(SecureBoot)
}

// LXDVersion describes one Canonical LXD image build.
#LXDVersion: {
	items!: [string]: #LXDItem
}

// LXDItem describes Canonical LXD image artifacts.
#LXDItem: {
	ftype!:                          #LXDFileType  @go(FileType)
	path!:                           #RelativePath @go(,type=RelativePath)
	size!:                           #Size
	sha256!:                         #Checksum @go(SHA256)
	delta_base?:                     string    @go(DeltaBase)
	combined_squashfs_sha256?:       #Checksum @go(CombinedSquashfsSHA256)
	"combined_disk-kvm-img_sha256"?: #Checksum @go(CombinedDiskKVMImageSHA256)
}

// LXDFileType is the current Canonical LXD artifact vocabulary.
#LXDFileType: "lxd.tar.xz" | "squashfs" | "disk-kvm.img" | "disk-kvm.img.vcdiff" | "squashfs.vcdiff"

// StringBool is a boolean encoded as a JSON string, matching current image streams.
#StringBool: "true" | "false"
