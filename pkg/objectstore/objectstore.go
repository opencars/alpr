package objectstore

import (
	"context"
	"io"
)

// ObjectStore is responsible for uploading objects.
type ObjectStore interface {
	Put(ctx context.Context, key string, r io.Reader) error
}
