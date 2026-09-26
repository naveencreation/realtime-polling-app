package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// PollStatus represents the lifecycle state of a poll.
type PollStatus string

const (
	StatusOpen   PollStatus = "open"
	StatusClosed PollStatus = "closed"
)

// User represents an application user with credentials.
type User struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Username     string             `bson:"username" json:"username"`
	Email        string             `bson:"email" json:"-"`
	PasswordHash string             `bson:"passwordHash" json:"-"`
	CreatedAt    time.Time          `bson:"createdAt" json:"createdAt"`
}

// Option represents a selectable choice within a poll.
type Option struct {
	ID   string `bson:"id" json:"id"`
	Text string `bson:"text" json:"text"`
}

// Poll represents a question and its options created by a user.
type Poll struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	CreatorID primitive.ObjectID `bson:"creatorId" json:"-"`
	Question  string             `bson:"question" json:"question"`
	Options   []Option           `bson:"options" json:"options"`
	Status    PollStatus         `bson:"status" json:"status"`
	ExpiresAt *time.Time         `bson:"expiresAt,omitempty" json:"expiresAt,omitempty"`
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt"`
}

// Vote represents an audit log entry for a cast vote.
type Vote struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	PollID     primitive.ObjectID `bson:"pollId"`
	OptionID   string             `bson:"optionId"`
	VoterToken string             `bson:"voterToken"`
	VotedAt    time.Time          `bson:"votedAt"`
}
