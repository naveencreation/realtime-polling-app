package vote

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

// VotersKey returns the Redis set key tracking unique voter tokens for a poll.
func VotersKey(pollID string) string { return "poll:" + pollID + ":voters" }

// VotesKey returns the Redis hash key storing option vote tallies for a poll.
func VotesKey(pollID string) string { return "poll:" + pollID + ":votes" }

// UpdatesKey returns the Redis pub/sub channel used for real-time poll updates.
func UpdatesKey(pollID string) string { return "poll:" + pollID + ":updates" }

// NewVoterToken generates a unique UUID string identifying an anonymous voter.
func NewVoterToken() string {
	return uuid.NewString()
}

// registerVoteScript executes an atomic check-and-cast in Redis:
// 1. Checks if the voter token already exists in KEYS[1] (voters set).
// 2. If present, returns 0 without making changes.
// 3. Otherwise, adds the token to KEYS[1], increments the count in KEYS[2] (hash), and returns 1.
var registerVoteScript = redis.NewScript(`
local exists = redis.call('SISMEMBER', KEYS[1], ARGV[1])
if exists == 1 then
    return 0
end
redis.call('SADD', KEYS[1], ARGV[1])
redis.call('HINCRBY', KEYS[2], ARGV[2], 1)
return 1
`)

// Service manages live poll tallies, voter deduplication, and real-time broadcasting via Redis.
type Service struct {
	Redis *redis.Client
}

// HasVoted checks whether voterToken has already cast a ballot for pollID.
func (service Service) HasVoted(ctx context.Context, pollID, voterToken string) (bool, error) {
	return service.Redis.SIsMember(ctx, VotersKey(pollID), voterToken).Result()
}

// Counts returns the current vote distribution across all options for pollID.
func (service Service) Counts(ctx context.Context, pollID string) (map[string]int, error) {
	raw, err := service.Redis.HGetAll(ctx, VotesKey(pollID)).Result()
	if err != nil {
		return nil, err
	}
	out := make(map[string]int, len(raw))
	for optionID, valStr := range raw {
		if count, parseErr := strconv.Atoi(valStr); parseErr == nil {
			out[optionID] = count
		}
	}
	return out, nil
}

// RegisterVote atomically records a vote and publishes the updated tallies to subscribers.
// Returns (accepted: false, counts: nil, err: nil) if the voter has already voted.
func (service Service) RegisterVote(ctx context.Context, pollID, voterToken, optionID string) (bool, map[string]int, error) {
	keys := []string{VotersKey(pollID), VotesKey(pollID)}
	args := []any{voterToken, optionID}

	res, err := registerVoteScript.Run(ctx, service.Redis, keys, args...).Int64()
	if err != nil {
		return false, nil, err
	}
	if res == 0 {
		return false, nil, nil
	}

	counts, err := service.Counts(ctx, pollID)
	if err != nil {
		return false, nil, err
	}

	payload, err := json.Marshal(map[string]any{
		"event":    "results",
		"optionId": optionID,
		"counts":   counts,
	})
	if err != nil {
		return false, nil, err
	}

	if err = service.Redis.Publish(ctx, UpdatesKey(pollID), payload).Err(); err != nil {
		return false, nil, err
	}

	return true, counts, nil
}

// PublishClosed broadcasts a closure event to any active SSE subscribers for pollID.
func (service Service) PublishClosed(ctx context.Context, pollID string) error {
	return service.Redis.Publish(ctx, UpdatesKey(pollID), `{"event":"closed"}`).Err()
}

// DeletePollData purges Redis keys for pollID and broadcasts a deleted event to subscribers.
func (service Service) DeletePollData(ctx context.Context, pollID string) error {
	_ = service.Redis.Publish(ctx, UpdatesKey(pollID), `{"event":"deleted"}`)
	return service.Redis.Del(ctx, VotersKey(pollID), VotesKey(pollID)).Err()
}
