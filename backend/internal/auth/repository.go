package auth

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"polling-backend/internal/models"
)

type Repository struct{ Users *mongo.Collection }

func (r Repository) FindByEmailOrUsername(ctx context.Context, email, username string) (*models.User, error) {
	var u models.User
	err := r.Users.FindOne(ctx, bson.M{"$or": []bson.M{{"email": email}, {"username": username}}}).Decode(&u)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, nil
	}
	return &u, err
}
func (r Repository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	var u models.User
	err := r.Users.FindOne(ctx, bson.M{"email": email}).Decode(&u)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, nil
	}
	return &u, err
}
func (r Repository) Create(ctx context.Context, u *models.User) error {
	_, err := r.Users.InsertOne(ctx, u)
	if mongo.IsDuplicateKeyError(err) {
		return ErrDuplicateUser
	}
	return err
}
