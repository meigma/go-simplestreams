package simplestreams

import (
	"errors"
	"maps"
	"reflect"
	"slices"
	"sort"
	"strings"
)

// ItemRef identifies one item in a Simple Streams product tree.
type ItemRef struct {
	// ContentID is the product file content ID.
	ContentID string

	// ProductName is the product map key.
	ProductName string

	// VersionName is the version map key.
	VersionName string

	// ItemName is the item map key.
	ItemName string
}

// String returns the canonical slash-delimited item identity.
func (ref ItemRef) String() string {
	return strings.Join([]string{ref.ContentID, ref.ProductName, ref.VersionName, ref.ItemName}, "/")
}

// ItemView is an effective item view with inherited metadata applied.
type ItemView struct {
	// Ref identifies the item in the product tree.
	Ref ItemRef

	// ProductFile is the product document that contains the item.
	ProductFile *ProductFile

	// Product is the parent product.
	Product *Product

	// Version is the parent version.
	Version *Version

	// Item is the item data.
	Item *Item

	// Metadata is the effective metadata for the item.
	Metadata map[string]any
}

// WalkItems calls visit once for each item in deterministic product/version/item order.
func (productFile *ProductFile) WalkItems(visit func(ItemView) error) error {
	if productFile == nil {
		return errors.New("simplestreams: nil product file")
	}
	if visit == nil {
		return errors.New("simplestreams: nil item visitor")
	}
	for _, productName := range sortedKeys(productFile.Products) {
		product := productFile.Products[productName]
		for _, versionName := range sortedKeys(product.Versions) {
			version := product.Versions[versionName]
			for _, itemName := range sortedKeys(version.Items) {
				view := makeItemView(productFile, product, version, version.Items[itemName])
				if err := visit(view); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

// Items returns all items in deterministic product/version/item order.
func (productFile *ProductFile) Items() []ItemView {
	if productFile == nil {
		return nil
	}
	items := []ItemView{}
	_ = productFile.WalkItems(func(item ItemView) error {
		items = append(items, item)
		return nil
	})
	return items
}

// Artifact returns the item's primary artifact reference when the item has a path.
func (view ItemView) Artifact() (ArtifactRef, bool) {
	if view.Item == nil || view.Item.Path == "" {
		return ArtifactRef{}, false
	}
	return ArtifactRef{
		Path:      view.Item.Path,
		Mirrors:   slices.Clone(view.Item.Mirrors),
		Size:      cloneInt64(view.Item.Size),
		Checksums: Checksums(view.Item),
	}, true
}

func makeItemView(productFile *ProductFile, product *Product, version *Version, item *Item) ItemView {
	ref := ItemRef{
		ContentID:   productFile.ContentID,
		ProductName: product.Name,
		VersionName: version.Name,
		ItemName:    item.Name,
	}
	metadata := effectiveMetadata(productFile, product, version, item, ref)
	return ItemView{
		Ref:         ref,
		ProductFile: productFile,
		Product:     product,
		Version:     version,
		Item:        item,
		Metadata:    metadata,
	}
}

func effectiveMetadata(
	productFile *ProductFile,
	product *Product,
	version *Version,
	item *Item,
	ref ItemRef,
) map[string]any {
	metadata := map[string]any{}

	metadata["format"] = productFile.Format
	metadata["content_id"] = productFile.ContentID
	if productFile.DataType != "" {
		metadata["datatype"] = productFile.DataType
	}
	if productFile.Updated != "" {
		metadata["updated"] = productFile.Updated
	}
	copyMetadata(metadata, productFile.Metadata)
	copyMetadata(metadata, product.Metadata)
	copyMetadata(metadata, version.Metadata)

	if item.FileType != "" {
		metadata["ftype"] = item.FileType
	}
	if item.Path != "" {
		metadata["path"] = item.Path.String()
	}
	if item.Size != nil {
		metadata["size"] = *item.Size
	}
	if item.MD5 != "" {
		metadata["md5"] = item.MD5
	}
	if item.SHA256 != "" {
		metadata["sha256"] = item.SHA256
	}
	if item.SHA512 != "" {
		metadata["sha512"] = item.SHA512
	}
	if len(item.Mirrors) > 0 {
		mirrors := make([]string, len(item.Mirrors))
		for i, mirror := range item.Mirrors {
			mirrors[i] = string(mirror)
		}
		metadata["mirrors"] = mirrors
	}
	copyMetadata(metadata, item.Metadata)

	metadata["product_name"] = ref.ProductName
	metadata["version_name"] = ref.VersionName
	metadata["item_name"] = ref.ItemName
	return metadata
}

func copyMetadata(dst map[string]any, src map[string]any) {
	maps.Copy(dst, src)
}

func sortedKeys[V any](values map[string]V) []string {
	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

func cloneInt64(value *int64) *int64 {
	if value == nil {
		return nil
	}
	clone := *value
	return &clone
}

// Predicate decides whether an item should be selected.
type Predicate func(ItemView) bool

// FilterItems returns items accepted by predicate in their original order.
func FilterItems(items []ItemView, predicate Predicate) []ItemView {
	if predicate == nil {
		return slices.Clone(items)
	}
	filtered := make([]ItemView, 0, len(items))
	for _, item := range items {
		if predicate(item) {
			filtered = append(filtered, item)
		}
	}
	return filtered
}

// And returns a predicate that requires all predicates to match.
func And(predicates ...Predicate) Predicate {
	return func(item ItemView) bool {
		for _, predicate := range predicates {
			if predicate != nil && !predicate(item) {
				return false
			}
		}
		return true
	}
}

// MatchDataType matches items from product files with datatype.
func MatchDataType(datatype string) Predicate {
	return func(item ItemView) bool {
		return item.ProductFile != nil && item.ProductFile.DataType == datatype
	}
}

// MatchProductName matches items with the provided product name.
func MatchProductName(name string) Predicate {
	return func(item ItemView) bool {
		return item.Ref.ProductName == name
	}
}

// MatchVersionName matches items with the provided version name.
func MatchVersionName(name string) Predicate {
	return func(item ItemView) bool {
		return item.Ref.VersionName == name
	}
}

// MatchItemName matches items with the provided item name.
func MatchItemName(name string) Predicate {
	return func(item ItemView) bool {
		return item.Ref.ItemName == name
	}
}

// MatchFileType matches items with the provided ftype value.
func MatchFileType(fileType string) Predicate {
	return func(item ItemView) bool {
		return item.Item != nil && item.Item.FileType == fileType
	}
}

// HasPath matches items that reference a primary artifact path.
func HasPath() Predicate {
	return func(item ItemView) bool {
		return item.Item != nil && item.Item.Path != ""
	}
}

// MetadataEquals matches items whose effective metadata key equals value.
func MetadataEquals(key string, value any) Predicate {
	return func(item ItemView) bool {
		return reflect.DeepEqual(item.Metadata[key], value)
	}
}

// LatestByVersion returns the item with the greatest bytewise lexical version name.
func LatestByVersion(items []ItemView) (ItemView, bool) {
	if len(items) == 0 {
		return ItemView{}, false
	}
	latest := items[0]
	for _, item := range items[1:] {
		if item.Ref.VersionName > latest.Ref.VersionName {
			latest = item
		}
	}
	return latest, true
}
