// Package fsmirror provides a filesystem-backed Simple Streams source.
package fsmirror

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	simplestreams "github.com/meigma/go-simplestreams"
)

// Source opens Simple Streams content from a filesystem root.
type Source struct {
	root string
}

// New constructs a filesystem source rooted at root.
func New(root string) (*Source, error) {
	if root == "" {
		return nil, errors.New("fsmirror: empty root")
	}
	absoluteRoot, err := filepath.Abs(root)
	if err != nil {
		return nil, err
	}
	return &Source{root: absoluteRoot}, nil
}

// Open opens path under the source root.
func (source *Source) Open(ctx context.Context, path simplestreams.RelativePath) (io.ReadCloser, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if err := path.Validate(); err != nil {
		return nil, err
	}

	fullPath := filepath.Join(source.root, filepath.FromSlash(path.String()))
	cleanRoot := filepath.Clean(source.root)
	cleanFullPath := filepath.Clean(fullPath)
	if cleanFullPath != cleanRoot && !strings.HasPrefix(cleanFullPath, cleanRoot+string(os.PathSeparator)) {
		return nil, fmt.Errorf("fsmirror: path %q escapes root", path)
	}
	file, err := os.Open(cleanFullPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, fmt.Errorf("%w: %s", simplestreams.ErrNotFound, path)
		}
		return nil, err
	}
	return file, nil
}
