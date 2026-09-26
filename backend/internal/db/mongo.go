package db

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ConnectMongo initializes and validates a MongoDB client connection.
func ConnectMongo(ctx context.Context, mongoURI, databaseName string) (*mongo.Client, *mongo.Database, error) {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		return nil, nil, err
	}

	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err = client.Ping(pingCtx, nil); err != nil {
		_ = client.Disconnect(ctx)
		return nil, nil, err
	}
	return client, client.Database(databaseName), nil
}

type indexSpec struct {
	collection string
	keys       bson.D
	unique     bool
	name       string
}

// EnsureIndexes creates necessary uniqueness and lookup indexes in MongoDB.
func EnsureIndexes(ctx context.Context, database *mongo.Database) error {
	specs := []indexSpec{
		{
			collection: "users",
			keys:       bson.D{{Key: "username", Value: 1}},
			unique:     true,
			name:       "users_username_unique",
		},
		{
			collection: "users",
			keys:       bson.D{{Key: "email", Value: 1}},
			unique:     true,
			name:       "users_email_unique",
		},
		{
			collection: "polls",
			keys:       bson.D{{Key: "creatorId", Value: 1}, {Key: "createdAt", Value: -1}},
			unique:     false,
			name:       "polls_creator_created",
		},
	}

	for _, spec := range specs {
		model := mongo.IndexModel{
			Keys: spec.keys,
			Options: options.Index().
				SetUnique(spec.unique).
				SetName(spec.name),
		}
		if _, err := database.Collection(spec.collection).Indexes().CreateOne(ctx, model); err != nil {
			return err
		}
	}
	return nil
}
