package database

import (
	"reflect"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type MongoRawVideoMetadata struct {
	Id           bson.ObjectID `bson:"_id"`
	Filename     string        `bson:"filename"`
	ContentType  string        `bson:"content_type"`
	FileSize     int64         `bson:"file_size"`
	// FileLocation string        `bson:"-"` // object location // INFO: currently not implementing it as location is hardcoded in backend
	CreatedAt    time.Time     `bson:"created_at"`
}

func FromRawVideoMetadataToMongoRawVideoMetadataIgnoringId(metadata RawVideoMetadata) MongoRawVideoMetadata {
	id := bson.NewObjectID()
	return MongoRawVideoMetadata{
		Id:           id,
		Filename:     metadata.Filename,
		ContentType:  metadata.ContentType,
		FileSize:     metadata.FileSize,
		// FileLocation: metadata.FileLocation,
		CreatedAt:    metadata.CreatedAt,
	}
}

func (metadata *MongoRawVideoMetadata) ToRawVideoMetadata() RawVideoMetadata {
	return RawVideoMetadata{
		Id:           metadata.Id.Hex(),
		Filename:     metadata.Filename,
		ContentType:  metadata.ContentType,
		FileSize:     metadata.FileSize,
		// FileLocation: metadata.FileLocation,
		CreatedAt:    metadata.CreatedAt,
	}
}

// we are retrieving all the fields with json tag because the type is fixed and hence all
// fields are valid
func RetrieveFieldsToBeUpdated(metadata UpdateRawVideoMetadata) map[string]interface{} {
	fields := map[string]interface{}{}

	metadataVal := reflect.ValueOf(metadata)
	metadataTyp := reflect.TypeOf(metadata)

	for idx := range metadataVal.NumField() {
		field := metadataVal.Field(idx)
		fieldTyp := metadataTyp.Field(idx)

		if field.Kind() == reflect.Ptr && !field.IsNil() {
			jsonTag := fieldTyp.Tag.Get("json")
			if jsonTag == "" {
				jsonTag = fieldTyp.Name
			}
			fields[jsonTag] = field.Elem().Interface()
		}
	}

	return fields
}
