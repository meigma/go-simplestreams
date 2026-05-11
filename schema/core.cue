package schema

// ProductFileCore is an open Simple Streams product document shape for profiles.
#ProductFileCore: {
	@go(-)

	format!:     #ProductsFormat
	content_id!: string
	datatype?:   #DataType
	updated?:    #Timestamp
	products!: [string]: #ProductCore
	...
}

// ProductCore is an open product shape for profile-specific metadata.
#ProductCore: {
	@go(-)

	versions!: [string]: #VersionCore
	...
}

// VersionCore is an open version shape for profile-specific metadata.
#VersionCore: {
	@go(-)

	items!: [string]: #ItemCore
	...
}

// ItemCore is an open item shape for profile-specific artifact metadata.
#ItemCore: {
	@go(-)

	ftype?:  string
	path?:   #RelativePath
	size?:   #Size
	md5?:    #Checksum
	sha256?: #Checksum
	sha512?: #Checksum
	mirrors?: [...#Mirror]

	if mirrors != _|_ {
		path!: #RelativePath
	}

	...
}
