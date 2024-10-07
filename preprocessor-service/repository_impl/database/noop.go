package database

import (
)

type NoopDatabase struct {}

func NewNoopDatabse() *NoopDatabase {
	return &NoopDatabase{}
}

func (nd *NoopDatabase) Close() error { return nil }
