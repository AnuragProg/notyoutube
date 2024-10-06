package mq

import (
	"time"

	"github.com/anuragprog/notyoutube/file-service/types/database"
)

func FromRawVideoMetadataToProtoRawVideoMetadata(metadata database.RawVideoMetadata) *RawVideoMetadata{
	return &RawVideoMetadata{
		Id: metadata.Id,
		Filename: metadata.Filename,
		ContentType: metadata.ContentType,
		RequestId: metadata.RequestId,
		FileSize: metadata.FileSize,
		CreatedAt: metadata.CreatedAt.Format(time.RFC3339),
	}
}
