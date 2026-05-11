package schema

// IndexFile is a Simple Streams index document.
#IndexFile: {
	// Format identifies this document as a Simple Streams index.
	format!: #IndexFormat

	// Updated is the timestamp for the index document.
	updated?: #Timestamp

	// Index maps content IDs to product documents or nested indexes.
	index!: [string]: #IndexEntry
}

// IndexEntry describes one entry in a Simple Streams index document.
#IndexEntry: {
	// Format identifies the referenced document format when provided.
	format?: #DocumentFormat

	// DataType describes the referenced content domain when provided.
	datatype?: #DataType @go(DataType)

	// Path is the stream-root-relative path to the referenced document.
	path!: #MetadataPath

	// Updated is the timestamp for the referenced document.
	updated?: #Timestamp

	// Products lists product names available through the referenced document.
	products?: [...string]
}
