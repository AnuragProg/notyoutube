package store

import (
	"context"
	"io"
)

type Store interface {

	// to make sure the handler closes the connection properly
	io.Closer

	// Uploads a file or data to the object store.
	Upload(ctx context.Context, bucketName string, objectName string, reader io.Reader, objectSize int64, contentType string) error

	// Downloads an object from the object store.
	Download(ctx context.Context, bucketName string, objectName string) (io.ReadCloser, error)

	// Deletes an object from the object store.
	Delete(ctx context.Context, bucketName string, objectName string) error

	// Lists objects in a specific bucket or prefix.
	ListObjects(ctx context.Context, bucketName string, prefix string) ([]string, error)

	// Checks if an object exists.
	ObjectExists(ctx context.Context, bucketName string, objectName string) (bool, error)
}
