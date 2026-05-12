package simplestreams

import (
	"errors"
	"fmt"
	"strings"
)

const (
	// DefaultIndexPath is the conventional unsigned Simple Streams index path.
	DefaultIndexPath RelativePath = "streams/v1/index.json"

	// DefaultSignedIndexPath is the conventional cleartext signed Simple Streams index path.
	DefaultSignedIndexPath RelativePath = "streams/v1/index.sjson"
)

// RelativePath is a backend-neutral Simple Streams path relative to a mirror root.
type RelativePath string

// ParseRelativePath validates path as a Simple Streams mirror-relative path.
func ParseRelativePath(path string) (RelativePath, error) {
	relativePath := RelativePath(path)
	if err := relativePath.Validate(); err != nil {
		return "", err
	}
	return relativePath, nil
}

// String returns the path as it appears in Simple Streams metadata.
func (p RelativePath) String() string {
	return string(p)
}

// Validate checks that p is non-empty, relative, and cannot traverse upward.
func (p RelativePath) Validate() error {
	value := string(p)
	if value == "" {
		return errors.New("simplestreams: empty relative path")
	}
	if strings.HasPrefix(value, "/") {
		return fmt.Errorf("simplestreams: relative path %q is absolute", value)
	}
	if strings.Contains(value, "\\") {
		return fmt.Errorf("simplestreams: relative path %q contains backslash", value)
	}

	for segment := range strings.SplitSeq(value, "/") {
		switch segment {
		case "", ".", "..", "...":
			return fmt.Errorf("simplestreams: relative path %q contains unsafe segment %q", value, segment)
		}
	}
	return nil
}

// IsJSON reports whether p names an unsigned JSON metadata document.
func (p RelativePath) IsJSON() bool {
	return strings.HasSuffix(string(p), ".json")
}

// IsSignedJSON reports whether p names a cleartext signed JSON metadata document.
func (p RelativePath) IsSignedJSON() bool {
	return strings.HasSuffix(string(p), ".sjson")
}

// IsMetadataPath reports whether p names a supported Simple Streams metadata document.
func (p RelativePath) IsMetadataPath() bool {
	return p.IsJSON() || p.IsSignedJSON()
}
