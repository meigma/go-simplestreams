package simplestreams

// NewProductFile creates a product document with the Simple Streams products format.
func NewProductFile(contentID string) *ProductFile {
	return &ProductFile{
		Format:    ProductsFormat,
		ContentID: contentID,
		Products:  map[string]*Product{},
	}
}

// NewProduct creates a product with an initialized version map.
func NewProduct(name string) *Product {
	return &Product{
		Name:     name,
		Versions: map[string]*Version{},
	}
}

// NewVersion creates a product version with an initialized item map.
func NewVersion(name string) *Version {
	return &Version{
		Name:  name,
		Items: map[string]*Item{},
	}
}

// NewItem creates an item with its map key name recorded.
func NewItem(name string) *Item {
	return &Item{Name: name}
}

// SetProduct stores product under name and updates the product name and parent.
//
// If product is nil, SetProduct creates one. The product map is initialized
// when needed.
func (productFile *ProductFile) SetProduct(name string, product *Product) *Product {
	if product == nil {
		product = NewProduct(name)
	}
	product.Name = name
	if product.Versions == nil {
		product.Versions = map[string]*Version{}
	}
	product.parent = productFile

	if productFile.Products == nil {
		productFile.Products = map[string]*Product{}
	}
	productFile.Products[name] = product
	return product
}

// SetVersion stores version under name and updates the version name and parent.
//
// If version is nil, SetVersion creates one. The version map is initialized
// when needed.
func (product *Product) SetVersion(name string, version *Version) *Version {
	if version == nil {
		version = NewVersion(name)
	}
	version.Name = name
	if version.Items == nil {
		version.Items = map[string]*Item{}
	}
	version.parent = product

	if product.Versions == nil {
		product.Versions = map[string]*Version{}
	}
	product.Versions[name] = version
	return version
}

// SetItem stores item under name and updates the item name and parent.
//
// If item is nil, SetItem creates one. The item map is initialized when needed.
func (version *Version) SetItem(name string, item *Item) *Item {
	if item == nil {
		item = NewItem(name)
	}
	item.Name = name
	item.parent = version

	if version.Items == nil {
		version.Items = map[string]*Item{}
	}
	version.Items[name] = item
	return item
}
