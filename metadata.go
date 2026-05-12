package simplestreams

// SetMetadata sets an unknown product-file metadata field.
func (productFile *ProductFile) SetMetadata(key string, value any) {
	productFile.Metadata = setMetadata(productFile.Metadata, key, value)
}

// MetadataValue returns an unknown product-file metadata field.
func (productFile *ProductFile) MetadataValue(key string) (any, bool) {
	if productFile == nil {
		return nil, false
	}
	return metadataValue(productFile.Metadata, key)
}

// SetMetadata sets an unknown product metadata field.
func (product *Product) SetMetadata(key string, value any) {
	product.Metadata = setMetadata(product.Metadata, key, value)
}

// MetadataValue returns an unknown product metadata field.
func (product *Product) MetadataValue(key string) (any, bool) {
	if product == nil {
		return nil, false
	}
	return metadataValue(product.Metadata, key)
}

// SetMetadata sets an unknown version metadata field.
func (version *Version) SetMetadata(key string, value any) {
	version.Metadata = setMetadata(version.Metadata, key, value)
}

// MetadataValue returns an unknown version metadata field.
func (version *Version) MetadataValue(key string) (any, bool) {
	if version == nil {
		return nil, false
	}
	return metadataValue(version.Metadata, key)
}

// SetMetadata sets an unknown item metadata field.
func (item *Item) SetMetadata(key string, value any) {
	item.Metadata = setMetadata(item.Metadata, key, value)
}

// MetadataValue returns an unknown item metadata field.
func (item *Item) MetadataValue(key string) (any, bool) {
	if item == nil {
		return nil, false
	}
	return metadataValue(item.Metadata, key)
}

func setMetadata(metadata map[string]any, key string, value any) map[string]any {
	if metadata == nil {
		metadata = map[string]any{}
	}
	metadata[key] = value
	return metadata
}

func metadataValue(metadata map[string]any, key string) (any, bool) {
	if metadata == nil {
		return nil, false
	}
	value, ok := metadata[key]
	return value, ok
}
