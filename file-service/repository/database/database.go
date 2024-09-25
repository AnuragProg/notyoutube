package database

import (
	"context"
	"io"

	databaseType "github.com/anuragprog/notyoutube/file-service/types/database"
)

type Database interface {

	// to make sure the handler closes the connection properly
	io.Closer

    CreateRawVideoMetadata(ctx context.Context, metadata databaseType.RawVideoMetadata) (storedMetadata databaseType.RawVideoMetadata, err error)

    GetRawVideoMetadata(ctx context.Context, id string) (*databaseType.RawVideoMetadata, error)

    UpdateRawVideoMetadata(ctx context.Context, id string, metadata databaseType.UpdateRawVideoMetadata) error

    DeleteRawVideoMetadata(ctx context.Context, id string) error

    ListRawVideosMetadata(ctx context.Context, limit int, offset int) ([]databaseType.RawVideoMetadata, error)
}
