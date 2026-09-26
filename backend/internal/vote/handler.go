package vote

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"polling-backend/internal/middleware"
	"polling-backend/internal/models"
	"polling-backend/internal/poll"
	"polling-backend/pkg/response"
)

// Handler processes cast ballots and updates poll vote tallies.
type Handler struct {
	Polls       poll.Repository
	Service     Service
	Audit       *mongo.Collection
	VoterCookie string
	Secure      bool
}

type voteInput struct {
	OptionID string `json:"optionId"`
}

// Vote handles POST /api/polls/:id/vote.
// It verifies poll openness, assigns or reads anonymous voter tokens, registers the vote
// atomically in Redis, and asynchronously logs the audit trail in MongoDB.
func (handler Handler) Vote(ginCtx *gin.Context) {
	logger := middleware.GetLogger(ginCtx)

	pollID, err := primitive.ObjectIDFromHex(ginCtx.Param("id"))
	if err != nil {
		response.Error(ginCtx, http.StatusNotFound, "not_found", "poll not found")
		return
	}

	pollItem, err := handler.Polls.FindByID(ginCtx, pollID)
	if err != nil {
		logger.Error("database error retrieving poll for vote", "poll_id", pollID.Hex(), "error", err)
		response.Error(ginCtx, http.StatusInternalServerError, "server_error", "something went wrong")
		return
	}
	if pollItem == nil {
		response.Error(ginCtx, http.StatusNotFound, "not_found", "poll not found")
		return
	}

	if err = handler.Polls.ResolveLazyExpiry(ginCtx, pollItem); err != nil {
		logger.Error("database error resolving poll expiry in vote", "poll_id", pollID.Hex(), "error", err)
		response.Error(ginCtx, http.StatusInternalServerError, "server_error", "something went wrong")
		return
	}
	if pollItem.Status != models.StatusOpen {
		response.Error(ginCtx, http.StatusConflict, "poll_closed", "poll is closed")
		return
	}

	var req voteInput
	if err := ginCtx.ShouldBindJSON(&req); err != nil || strings.TrimSpace(req.OptionID) == "" {
		response.Error(ginCtx, http.StatusBadRequest, "validation_error", "optionId is required")
		return
	}

	isValidOption := false
	for _, option := range pollItem.Options {
		if option.ID == req.OptionID {
			isValidOption = true
			break
		}
	}
	if !isValidOption {
		response.Error(ginCtx, http.StatusBadRequest, "invalid_option", "option is not valid for this poll")
		return
	}

	// Resolve voter token from cookie or request header; issue a new one if missing
	voterToken, err := ginCtx.Cookie(handler.VoterCookie)
	if err != nil || voterToken == "" {
		voterToken = ginCtx.GetHeader("X-Voter-Token")
	}
	if voterToken == "" {
		voterToken = NewVoterToken()
		sameSite := http.SameSiteLaxMode
		if handler.Secure {
			sameSite = http.SameSiteNoneMode
		}
		http.SetCookie(ginCtx.Writer, &http.Cookie{
			Name:     handler.VoterCookie,
			Value:    voterToken,
			MaxAge:   31536000,
			HttpOnly: true,
			Secure:   handler.Secure,
			SameSite: sameSite,
			Path:     "/",
		})
	}
	ginCtx.Header("X-Voter-Token", voterToken)

	accepted, counts, err := handler.Service.RegisterVote(ginCtx, pollItem.ID.Hex(), voterToken, req.OptionID)
	if err != nil {
		logger.Error("redis error registering vote", "poll_id", pollItem.ID.Hex(), "error", err)
		response.Error(ginCtx, http.StatusServiceUnavailable, "dependency_unavailable", "voting is temporarily unavailable")
		return
	}
	if !accepted {
		response.Error(ginCtx, http.StatusConflict, "already_voted", "you have already voted")
		return
	}

	// Persist an audit record in MongoDB for auditing and analytics
	if handler.Audit != nil {
		auditCtx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
		defer cancel()
		if _, auditErr := handler.Audit.InsertOne(auditCtx, models.Vote{
			ID:         primitive.NewObjectID(),
			PollID:     pollItem.ID,
			OptionID:   req.OptionID,
			VoterToken: voterToken,
			VotedAt:    time.Now().UTC(),
		}); auditErr != nil {
			logger.Warn("mongodb audit log write failed", "poll_id", pollItem.ID.Hex(), "error", auditErr)
		}
	}

	response.JSON(ginCtx, http.StatusOK, gin.H{"accepted": true, "counts": counts})
}
