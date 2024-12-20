package database

import (
)

type NoopDatabase struct {}

func NewNoopDatabase() *NoopDatabase {
	return &NoopDatabase{}
}

func (nd *NoopDatabase) Close() error { return nil }
