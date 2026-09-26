package vote

import (
	"testing"

	"github.com/google/uuid"
)

func TestKeyHelpers(t *testing.T) {
	pollID := "651a2b3c4d5e6f7a8b9c0d1e"

	if got := VotersKey(pollID); got != "poll:651a2b3c4d5e6f7a8b9c0d1e:voters" {
		t.Errorf("unexpected VotersKey: %s", got)
	}
	if got := VotesKey(pollID); got != "poll:651a2b3c4d5e6f7a8b9c0d1e:votes" {
		t.Errorf("unexpected VotesKey: %s", got)
	}
	if got := UpdatesKey(pollID); got != "poll:651a2b3c4d5e6f7a8b9c0d1e:updates" {
		t.Errorf("unexpected UpdatesKey: %s", got)
	}
}

func TestNewVoterToken(t *testing.T) {
	token1 := NewVoterToken()
	token2 := NewVoterToken()

	if token1 == "" || token2 == "" {
		t.Fatal("generated voter token should not be empty")
	}
	if token1 == token2 {
		t.Fatal("subsequent voter tokens should be unique")
	}

	if _, err := uuid.Parse(token1); err != nil {
		t.Fatalf("voter token is not a valid UUID: %v", err)
	}
}
