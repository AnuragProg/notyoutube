/*
Store Manager has been created purely for deciding where to store a certain resource in the store
Except for deciding the location of the resource StoreManager has no other purpose and should not be coded that way
*/

package store

import (
	"context"
	"errors"
	"fmt"
	"io"
)

type Resource string

const (
	RAW_VIDEO Resource = "raw_video"
)

var resourceLocations = map[Resource]string{
	RAW_VIDEO: "raw_videos",
}

// Store manager book keeps the resource's location
type StoreManager struct {
	bucket string
	store  Store
}

func NewStoreManager(bucket string, store Store) *StoreManager {
	return &StoreManager{
		bucket: bucket,
		store:  store,
	}
}

func (sm *StoreManager) getObjectResourceLocation(resource Resource, objectName string) (string, error) {
	resourceLocation, ok := resourceLocations[resource]
	if !ok {
		return "", errors.New("invalid resource")
	}

	return fmt.Sprintf("%v/%v", resourceLocation, objectName), nil
}

func (sm *StoreManager) Upload(ctx context.Context, resource Resource, objectName string, reader io.Reader, objectSize int64, contentType string) error {
	objectPath, err := sm.getObjectResourceLocation(resource, objectName)
	if err != nil {
		return err
	}
	return sm.store.Upload(ctx, sm.bucket, objectPath, reader, objectSize, contentType)
}

func (sm *StoreManager) Download(ctx context.Context, resource Resource, objectName string) (io.ReadCloser, error) {
	objectPath, err := sm.getObjectResourceLocation(resource, objectName)
	if err != nil {
		return nil, err
	}
	return sm.store.Download(ctx, sm.bucket, objectPath)
}
