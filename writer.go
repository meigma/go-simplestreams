package simplestreams

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"slices"
)

// BuildIndexEntry is one entry used to build an index document.
type BuildIndexEntry struct {
	// ContentID is the index map key for the entry.
	ContentID string

	// Path is the product or nested index metadata path.
	Path RelativePath

	// Format identifies the referenced document format when provided.
	Format string

	// DataType describes the referenced content domain when provided.
	DataType string

	// Updated is the timestamp for the referenced document.
	Updated string

	// Products lists product names available through the referenced document.
	Products []string
}

// BuildIndex creates an index document from entries.
func BuildIndex(entries []BuildIndexEntry, updated string) (*Index, error) {
	index := &Index{
		Format:  indexFormat,
		Updated: updated,
		Entries: make(map[string]*IndexEntry, len(entries)),
	}
	for _, entry := range entries {
		if entry.ContentID == "" {
			return nil, errors.New("simplestreams: index entry missing content ID")
		}
		if _, exists := index.Entries[entry.ContentID]; exists {
			return nil, fmt.Errorf("simplestreams: duplicate index content ID %q", entry.ContentID)
		}
		if err := entry.Path.Validate(); err != nil {
			return nil, err
		}
		if !entry.Path.IsMetadataPath() {
			return nil, fmt.Errorf(
				"simplestreams: index entry %q path %q is not .json or .sjson",
				entry.ContentID,
				entry.Path,
			)
		}
		index.Entries[entry.ContentID] = &IndexEntry{
			ContentID: entry.ContentID,
			Format:    entry.Format,
			DataType:  entry.DataType,
			Path:      entry.Path,
			Updated:   entry.Updated,
			Products:  slices.Clone(entry.Products),
			parent:    index,
		}
	}
	return index, nil
}

// CheckDuplicateItemRefs rejects repeated item identities.
func CheckDuplicateItemRefs(refs []ItemRef) error {
	seen := make(map[ItemRef]struct{}, len(refs))
	for _, ref := range refs {
		if _, ok := seen[ref]; ok {
			return fmt.Errorf("simplestreams: duplicate item identity %s", ref)
		}
		seen[ref] = struct{}{}
	}
	return nil
}

// MarshalJSONDocument renders v as deterministic indented JSON with a trailing newline.
func MarshalJSONDocument(v any) ([]byte, error) {
	switch v.(type) {
	case *Index, *ProductFile:
	default:
		return nil, fmt.Errorf("simplestreams: unsupported JSON document type %T", v)
	}
	body, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return nil, err
	}
	var out bytes.Buffer
	out.Write(body)
	out.WriteByte('\n')
	return out.Bytes(), nil
}
