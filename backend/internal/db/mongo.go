package db

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

func ConnectMongo(ctx context.Context, uri, name string) (*mongo.Client, *mongo.Database, error) {
	c, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, nil, err
	}
	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err = c.Ping(pingCtx, nil); err != nil {
		_ = c.Disconnect(ctx)
		return nil, nil, err
	}
	return c, c.Database(name), nil
}
func EnsureIndexes(ctx context.Context, database *mongo.Database) error {
	for _, spec := range []struct {
		collection string
		keys       bson.D
		unique     bool
		name       string
	}{{"users", bson.D{{Key: "username", Value: 1}}, true, "users_username_unique"}, {"users", bson.D{{Key: "email", Value: 1}}, true, "users_email_unique"}, {"polls", bson.D{{Key: "creatorId", Value: 1}, {Key: "createdAt", Value: -1}}, false, "polls_creator_created"}} {
		_, err := database.Collection(spec.collection).Indexes().CreateOne(ctx, mongo.IndexModel{Keys: spec.keys, Options: options.Index().SetUnique(spec.unique).SetName(spec.name)})
		if err != nil {
			return err
		}
	}
	return nil
}
