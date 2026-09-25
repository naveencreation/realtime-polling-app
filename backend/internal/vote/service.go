package vote

import (
	"context"
	"encoding/json"
	"github.com/redis/go-redis/v9"
	"strconv"
)

type Service struct{ Redis *redis.Client }

func key(p, s string) string { return "poll:" + p + ":" + s }
func (s Service) HasVoted(ctx context.Context, p, t string) (bool, error) {
	return s.Redis.SIsMember(ctx, key(p, "voters"), t).Result()
}
func (s Service) Counts(ctx context.Context, p string) (map[string]int, error) {
	raw, err := s.Redis.HGetAll(ctx, key(p, "votes")).Result()
	out := map[string]int{}
	for k, v := range raw {
		n, e := strconv.Atoi(v)
		if e == nil {
			out[k] = n
		}
	}
	return out, err
}
func (s Service) RegisterVote(ctx context.Context, p, t, o string) (bool, map[string]int, error) {
	n, err := s.Redis.SAdd(ctx, key(p, "voters"), t).Result()
	if err != nil {
		return false, nil, err
	}
	if n == 0 {
		return false, nil, nil
	}
	if err = s.Redis.HIncrBy(ctx, key(p, "votes"), o, 1).Err(); err != nil {
		return false, nil, err
	}
	counts, err := s.Counts(ctx, p)
	if err != nil {
		return false, nil, err
	}
	b, _ := json.Marshal(map[string]any{"event": "results", "optionId": o, "counts": counts})
	if err = s.Redis.Publish(ctx, key(p, "updates"), b).Err(); err != nil {
		return false, nil, err
	}
	return true, counts, nil
}
func (s Service) PublishClosed(ctx context.Context, p string) error {
	return s.Redis.Publish(ctx, key(p, "updates"), `{"event":"closed"}`).Err()
}
