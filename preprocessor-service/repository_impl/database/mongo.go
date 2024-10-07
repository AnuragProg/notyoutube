package database

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

)

type MongoDatabase struct {
	client   *mongo.Client
	db       *mongo.Database
}

func NewMongoDatabase(uri, dbName string) (*MongoDatabase, error) {
	client, err := mongo.Connect(options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}
	db := client.Database(dbName)
	return &MongoDatabase{
		client:   client,
		db:       db,
	}, nil
}

func MustNewMongoDatabase(uri, dbName string) *MongoDatabase {
	client, err := mongo.Connect(options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}
	db := client.Database(dbName)
	return &MongoDatabase{
		client:   client,
		db:       db,
	}
}

func (md *MongoDatabase) Close() error {
	return md.client.Disconnect(context.TODO())
}
