package schema

// ProductFile is a Simple Streams product document.
#ProductFile: {
	// Format identifies this document as a Simple Streams product document.
	format!: #ProductsFormat

	// ContentID identifies this product document.
	content_id!: string @go(ContentID)

	// DataType describes the product content domain when provided.
	datatype?: #DataType @go(DataType)

	// Updated is the timestamp for the product document.
	updated?: #Timestamp

	// Products maps product names to product metadata and versions.
	products!: [string]: #Product
}

// Product describes one product in a Simple Streams product document.
#Product: {
	// Versions maps version names to version metadata and items.
	versions!: [string]: #Version
}

// Version describes one product version in a Simple Streams product document.
#Version: {
	// Items maps item names to item metadata.
	items!: [string]: #Item
}

// Item describes one item in a Simple Streams product version.
#Item: {
	// FileType identifies the item file type when provided.
	ftype?: string @go(FileType)

	// Path is the stream-root-relative path to the artifact when provided.
	path?: #RelativePath @go(,type=RelativePath)

	// Size is the declared byte size of the artifact when provided.
	size?: #Size

	// MD5 is the artifact MD5 checksum when provided.
	md5?: #Checksum @go(MD5)

	// SHA256 is the artifact SHA-256 checksum when provided.
	sha256?: #Checksum @go(SHA256)

	// SHA512 is the artifact SHA-512 checksum when provided.
	sha512?: #Checksum @go(SHA512)

	// Mirrors lists alternate mirror prefixes for the artifact path.
	mirrors?: [...#Mirror]

	if mirrors != _|_ {
		path!: #RelativePath
	}
}
