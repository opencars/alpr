package minio

import (
	"context"
	"io"

	"github.com/minio/minio-go"

	"github.com/opencars/alpr/pkg/config"
)

// ObjectStore is respnosible for uploading objects.
type ObjectStore struct {
	client *minio.Client
	bucket string
}

// New returns either newly allocated store, or error.
func New(s3 *config.S3) (*ObjectStore, error) {
	// Initiate a client using DigitalOcean Spaces.
	client, err := minio.New(s3.Endpoint, s3.AccessKeyID, s3.SecretAccessKey, s3.SSL)
	if err != nil {
		return nil, err
	}

	return &ObjectStore{
		client: client,
		bucket: s3.Bucket,
	}, nil
}

// Put uploads new image to S3 bucket.
func (s *ObjectStore) Put(ctx context.Context, key string, r io.Reader) error {
	userMetaData := map[string]string{"x-amz-acl": "public-read"}
	opts := minio.PutObjectOptions{
		UserMetadata: userMetaData,
	}

	_, err := s.client.PutObjectWithContext(ctx, s.bucket, "plates/"+key, r, -1, opts)
	if err != nil {
		return err
	}

	return nil
}
