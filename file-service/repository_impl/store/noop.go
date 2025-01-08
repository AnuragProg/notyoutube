package store

import (
	"context"
	"io"
	"strings"

	storeTypes "github.com/anuragprog/notyoutube/file-service/types/store"
)

type NoopStore struct {}

func NewNoopStore() *NoopStore { return &NoopStore{} }

func (*NoopStore) Upload(context.Context, string, string, io.Reader, int64, string) error { return nil }

func (*NoopStore) Download(context.Context, string, string) (io.ReadCloser, error) { return io.NopCloser(strings.NewReader("")), nil }

func (*NoopStore) Delete(context.Context, string, string) error { return nil }

func (*NoopStore) ListObjects(context.Context, string, string) ([]string, error) { return []string{}, nil }

func (*NoopStore) ObjectExists(ctx context.Context, bucketName string, objectName string) (bool, error) { return false, nil }

func (*NoopStore) GetPresignedUrl(ctx context.Context, bucketName string, objectName string) (storeTypes.PresignUrlResult, error) { return storeTypes.PresignUrlResult{}, nil }

func (*NoopStore) Close() error { return nil }
