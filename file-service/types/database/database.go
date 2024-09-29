package database

import (
	"time"
)

// id will be ignored in most of the scenarios
// e.g. id will be generated on insertion, and id will not be considered for updation
type RawVideoMetadata struct {
	Id          string `json:"_id"`
	Filename    string `json:"filename"`
	ContentType string `json:"content_type"`
	FileSize    int64  `json:"file_size"`
	// FileLocation string        `json:"-"` // object location // INFO: currently not implementing it as location is hardcoded in backend
	CreatedAt time.Time `json:"created_at"`
}

// fields should have same json tag as that in raw video metadata to maintain consistency
type UpdateRawVideoMetadata struct {
	Filename *string `json:"filename"`
}
