package simplestreams

import (
	"crypto/md5" // #nosec G501 -- Simple Streams metadata can publish MD5 checksums for compatibility.
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"errors"
	"fmt"
	"hash"
	"io"
	"strings"
)

// ArtifactRef describes downloadable artifact content referenced by an item.
type ArtifactRef struct {
	// Path is the primary artifact path.
	Path RelativePath

	// Mirrors are alternate mirror prefixes for Path.
	Mirrors []MirrorPrefix

	// Size is the declared content size when available.
	Size *int64

	// Checksums contains supported checksum values by algorithm.
	Checksums map[string]string
}

// Paths returns the primary path followed by mirror-expanded alternate paths.
func (artifact ArtifactRef) Paths() []string {
	paths := []string{artifact.Path.String()}
	for _, mirror := range artifact.Mirrors {
		prefix := strings.TrimRight(string(mirror), "/")
		if prefix == "" {
			continue
		}
		paths = append(paths, prefix+"/"+artifact.Path.String())
	}
	return paths
}

// VerifyReader verifies r against artifact size and strongest available checksum.
func (artifact ArtifactRef) VerifyReader(r io.Reader) error {
	return VerifyReader(r, artifact.Checksums, artifact.Size)
}

// Checksum is one checksum assertion for artifact bytes.
type Checksum struct {
	// Algorithm is the checksum algorithm name.
	Algorithm string

	// Value is the expected lowercase hex digest.
	Value string
}

// Checksums returns supported checksum fields from item.
func Checksums(item *Item) map[string]string {
	if item == nil {
		return map[string]string{}
	}
	checksums := map[string]string{}
	if item.MD5 != "" {
		checksums["md5"] = item.MD5
	}
	if item.SHA256 != "" {
		checksums["sha256"] = item.SHA256
	}
	if item.SHA512 != "" {
		checksums["sha512"] = item.SHA512
	}
	return checksums
}

// BestChecksum returns the strongest supported checksum in checksums.
func BestChecksum(checksums map[string]string) (Checksum, bool) {
	for _, algorithm := range []string{"sha512", "sha256", "md5"} {
		if value := checksums[algorithm]; value != "" {
			return Checksum{Algorithm: algorithm, Value: strings.ToLower(value)}, true
		}
	}
	return Checksum{}, false
}

// VerifyReader verifies r against size and the strongest supported checksum.
func VerifyReader(r io.Reader, checksums map[string]string, size *int64) error {
	if r == nil {
		return errors.New("simplestreams: nil artifact reader")
	}

	checksum, hasChecksum := BestChecksum(checksums)
	var hasher hash.Hash
	if hasChecksum {
		var err error
		hasher, err = newHasher(checksum.Algorithm)
		if err != nil {
			return err
		}
	}

	var written int64
	var err error
	if hasChecksum {
		written, err = io.Copy(hasher, r)
	} else {
		written, err = io.Copy(io.Discard, r)
	}
	if err != nil {
		return fmt.Errorf("simplestreams: read artifact: %w", err)
	}
	if size != nil && written != *size {
		return fmt.Errorf("simplestreams: artifact size mismatch: got %d, want %d", written, *size)
	}
	if !hasChecksum {
		return nil
	}

	actual := hex.EncodeToString(hasher.Sum(nil))
	if actual != checksum.Value {
		return fmt.Errorf(
			"simplestreams: artifact %s mismatch: got %s, want %s",
			checksum.Algorithm,
			actual,
			checksum.Value,
		)
	}
	return nil
}

func newHasher(algorithm string) (hash.Hash, error) {
	switch algorithm {
	case "md5":
		return md5.New(), nil // #nosec G401 -- MD5 is used only for protocol checksum verification.
	case "sha256":
		return sha256.New(), nil
	case "sha512":
		return sha512.New(), nil
	default:
		return nil, fmt.Errorf("simplestreams: unsupported checksum algorithm %q", algorithm)
	}
}
