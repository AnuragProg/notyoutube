package database

import (
	"strings"
	"time"
)

// id will be ignored in most of the scenarios
// e.g. id will be generated on insertion, and id will not be considered for updation
type RawVideoMetadata struct {
	Id          string `json:"_id"`
	Filename    string `json:"filename"`
	ContentType string `json:"content_type"`
	FileSize    int64  `json:"file_size"`
	RequestId   string `json:"request_id"`
	// FileLocation string        `json:"-"` // object location // INFO: currently not implementing it is hardcoded in backend
	CreatedAt time.Time `json:"created_at"`
}

// thanks to fiber's no reuse of string policy
func (metadata *RawVideoMetadata) DeepCopy() RawVideoMetadata {
	return RawVideoMetadata{
		Id: strings.Clone(metadata.Id),
		Filename: strings.Clone(metadata.Filename),
		ContentType: strings.Clone(metadata.ContentType),
		FileSize: metadata.FileSize,
		RequestId: strings.Clone(metadata.RequestId),
		CreatedAt: metadata.CreatedAt,
	}
}


// fields should have same json tag as that in raw video metadata to maintain consistency
type UpdateRawVideoMetadata struct {
	Filename *string `json:"filename"`
}
