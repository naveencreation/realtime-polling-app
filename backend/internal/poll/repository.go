package poll

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"polling-backend/internal/models"
	"time"
)

type Repository struct{ Collection *mongo.Collection }

func (r Repository) Create(ctx context.Context, p *models.Poll) error {
	p.Status = models.StatusOpen
	p.CreatedAt = time.Now().UTC()
	_, err := r.Collection.InsertOne(ctx, p)
	return err
}
func (r Repository) FindByID(ctx context.Context, id primitive.ObjectID) (*models.Poll, error) {
	var p models.Poll
	err := r.Collection.FindOne(ctx, bson.M{"_id": id}).Decode(&p)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	return &p, err
}
func (r Repository) FindByCreator(ctx context.Context, id primitive.ObjectID) ([]models.Poll, error) {
	cur, err := r.Collection.Find(ctx, bson.M{"creatorId": id}, options.Find().SetSort(bson.D{{Key: "createdAt", Value: -1}}))
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	var out []models.Poll
	err = cur.All(ctx, &out)
	return out, err
}
func (r Repository) SetStatus(ctx context.Context, id primitive.ObjectID, status models.PollStatus) error {
	_, err := r.Collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{"status": status}})
	return err
}
func (r Repository) ResolveLazyExpiry(ctx context.Context, p *models.Poll) error {
	if p.Status == models.StatusOpen && p.ExpiresAt != nil && time.Now().After(*p.ExpiresAt) {
		if err := r.SetStatus(ctx, p.ID, models.StatusClosed); err != nil {
			return err
		}
		p.Status = models.StatusClosed
	}
	return nil
}
