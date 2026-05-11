package schema

// IndexFormat is the format marker for a Simple Streams index document.
#IndexFormat: "index:1.0"

// ProductsFormat is the format marker for a Simple Streams product document.
#ProductsFormat: "products:1.0"

// DocumentFormat is any Simple Streams document format marker modeled by this package.
#DocumentFormat: #IndexFormat | #ProductsFormat

// DataType describes the kind of content referenced by an index or product document.
#DataType: string
