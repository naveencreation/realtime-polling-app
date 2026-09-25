package vote

import (
	"context"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"polling-backend/internal/auth"
	"polling-backend/internal/models"
	"polling-backend/internal/poll"
	"polling-backend/pkg/response"
	"strings"
	"time"
)

type Handler struct {
	Polls       poll.Repository
	Service     Service
	Audit       *mongo.Collection
	VoterCookie string
	Secure      bool
}
type input struct {
	OptionID string `json:"optionId"`
}

func (h Handler) Vote(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		response.Error(c, 404, "not_found", "poll not found")
		return
	}
	p, err := h.Polls.FindByID(c, id)
	if err != nil || p == nil {
		response.Error(c, 404, "not_found", "poll not found")
		return
	}
	if err = h.Polls.ResolveLazyExpiry(c, p); err != nil {
		response.Error(c, 500, "server_error", "something went wrong")
		return
	}
	if p.Status != models.StatusOpen {
		response.Error(c, 409, "poll_closed", "poll is closed")
		return
	}
	var in input
	if c.ShouldBindJSON(&in) != nil || strings.TrimSpace(in.OptionID) == "" {
		response.Error(c, 400, "validation_error", "optionId is required")
		return
	}
	valid := false
	for _, o := range p.Options {
		if o.ID == in.OptionID {
			valid = true
			break
		}
	}
	if !valid {
		response.Error(c, 400, "invalid_option", "option is not valid for this poll")
		return
	}
	token, err := c.Cookie(h.VoterCookie)
	if err != nil || token == "" {
		token = auth.NewVoterToken()
		http.SetCookie(c.Writer, &http.Cookie{Name: h.VoterCookie, Value: token, MaxAge: 31536000, HttpOnly: true, Secure: h.Secure, SameSite: http.SameSiteLaxMode, Path: "/"})
	}
	accepted, counts, err := h.Service.RegisterVote(c, p.ID.Hex(), token, in.OptionID)
	if err != nil {
		response.Error(c, 503, "dependency_unavailable", "voting is temporarily unavailable")
		return
	}
	if !accepted {
		response.Error(c, 409, "already_voted", "you have already voted")
		return
	}
	if h.Audit != nil {
		auditCtx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
		_, _ = h.Audit.InsertOne(auditCtx, models.Vote{ID: primitive.NewObjectID(), PollID: p.ID, OptionID: in.OptionID, VoterToken: token, VotedAt: time.Now().UTC()})
		cancel()
	}
	response.JSON(c, 200, gin.H{"accepted": true, "counts": counts})
}
