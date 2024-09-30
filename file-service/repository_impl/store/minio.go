package store

import (
	"context"
	"io"
	"net/http"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinioStore struct {
	client *minio.Client
}

func NewMinioStore(minioURI, minioServerAccessKey, minioServerSecretKey string) (*MinioStore, error) {
	// create client
	minioClient, err := minio.New(minioURI, &minio.Options{
		Creds: credentials.NewStaticV4(minioServerAccessKey, minioServerSecretKey, ""),
		Transport: &http.Transport{
			MaxIdleConns:    100,
			IdleConnTimeout: 60 * time.Second,
		},
	})
	if err != nil {
		return nil, err
	}

	return &MinioStore{
		client: minioClient,
	}, nil
}

func (ms *MinioStore) Upload(ctx context.Context, bucketName string, objectName string, reader io.Reader, objectSize int64, contentType string) error {
	_, err := ms.client.PutObject(ctx, bucketName, objectName, reader, objectSize, minio.PutObjectOptions{
		ContentType: contentType,
	})
	return err
}

func (ms *MinioStore) Download(ctx context.Context, bucketName string, objectName string) (io.ReadCloser, error) {
	object, err := ms.client.GetObject(ctx, bucketName, objectName, minio.GetObjectOptions{})
	return object, err
}

func (ms *MinioStore) Delete(ctx context.Context, bucketName string, objectName string) error {
	return ms.client.RemoveObject(ctx, bucketName, objectName, minio.RemoveObjectOptions{})
}

func (ms *MinioStore) ListObjects(ctx context.Context, bucketName string, prefix string) ([]string, error) {
	objectChan := ms.client.ListObjects(ctx, bucketName, minio.ListObjectsOptions{
		Prefix: prefix,
	})

	objectNames := []string{}
	for object := range objectChan {
		objectNames = append(objectNames, object.Key)
	}
	return objectNames, nil
}

func (ms *MinioStore) ObjectExists(ctx context.Context, bucketName string, objectName string) (bool, error) {
	_, err := ms.client.StatObject(ctx, bucketName, objectName, minio.GetObjectOptions{})
	if err != nil {
		if minio.ToErrorResponse(err).Code == "NoSuchKey" {
			return false, nil
		}
		return false, err
	}
	return true, nil
}


func (ms *MinioStore) Close() error {
	return nil
}
