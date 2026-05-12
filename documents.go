package simplestreams

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"slices"
)

const (
	indexFormat   = "index:1.0"
	productFormat = "products:1.0"
)

// Index is a Simple Streams index document.
type Index struct {
	// Format identifies this document as a Simple Streams index.
	Format string

	// Updated is the timestamp for the index document.
	Updated string

	// Metadata preserves index-level fields not modeled directly by this type.
	Metadata map[string]any

	// Entries maps content IDs to product documents or nested indexes.
	Entries map[string]*IndexEntry

	mirror *Mirror
	path   RelativePath
}

// IndexEntry describes one entry in a Simple Streams index document.
type IndexEntry struct {
	// ContentID is the index map key for this entry.
	ContentID string

	// Format identifies the referenced document format when provided.
	Format string

	// DataType describes the referenced content domain when provided.
	DataType string

	// Path is the stream-root-relative path to the referenced document.
	Path RelativePath

	// Updated is the timestamp for the referenced document.
	Updated string

	// Products lists product names available through the referenced document.
	Products []string

	// Metadata preserves index-entry fields not modeled directly by this type.
	Metadata map[string]any

	parent *Index
}

// ProductFile loads the product document referenced by entry.
func (entry *IndexEntry) ProductFile(ctx context.Context) (*ProductFile, error) {
	if entry == nil || entry.parent == nil || entry.parent.mirror == nil {
		return nil, errors.New("simplestreams: index entry is not attached to a mirror")
	}
	return entry.parent.mirror.ProductFile(ctx, entry)
}

// Index loads the nested index document referenced by entry.
func (entry *IndexEntry) Index(ctx context.Context) (*Index, error) {
	if entry == nil || entry.parent == nil || entry.parent.mirror == nil {
		return nil, errors.New("simplestreams: index entry is not attached to a mirror")
	}
	return entry.parent.mirror.IndexAt(ctx, entry.Path)
}

// ProductFile is a Simple Streams product document.
type ProductFile struct {
	// Format identifies this document as a Simple Streams product document.
	Format string

	// ContentID identifies this product document.
	ContentID string

	// DataType describes the product content domain when provided.
	DataType string

	// Updated is the timestamp for the product document.
	Updated string

	// Metadata preserves fields not modeled directly by this type.
	Metadata map[string]any

	// Products maps product names to product metadata and versions.
	Products map[string]*Product

	mirror *Mirror
	path   RelativePath
}

// Product describes one product in a Simple Streams product document.
type Product struct {
	// Name is the product map key.
	Name string

	// Metadata preserves product-level metadata fields.
	Metadata map[string]any

	// Versions maps version names to version metadata and items.
	Versions map[string]*Version

	parent *ProductFile
}

// Version describes one product version in a Simple Streams product document.
type Version struct {
	// Name is the version map key.
	Name string

	// Metadata preserves version-level metadata fields.
	Metadata map[string]any

	// Items maps item names to item metadata.
	Items map[string]*Item

	parent *Product
}

// Item describes one item in a Simple Streams product version.
type Item struct {
	// Name is the item map key.
	Name string

	// FileType identifies the item file type when provided.
	FileType string

	// Path is the stream-root-relative artifact path when provided.
	Path RelativePath

	// Size is the declared artifact byte size when provided.
	Size *int64

	// MD5 is the artifact MD5 checksum when provided.
	MD5 string

	// SHA256 is the artifact SHA-256 checksum when provided.
	SHA256 string

	// SHA512 is the artifact SHA-512 checksum when provided.
	SHA512 string

	// Mirrors lists alternate mirror prefixes for the artifact path.
	Mirrors []MirrorPrefix

	// Metadata preserves item-level metadata fields.
	Metadata map[string]any

	parent *Version
}

// MirrorPrefix is an alternate mirror prefix for an item path.
type MirrorPrefix string

type indexJSON struct {
	Format  string                     `json:"format"`
	Updated string                     `json:"updated,omitempty"`
	Index   map[string]json.RawMessage `json:"index"`
}

// MarshalJSON renders an index document using the Simple Streams wire shape.
func (index *Index) MarshalJSON() ([]byte, error) {
	if index == nil {
		return []byte("null"), nil
	}
	raw := map[string]any{}
	copyMetadata(raw, index.Metadata)
	raw["format"] = index.Format
	if index.Updated != "" {
		raw["updated"] = index.Updated
	}
	entries := make(map[string]any, len(index.Entries))
	for contentID, entry := range index.Entries {
		if entry == nil {
			return nil, fmt.Errorf("simplestreams: nil index entry %q", contentID)
		}
		entries[contentID] = entry.marshalJSON()
	}
	raw["index"] = entries
	return json.Marshal(raw)
}

