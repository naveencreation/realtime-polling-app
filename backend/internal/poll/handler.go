package poll

import (
	"context"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"polling-backend/internal/middleware"
	"polling-backend/internal/models"
	"polling-backend/pkg/response"
)

// LiveService defines the voting operations required by the poll handler.
type LiveService interface {
	HasVoted(ctx context.Context, pollID, voterToken string) (bool, error)
	Counts(ctx context.Context, pollID string) (map[string]int, error)
	PublishClosed(ctx context.Context, pollID string) error
}

// Handler coordinates poll creation, inspection, closure, and author listings.
type Handler struct {
	Repo            Repository
	VoteService     LiveService
	FrontendBaseURL string
	VoterCookie     string
	Secure          bool
}

type createInput struct {
	Question  string   `json:"question"`
	Options   []string `json:"options"`
	ExpiresAt *string  `json:"expiresAt"`
}

func parsePollID(ginCtx *gin.Context) (primitive.ObjectID, bool) {
	pollID, err := primitive.ObjectIDFromHex(ginCtx.Param("id"))
	return pollID, err == nil
}

func parseUserID(ginCtx *gin.Context) (primitive.ObjectID, bool) {
	rawVal, ok := ginCtx.Get("userId")
	if !ok {
		return primitive.NilObjectID, false
	}
	userIDStr, ok := rawVal.(string)
	if !ok {
		return primitive.NilObjectID, false
	}
	userID, err := primitive.ObjectIDFromHex(userIDStr)
	return userID, err == nil
}

func formatPollResponse(pollItem *models.Poll, frontendBaseURL string) gin.H {
	formattedOptions := make([]gin.H, len(pollItem.Options))
	for index, option := range pollItem.Options {
		formattedOptions[index] = gin.H{"id": option.ID, "text": option.Text}
	}

	primaryBase := strings.TrimSpace(strings.Split(frontendBaseURL, ",")[0])
	if primaryBase == "*" || primaryBase == "" {
		primaryBase = "http://localhost:5173"
	}

	return gin.H{
		"id":        pollItem.ID.Hex(),
		"question":  pollItem.Question,
		"options":   formattedOptions,
		"status":    pollItem.Status,
		"expiresAt": pollItem.ExpiresAt,
		"createdAt": pollItem.CreatedAt,
		"shareUrl":  strings.TrimRight(primaryBase, "/") + "/poll/" + pollItem.ID.Hex(),
	}
}

// Create handles POST /api/polls. Validates question length and option constraints (2-6 unique options).
func (handler Handler) Create(ginCtx *gin.Context) {
	creatorID, ok := parseUserID(ginCtx)
	if !ok {
		response.Error(ginCtx, http.StatusUnauthorized, "unauthorized", "authentication required")
		return
	}

	var req createInput
	if err := ginCtx.ShouldBindJSON(&req); err != nil {
		response.Error(ginCtx, http.StatusBadRequest, "validation_error", "invalid request body")
		return
	}

	req.Question = strings.TrimSpace(req.Question)
	if req.Question == "" || len([]rune(req.Question)) > 200 {
		response.Error(ginCtx, http.StatusBadRequest, "validation_error", "question is required and must be at most 200 characters")
		return
	}

	if len(req.Options) < 2 {
		response.Error(ginCtx, http.StatusBadRequest, "validation_error", "at least 2 options are required")
		return
	}
	if len(req.Options) > 6 {
		response.Error(ginCtx, http.StatusBadRequest, "validation_error", "at most 6 options are allowed")
		return
	}

	seenOptions := make(map[string]bool, len(req.Options))
	pollOptions := make([]models.Option, len(req.Options))
	for index, optionText := range req.Options {
		optionText = strings.TrimSpace(optionText)
		if optionText == "" {
			response.Error(ginCtx, http.StatusBadRequest, "validation_error", "options cannot be empty")
			return
		}
		loweredKey := strings.ToLower(optionText)
		if seenOptions[loweredKey] {
			response.Error(ginCtx, http.StatusBadRequest, "validation_error", "options must be unique")
			return
		}
		seenOptions[loweredKey] = true
		pollOptions[index] = models.Option{
			ID:   "opt_" + strconv.Itoa(index+1),
			Text: optionText,
		}
	}

	var expiryTime *time.Time
	if req.ExpiresAt != nil {
		parsedTime, err := time.Parse(time.RFC3339, *req.ExpiresAt)
		if err != nil || !parsedTime.After(time.Now()) {
			response.Error(ginCtx, http.StatusBadRequest, "validation_error", "expiresAt must be a future ISO 8601 timestamp")
			return
		}
		expiryTime = &parsedTime
	}

	newPoll := &models.Poll{
		ID:        primitive.NewObjectID(),
		CreatorID: creatorID,
		Question:  req.Question,
		Options:   pollOptions,
		ExpiresAt: expiryTime,
	}

	if err := handler.Repo.Create(ginCtx, newPoll); err != nil {
		middleware.GetLogger(ginCtx).Error("database error creating poll", "error", err)
		response.Error(ginCtx, http.StatusInternalServerError, "server_error", "something went wrong")
		return
	}

	response.JSON(ginCtx, http.StatusCreated, formatPollResponse(newPoll, handler.FrontendBaseURL))
}

