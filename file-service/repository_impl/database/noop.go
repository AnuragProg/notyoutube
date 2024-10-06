package database

import (
	"context"
	databaseType "github.com/anuragprog/notyoutube/file-service/types/database"
)

type NoopDatabase struct {}

func NewNoopDatabse() *NoopDatabase {
	return &NoopDatabase{}
}

func (nd *NoopDatabase) CreateRawVideoMetadata(context.Context, databaseType.RawVideoMetadata) (databaseType.RawVideoMetadata, error) { return databaseType.RawVideoMetadata{}, nil }

func (nd *NoopDatabase) GetRawVideoMetadata(context.Context, string) (*databaseType.RawVideoMetadata, error) { return &databaseType.RawVideoMetadata{}, nil }

func (nd *NoopDatabase) UpdateRawVideoMetadata(context.Context, string, databaseType.UpdateRawVideoMetadata) error { return nil }

func (nd *NoopDatabase) DeleteRawVideoMetadata(context.Context, string) error { return nil }

func (nd *NoopDatabase) ListRawVideosMetadata(context.Context, int, int) ([]databaseType.RawVideoMetadata, error) { return []databaseType.RawVideoMetadata{}, nil }

func (nd *NoopDatabase) Close() error { return nil }
