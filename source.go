package simplestreams

import (
	"context"
	"errors"
	"io"
)

// ErrNotFound reports that a requested mirror-relative path does not exist.
var ErrNotFound = errors.New("simplestreams: not found")

// Source opens Simple Streams content by mirror-relative path.
type Source interface {
	Open(ctx context.Context, path RelativePath) (io.ReadCloser, error)
}

// SourceFunc adapts a function into a Source.
type SourceFunc func(context.Context, RelativePath) (io.ReadCloser, error)

// Open calls f(ctx, path).
func (f SourceFunc) Open(ctx context.Context, path RelativePath) (io.ReadCloser, error) {
	return f(ctx, path)
}

// Store writes and removes Simple Streams content by mirror-relative path.
type Store interface {
	Put(ctx context.Context, path RelativePath, content io.Reader) error
	Delete(ctx context.Context, path RelativePath) error
}

// AtomicStore is implemented by stores that stage writes before publishing them.
type AtomicStore interface {
	Store
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
}