type indexEntryJSON struct {
	Format   string   `json:"format,omitempty"`
	DataType string   `json:"datatype,omitempty"`
	Path     string   `json:"path"`
	Updated  string   `json:"updated,omitempty"`
	Products []string `json:"products,omitempty"`
}

// MarshalJSON renders an index entry using the Simple Streams wire shape.
func (entry *IndexEntry) MarshalJSON() ([]byte, error) {
	if entry == nil {
		return []byte("null"), nil
	}
	return json.Marshal(entry.marshalJSON())
}

func (entry *IndexEntry) marshalJSON() map[string]any {
	raw := map[string]any{}
	copyMetadata(raw, entry.Metadata)
	if entry.Format != "" {
		raw["format"] = entry.Format
	}
	if entry.DataType != "" {
		raw["datatype"] = entry.DataType
	}
	raw["path"] = entry.Path.String()
	if entry.Updated != "" {
		raw["updated"] = entry.Updated
	}
	if len(entry.Products) > 0 {
		raw["products"] = slices.Clone(entry.Products)
	}
	return raw
}

type productFileJSON struct {
	Format    string                     `json:"format"`
	ContentID string                     `json:"content_id"`
	DataType  string                     `json:"datatype,omitempty"`
	Updated   string                     `json:"updated,omitempty"`
	Products  map[string]json.RawMessage `json:"products"`
}

// MarshalJSON renders a product file using the Simple Streams wire shape.
func (productFile *ProductFile) MarshalJSON() ([]byte, error) {
	if productFile == nil {
		return []byte("null"), nil
	}
	raw, err := productFile.marshalJSON()
	if err != nil {
		return nil, err
	}
	return json.Marshal(raw)
}

func (productFile *ProductFile) marshalJSON() (map[string]any, error) {
	raw := map[string]any{}
	copyMetadata(raw, productFile.Metadata)
	raw["format"] = productFile.Format
	raw["content_id"] = productFile.ContentID
	if productFile.DataType != "" {
		raw["datatype"] = productFile.DataType
	}
	if productFile.Updated != "" {
		raw["updated"] = productFile.Updated
	}

	products := make(map[string]any, len(productFile.Products))
	for name, product := range productFile.Products {
		if product == nil {
			return nil, fmt.Errorf("simplestreams: nil product %q", name)
		}
		productRaw, err := product.marshalJSON()
		if err != nil {
			return nil, err
		}
		products[name] = productRaw
	}
	raw["products"] = products
	return raw, nil
}

type productJSON struct {
	Versions map[string]json.RawMessage `json:"versions"`
}

// MarshalJSON renders a product using the Simple Streams wire shape.
func (product *Product) MarshalJSON() ([]byte, error) {
	if product == nil {
		return []byte("null"), nil
	}
	raw, err := product.marshalJSON()
	if err != nil {
		return nil, err
	}
	return json.Marshal(raw)
}

func (product *Product) marshalJSON() (map[string]any, error) {
	raw := map[string]any{}
	copyMetadata(raw, product.Metadata)

	versions := make(map[string]any, len(product.Versions))
	for name, version := range product.Versions {
		if version == nil {
			return nil, fmt.Errorf("simplestreams: nil version %q/%q", product.Name, name)
		}
		versionRaw, err := version.marshalJSON()
		if err != nil {
			return nil, err
		}
		versions[name] = versionRaw
	}
	raw["versions"] = versions
	return raw, nil
}

type versionJSON struct {
	Items map[string]json.RawMessage `json:"items"`
}

// MarshalJSON renders a version using the Simple Streams wire shape.
func (version *Version) MarshalJSON() ([]byte, error) {
	if version == nil {
		return []byte("null"), nil
	}
	raw, err := version.marshalJSON()
	if err != nil {
		return nil, err
	}
	return json.Marshal(raw)
}

func (version *Version) marshalJSON() (map[string]any, error) {
	raw := map[string]any{}
	copyMetadata(raw, version.Metadata)

	items := make(map[string]any, len(version.Items))
	for name, item := range version.Items {
		if item == nil {
			return nil, fmt.Errorf("simplestreams: nil item %q/%q", version.Name, name)
		}
		items[name] = item.marshalJSON()
	}
	raw["items"] = items
	return raw, nil
}