// Get handles GET /api/polls/:id. Resolves lazy expiry, checks voter status, and attaches live tallies.
func (handler Handler) Get(ginCtx *gin.Context) {
	pollID, ok := parsePollID(ginCtx)
	if !ok {
		response.Error(ginCtx, http.StatusNotFound, "not_found", "poll not found")
		return
	}

	pollItem, err := handler.Repo.FindByID(ginCtx, pollID)
	if err != nil {
		middleware.GetLogger(ginCtx).Error("database error finding poll", "poll_id", pollID.Hex(), "error", err)
		response.Error(ginCtx, http.StatusInternalServerError, "server_error", "something went wrong")
		return
	}
	if pollItem == nil {
		response.Error(ginCtx, http.StatusNotFound, "not_found", "poll not found")
		return
	}

	if err = handler.Repo.ResolveLazyExpiry(ginCtx, pollItem); err != nil {
		middleware.GetLogger(ginCtx).Error("database error resolving poll expiry", "poll_id", pollID.Hex(), "error", err)
		response.Error(ginCtx, http.StatusInternalServerError, "server_error", "something went wrong")
		return
	}

	voterToken, _ := ginCtx.Cookie(handler.VoterCookie)
	if voterToken == "" {
		voterToken = ginCtx.GetHeader("X-Voter-Token")
	}

	hasVoted := false
	if voterToken != "" {
		hasVoted, err = handler.VoteService.HasVoted(ginCtx, pollItem.ID.Hex(), voterToken)
		if err != nil {
			middleware.GetLogger(ginCtx).Error("redis error checking voter status", "poll_id", pollItem.ID.Hex(), "error", err)
			response.Error(ginCtx, http.StatusServiceUnavailable, "dependency_unavailable", "live results are unavailable")
			return
		}
	}

	counts, err := handler.VoteService.Counts(ginCtx, pollItem.ID.Hex())
	if err != nil {
		middleware.GetLogger(ginCtx).Error("redis error fetching vote counts", "poll_id", pollItem.ID.Hex(), "error", err)
		response.Error(ginCtx, http.StatusServiceUnavailable, "dependency_unavailable", "live results are unavailable")
		return
	}

	responsePayload := formatPollResponse(pollItem, handler.FrontendBaseURL)
	responsePayload["hasVoted"] = hasVoted
	responsePayload["results"] = counts
	response.JSON(ginCtx, http.StatusOK, responsePayload)
}

// Close handles PATCH /api/polls/:id/close. Only the poll creator can close an active poll.
func (handler Handler) Close(ginCtx *gin.Context) {
	userID, ok := parseUserID(ginCtx)
	if !ok {
		response.Error(ginCtx, http.StatusUnauthorized, "unauthorized", "authentication required")
		return
	}

	pollID, ok := parsePollID(ginCtx)
	if !ok {
		response.Error(ginCtx, http.StatusNotFound, "not_found", "poll not found")
		return
	}

	pollItem, err := handler.Repo.FindByID(ginCtx, pollID)
	if err != nil {
		middleware.GetLogger(ginCtx).Error("database error finding poll to close", "poll_id", pollID.Hex(), "error", err)
		response.Error(ginCtx, http.StatusInternalServerError, "server_error", "something went wrong")
		return
	}
	if pollItem == nil {
		response.Error(ginCtx, http.StatusNotFound, "not_found", "poll not found")
		return
	}

	if pollItem.CreatorID != userID {
		response.Error(ginCtx, http.StatusForbidden, "forbidden", "not authorized to manage this poll")
		return
	}

	if pollItem.Status == models.StatusOpen {
		if err = handler.Repo.SetStatus(ginCtx, pollID, models.StatusClosed); err != nil {
			middleware.GetLogger(ginCtx).Error("database error closing poll", "poll_id", pollID.Hex(), "error", err)
			response.Error(ginCtx, http.StatusInternalServerError, "server_error", "something went wrong")
			return
		}
		if err = handler.VoteService.PublishClosed(ginCtx, pollID.Hex()); err != nil {
			middleware.GetLogger(ginCtx).Error("redis error publishing poll closure", "poll_id", pollID.Hex(), "error", err)
			response.Error(ginCtx, http.StatusServiceUnavailable, "dependency_unavailable", "live results are unavailable")
			return
		}
	}

	response.JSON(ginCtx, http.StatusOK, gin.H{"id": pollID.Hex(), "status": models.StatusClosed})
}

// ListMine handles GET /api/polls/mine. Returns all polls authored by the authenticated user.
func (handler Handler) ListMine(ginCtx *gin.Context) {
	creatorID, ok := parseUserID(ginCtx)
	if !ok {
		response.Error(ginCtx, http.StatusUnauthorized, "unauthorized", "authentication required")
		return
	}

	polls, err := handler.Repo.FindByCreator(ginCtx, creatorID)
	if err != nil {
		middleware.GetLogger(ginCtx).Error("database error listing user polls", "creator_id", creatorID.Hex(), "error", err)
		response.Error(ginCtx, http.StatusInternalServerError, "server_error", "something went wrong")
		return
	}

	pollViews := make([]gin.H, 0, len(polls))
	for index := range polls {
		if err = handler.Repo.ResolveLazyExpiry(ginCtx, &polls[index]); err != nil {
			middleware.GetLogger(ginCtx).Error("database error resolving poll expiry in list", "poll_id", polls[index].ID.Hex(), "error", err)
			response.Error(ginCtx, http.StatusInternalServerError, "server_error", "something went wrong")
			return
		}
		pollViews = append(pollViews, formatPollResponse(&polls[index], handler.FrontendBaseURL))
	}

	response.JSON(ginCtx, http.StatusOK, pollViews)
}
