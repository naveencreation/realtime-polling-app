package auth

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"polling-backend/internal/models"
)

// Repository manages MongoDB persistence for user account documents.
type Repository struct {
	Users *mongo.Collection
}

// FindByEmailOrUsername looks up a user by either email or username to detect collisions.
func (repo Repository) FindByEmailOrUsername(ctx context.Context, email, username string) (*models.User, error) {
	var user models.User
	filter := bson.M{
		"$or": []bson.M{
			{"email": email},
			{"username": username},
		},
	}
	err := repo.Users.FindOne(ctx, filter).Decode(&user)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, nil
	}
	return &user, err
}

// FindByEmail retrieves a user by exact email match. Returns (nil, nil) if no user exists.
func (repo Repository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	err := repo.Users.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, nil
	}
	return &user, err
}

// Create inserts a new user record. Returns ErrDuplicateUser if email or username violates uniqueness.
func (repo Repository) Create(ctx context.Context, user *models.User) error {
	_, err := repo.Users.InsertOne(ctx, user)
	if mongo.IsDuplicateKeyError(err) {
		return ErrDuplicateUser
	}
	return err
}
