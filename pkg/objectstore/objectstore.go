package objectstore

import (
	"bytes"
	"context"
)

// ObjectStore is responsible for uploading objects.
type ObjectStore interface {
	Put(ctx context.Context, key string, r *bytes.Reader) error
}
