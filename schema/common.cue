package schema

// RelativePath is a Simple Streams path relative to a stream root.
#RelativePath: string & !="" & !~"^/"

// MetadataPath is a Simple Streams metadata path ending in a supported file format.
#MetadataPath: #RelativePath & (=~"\\.json$" | =~"\\.sjson$")

// Timestamp is a Simple Streams timestamp string.
#Timestamp: string

// Size is the declared byte size of an artifact item.
#Size: int

// Checksum is an artifact checksum value.
#Checksum: string

// Mirror is an alternate mirror prefix for an item path.
#Mirror: string
