package poll

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"polling-backend/internal/models"
)

// Repository encapsulates MongoDB operations for poll records.
type Repository struct {
	Collection *mongo.Collection
}

// Create persists a new poll document with an initial open status and UTC timestamp.
func (repo Repository) Create(ctx context.Context, pollItem *models.Poll) error {
	pollItem.Status = models.StatusOpen
	pollItem.CreatedAt = time.Now().UTC()
	_, err := repo.Collection.InsertOne(ctx, pollItem)
	return err
}

// FindByID retrieves a single poll by its ObjectID. Returns (nil, nil) if not found.
func (repo Repository) FindByID(ctx context.Context, pollID primitive.ObjectID) (*models.Poll, error) {
	var pollItem models.Poll
	err := repo.Collection.FindOne(ctx, bson.M{"_id": pollID}).Decode(&pollItem)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, nil
	}
	return &pollItem, err
}

// FindByCreator lists all polls belonging to a creator, sorted newest first.
func (repo Repository) FindByCreator(ctx context.Context, creatorID primitive.ObjectID) ([]models.Poll, error) {
	cursor, err := repo.Collection.Find(
		ctx,
		bson.M{"creatorId": creatorID},
		options.Find().SetSort(bson.D{{Key: "createdAt", Value: -1}}),
	)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var polls []models.Poll
	err = cursor.All(ctx, &polls)
	return polls, err
}

// SetStatus updates a poll's status (e.g. open to closed).
func (repo Repository) SetStatus(ctx context.Context, pollID primitive.ObjectID, status models.PollStatus) error {
	_, err := repo.Collection.UpdateOne(ctx, bson.M{"_id": pollID}, bson.M{"$set": bson.M{"status": status}})
	return err
}

// ResolveLazyExpiry checks if an open poll has passed its expiration deadline.
// If expired, it closes the poll in the database and updates the in-memory struct.
func (repo Repository) ResolveLazyExpiry(ctx context.Context, pollItem *models.Poll) error {
	if pollItem.Status == models.StatusOpen && pollItem.ExpiresAt != nil && time.Now().After(*pollItem.ExpiresAt) {
		if err := repo.SetStatus(ctx, pollItem.ID, models.StatusClosed); err != nil {
			return err
		}
		pollItem.Status = models.StatusClosed
	}
	return nil
}