type itemJSON struct {
	FileType string         `json:"ftype,omitempty"`
	Path     string         `json:"path,omitempty"`
	Size     *int64         `json:"size,omitempty"`
	MD5      string         `json:"md5,omitempty"`
	SHA256   string         `json:"sha256,omitempty"`
	SHA512   string         `json:"sha512,omitempty"`
	Mirrors  []MirrorPrefix `json:"mirrors,omitempty"`
}

// MarshalJSON renders an item using the Simple Streams wire shape.
func (item *Item) MarshalJSON() ([]byte, error) {
	if item == nil {
		return []byte("null"), nil
	}
	return json.Marshal(item.marshalJSON())
}

func (item *Item) marshalJSON() map[string]any {
	raw := map[string]any{}
	copyMetadata(raw, item.Metadata)
	if item.FileType != "" {
		raw["ftype"] = item.FileType
	}
	if item.Path != "" {
		raw["path"] = item.Path.String()
	}
	if item.Size != nil {
		raw["size"] = *item.Size
	}
	if item.MD5 != "" {
		raw["md5"] = item.MD5
	}
	if item.SHA256 != "" {
		raw["sha256"] = item.SHA256
	}
	if item.SHA512 != "" {
		raw["sha512"] = item.SHA512
	}
	if len(item.Mirrors) > 0 {
		raw["mirrors"] = slices.Clone(item.Mirrors)
	}
	return raw
}

func decodeIndex(data []byte, mirror *Mirror, path RelativePath) (*Index, error) {
	var raw indexJSON
	if err := decodeStrict(data, &raw); err != nil {
		return nil, fmt.Errorf("simplestreams: decode index %q: %w", path, err)
	}
	if raw.Format != indexFormat {
		return nil, fmt.Errorf("simplestreams: index %q has format %q, want %q", path, raw.Format, indexFormat)
	}
	if raw.Index == nil {
		return nil, fmt.Errorf("simplestreams: index %q missing index entries", path)
	}

	metadata, err := unknownMetadata(data, "format", "updated", "index")
	if err != nil {
		return nil, err
	}
	index := &Index{
		Format:   raw.Format,
		Updated:  raw.Updated,
		Metadata: metadata,
		Entries:  make(map[string]*IndexEntry, len(raw.Index)),
		mirror:   mirror,
		path:     path,
	}
	for contentID, entryData := range raw.Index {
		var entryRaw indexEntryJSON
		if err := decodeStrict(entryData, &entryRaw); err != nil {
			return nil, fmt.Errorf("simplestreams: decode index entry %q: %w", contentID, err)
		}
		entryPath, err := ParseRelativePath(entryRaw.Path)
		if err != nil {
			return nil, fmt.Errorf("simplestreams: index entry %q path: %w", contentID, err)
		}
		if !entryPath.IsMetadataPath() {
			return nil, fmt.Errorf("simplestreams: index entry %q path %q is not .json or .sjson", contentID, entryPath)
		}
		metadata, err := unknownMetadata(entryData, "format", "datatype", "path", "updated", "products")
		if err != nil {
			return nil, fmt.Errorf("simplestreams: decode index entry %q metadata: %w", contentID, err)
		}
		index.Entries[contentID] = &IndexEntry{
			ContentID: contentID,
			Format:    entryRaw.Format,
			DataType:  entryRaw.DataType,
			Path:      entryPath,
			Updated:   entryRaw.Updated,
			Products:  slices.Clone(entryRaw.Products),
			Metadata:  metadata,
			parent:    index,
		}
	}
	return index, nil
}

func decodeProductFile(data []byte, mirror *Mirror, path RelativePath) (*ProductFile, error) {
	var raw productFileJSON
	if err := decodeStrict(data, &raw); err != nil {
		return nil, fmt.Errorf("simplestreams: decode product file %q: %w", path, err)
	}
	if raw.Format != productFormat {
		return nil, fmt.Errorf("simplestreams: product file %q has format %q, want %q", path, raw.Format, productFormat)
	}
	if raw.ContentID == "" {
		return nil, fmt.Errorf("simplestreams: product file %q missing content_id", path)
	}
	if raw.Products == nil {
		return nil, fmt.Errorf("simplestreams: product file %q missing products", path)
	}

	metadata, err := unknownMetadata(data, "format", "content_id", "datatype", "updated", "products")
	if err != nil {
		return nil, err
	}
	productFile := &ProductFile{
		Format:    raw.Format,
		ContentID: raw.ContentID,
		DataType:  raw.DataType,
		Updated:   raw.Updated,
		Metadata:  metadata,
		Products:  make(map[string]*Product, len(raw.Products)),
		mirror:    mirror,
		path:      path,
	}
	for name, productData := range raw.Products {
		product, err := decodeProduct(name, productData, productFile)
		if err != nil {
			return nil, err
		}
		productFile.Products[name] = product
	}
	return productFile, nil
}

