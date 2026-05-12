package simplestreams

import (
	"context"
	"errors"
	"fmt"
	"io"
	"sync"
)

// Mirror is a lazy, backend-neutral view of a Simple Streams mirror.
type Mirror struct {
	source Source

	indexPath              RelativePath
	indexPathExplicit      bool
	allowContentIDMismatch bool
	signaturePolicy        SignaturePolicy
	verifier               Verifier

	mu       sync.Mutex
	indexes  map[RelativePath]*Index
	products map[RelativePath]*ProductFile
}

// NewMirror constructs a Mirror backed by source.
func NewMirror(source Source, options ...Option) (*Mirror, error) {
	if source == nil {
		return nil, errors.New("simplestreams: nil source")
	}

	config := defaultMirrorOptions()
	for _, option := range options {
		if option == nil {
			continue
		}
		if err := option(&config); err != nil {
			return nil, err
		}
	}

	return &Mirror{
		source:                 source,
		indexPath:              config.indexPath,
		indexPathExplicit:      config.indexPathExplicit,
		allowContentIDMismatch: config.allowContentIDMismatch,
		signaturePolicy:        config.signaturePolicy,
		verifier:               config.verifier,
		indexes:                map[RelativePath]*Index{},
		products:               map[RelativePath]*ProductFile{},
	}, nil
}

// Index loads and returns the mirror index.
func (mirror *Mirror) Index(ctx context.Context) (*Index, error) {
	if mirror.indexPathExplicit {
		return mirror.IndexAt(ctx, mirror.indexPath)
	}
	switch mirror.signaturePolicy {
	case SignatureDisabled:
		return mirror.IndexAt(ctx, DefaultIndexPath)
	case SignatureRequired:
		return mirror.IndexAt(ctx, DefaultSignedIndexPath)
	case SignaturePreferred:
		index, err := mirror.IndexAt(ctx, DefaultSignedIndexPath)
		if err == nil {
			return index, nil
		}
		if errors.Is(err, ErrNotFound) {
			return mirror.IndexAt(ctx, DefaultIndexPath)
		}
		return nil, err
	default:
		return nil, fmt.Errorf("simplestreams: unknown signature policy %d", mirror.signaturePolicy)
	}
}

// IndexAt loads and returns the index document at path.
func (mirror *Mirror) IndexAt(ctx context.Context, path RelativePath) (*Index, error) {
	mirror.mu.Lock()
	if index := mirror.indexes[path]; index != nil {
		mirror.mu.Unlock()
		return index, nil
	}
	mirror.mu.Unlock()

	data, err := mirror.readMetadata(ctx, path)
	if err != nil {
		return nil, err
	}
	index, err := decodeIndex(data, mirror, path)
	if err != nil {
		return nil, err
	}

	mirror.mu.Lock()
	if cached := mirror.indexes[path]; cached != nil {
		mirror.mu.Unlock()
		return cached, nil
	}
	mirror.indexes[path] = index
	mirror.mu.Unlock()
	return index, nil
}

// ProductFile loads and returns the product document referenced by entry.
func (mirror *Mirror) ProductFile(ctx context.Context, entry *IndexEntry) (*ProductFile, error) {
	if entry == nil {
		return nil, errors.New("simplestreams: nil index entry")
	}
	if entry.Format == indexFormat {
		return nil, fmt.Errorf("simplestreams: index entry %q references nested index %q", entry.ContentID, entry.Path)
	}

	mirror.mu.Lock()
	if productFile := mirror.products[entry.Path]; productFile != nil {
		mirror.mu.Unlock()
		if err := mirror.checkContentID(entry, productFile); err != nil {
			return nil, err
		}
		return productFile, nil
	}
	mirror.mu.Unlock()

	data, err := mirror.readMetadata(ctx, entry.Path)
	if err != nil {
		return nil, err
	}
	productFile, err := decodeProductFile(data, mirror, entry.Path)
	if err != nil {
		return nil, err
	}
	if err := mirror.checkContentID(entry, productFile); err != nil {
		return nil, err
	}

	mirror.mu.Lock()
	if cached := mirror.products[entry.Path]; cached != nil {
		mirror.mu.Unlock()
		return cached, nil
	}
	mirror.products[entry.Path] = productFile
	mirror.mu.Unlock()
	return productFile, nil
}

func (mirror *Mirror) readMetadata(ctx context.Context, path RelativePath) ([]byte, error) {
	if err := path.Validate(); err != nil {
		return nil, err
	}
	if !path.IsMetadataPath() {
		return nil, fmt.Errorf("simplestreams: metadata path %q is not .json or .sjson", path)
	}
	if path.IsJSON() && mirror.signaturePolicy == SignatureRequired {
		return nil, fmt.Errorf("simplestreams: signature required for metadata path %q", path)
	}

	reader, err := mirror.source.Open(ctx, path)
	if err != nil {
		return nil, fmt.Errorf("simplestreams: open %q: %w", path, err)
	}
	defer func() { _ = reader.Close() }()

	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("simplestreams: read %q: %w", path, err)
	}
	if path.IsSignedJSON() {
		return mirror.decodeSigned(ctx, data)
	}
	if !path.IsJSON() {
		return nil, fmt.Errorf("simplestreams: unsupported metadata path %q", path)
	}
	return data, nil
}

func (mirror *Mirror) checkContentID(entry *IndexEntry, productFile *ProductFile) error {
	if mirror.allowContentIDMismatch || entry.ContentID == "" || productFile.ContentID == entry.ContentID {
		return nil
	}
	return fmt.Errorf(
		"simplestreams: product file %q content_id %q does not match index content ID %q",
		entry.Path,
		productFile.ContentID,
		entry.ContentID,
	)
}
