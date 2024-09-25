package database

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	databaseType "github.com/anuragprog/notyoutube/file-service/types/database"
)

type MongoDatabase struct {
	client   *mongo.Client
	db       *mongo.Database
	rawVideoCol *mongo.Collection
}

func NewMongoDatabase(uri string) (*MongoDatabase, error) {
	client, err := mongo.Connect(options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}
	db := client.Database("not-youtube")
	return &MongoDatabase{
		client:   client,
		db:       db,
		rawVideoCol: db.Collection("raw-videos"),
	}, nil
}

func (md *MongoDatabase) Close() error {
	return md.client.Disconnect(context.TODO())
}

func (md *MongoDatabase) CreateRawVideoMetadata(ctx context.Context, metadata databaseType.RawVideoMetadata) (databaseType.RawVideoMetadata, error) {
	mongoRawVideoMetadata := databaseType.FromRawVideoMetadataToMongoRawVideoMetadataIgnoringId(metadata)
	_, err := md.rawVideoCol.InsertOne(ctx, mongoRawVideoMetadata)
	if err != nil {
		return databaseType.RawVideoMetadata{}, err
	}
	return mongoRawVideoMetadata.ToRawVideoMetadata(), nil
}

func (md *MongoDatabase) GetRawVideoMetadata(ctx context.Context, id string) (*databaseType.RawVideoMetadata, error) {
	objectId, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid video id")
	}
	filter := bson.M{"_id": objectId}

	var mongoRawVideoMetadata databaseType.MongoRawVideoMetadata
	err = md.rawVideoCol.FindOne(ctx, filter).Decode(&mongoRawVideoMetadata)
	if err != nil {
		switch {
		case errors.Is(err, mongo.ErrNoDocuments):
			return nil, errors.New("video metadata not found")
		default:
			return nil, err
		}
	}

	rawVideoMetadata := mongoRawVideoMetadata.ToRawVideoMetadata()
	return &rawVideoMetadata, nil
}

func (md *MongoDatabase) UpdateRawVideoMetadata(ctx context.Context, id string, metadata databaseType.UpdateRawVideoMetadata) error {
	objectId, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid video id")
	}

	update := bson.M{"$set": databaseType.RetrieveFieldsToBeUpdated(metadata)}
	result, err := md.rawVideoCol.UpdateByID(ctx, objectId, update)
	if err != nil {
		return err
	}
	if result.ModifiedCount == 0 {
		return errors.New("no video was updated")
	}
	return nil
}

func (md *MongoDatabase) DeleteRawVideoMetadata(ctx context.Context, id string) error {
	objectId, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid video id")
	}
	filter := bson.M{"_id": objectId}

	result, err := md.rawVideoCol.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return errors.New("no video was deleted")
	}
	return nil
}

func (md *MongoDatabase) ListRawVideosMetadata(ctx context.Context, limit int, offset int) ([]databaseType.RawVideoMetadata, error) {

	findOption := options.Find().
		SetLimit(int64(limit)).
		SetSkip(int64(offset))

	cursor, err := md.rawVideoCol.Find(ctx, bson.M{}, findOption)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	mongoVideoMetadatas := []databaseType.MongoRawVideoMetadata{}
	if err := cursor.All(ctx, &mongoVideoMetadatas); err != nil { return nil, err }

	videoMetadatas := make([]databaseType.RawVideoMetadata, 0, len(mongoVideoMetadatas))
	for _, metadata := range mongoVideoMetadatas {
		videoMetadatas = append(videoMetadatas, metadata.ToRawVideoMetadata())
	}

	return videoMetadatas, nil
}