func decodeProduct(name string, data []byte, parent *ProductFile) (*Product, error) {
	var raw productJSON
	if err := decodeStrict(data, &raw); err != nil {
		return nil, fmt.Errorf("simplestreams: decode product %q: %w", name, err)
	}
	if raw.Versions == nil {
		return nil, fmt.Errorf("simplestreams: product %q missing versions", name)
	}
	metadata, err := unknownMetadata(data, "versions")
	if err != nil {
		return nil, err
	}
	product := &Product{
		Name:     name,
		Metadata: metadata,
		Versions: make(map[string]*Version, len(raw.Versions)),
		parent:   parent,
	}
	for versionName, versionData := range raw.Versions {
		version, err := decodeVersion(versionName, versionData, product)
		if err != nil {
			return nil, err
		}
		product.Versions[versionName] = version
	}
	return product, nil
}

func decodeVersion(name string, data []byte, parent *Product) (*Version, error) {
	var raw versionJSON
	if err := decodeStrict(data, &raw); err != nil {
		return nil, fmt.Errorf("simplestreams: decode version %q/%q: %w", parent.Name, name, err)
	}
	if raw.Items == nil {
		return nil, fmt.Errorf("simplestreams: version %q/%q missing items", parent.Name, name)
	}
	metadata, err := unknownMetadata(data, "items")
	if err != nil {
		return nil, err
	}
	version := &Version{
		Name:     name,
		Metadata: metadata,
		Items:    make(map[string]*Item, len(raw.Items)),
		parent:   parent,
	}
	for itemName, itemData := range raw.Items {
		item, err := decodeItem(itemName, itemData, version)
		if err != nil {
			return nil, err
		}
		version.Items[itemName] = item
	}
	return version, nil
}

func decodeItem(name string, data []byte, parent *Version) (*Item, error) {
	var raw itemJSON
	if err := decodeStrict(data, &raw); err != nil {
		return nil, fmt.Errorf("simplestreams: decode item %q/%q/%q: %w", parent.parent.Name, parent.Name, name, err)
	}
	var itemPath RelativePath
	if raw.Path != "" {
		var err error
		itemPath, err = ParseRelativePath(raw.Path)
		if err != nil {
			return nil, fmt.Errorf("simplestreams: item %q/%q/%q path: %w", parent.parent.Name, parent.Name, name, err)
		}
	}
	metadata, err := unknownMetadata(data, "ftype", "path", "size", "md5", "sha256", "sha512", "mirrors")
	if err != nil {
		return nil, err
	}
	return &Item{
		Name:     name,
		FileType: raw.FileType,
		Path:     itemPath,
		Size:     raw.Size,
		MD5:      raw.MD5,
		SHA256:   raw.SHA256,
		SHA512:   raw.SHA512,
		Mirrors:  slices.Clone(raw.Mirrors),
		Metadata: metadata,
		parent:   parent,
	}, nil
}

func decodeStrict(data []byte, target any) error {
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.UseNumber()
	if err := decoder.Decode(target); err != nil {
		return err
	}
	if err := decoder.Decode(&struct{}{}); err != io.EOF {
		return errors.New("multiple JSON values")
	}
	return nil
}

func unknownMetadata(data []byte, known ...string) (map[string]any, error) {
	var raw map[string]json.RawMessage
	if err := decodeStrict(data, &raw); err != nil {
		return nil, err
	}
	knownSet := map[string]struct{}{}
	for _, field := range known {
		knownSet[field] = struct{}{}
	}
	metadata := make(map[string]any)
	for key, valueData := range raw {
		if _, ok := knownSet[key]; ok {
			continue
		}
		var value any
		if err := decodeStrict(valueData, &value); err != nil {
			return nil, fmt.Errorf("decode metadata field %q: %w", key, err)
		}
		metadata[key] = value
	}
	return metadata, nil
}
