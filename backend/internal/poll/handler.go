package poll

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"polling-backend/internal/models"
	"polling-backend/pkg/response"
	"strings"
	"time"
)

type LiveService interface {
	HasVoted(context.Context, string, string) (bool, error)
	Counts(context.Context, string) (map[string]int, error)
	PublishClosed(context.Context, string) error
}
type Handler struct {
	Repo                         Repository
	Redis                        *redis.Client
	VoteService                  LiveService
	FrontendBaseURL, VoterCookie string
	Secure                       bool
}
type createInput struct {
	Question  string   `json:"question"`
	Options   []string `json:"options"`
	ExpiresAt *string  `json:"expiresAt"`
}

func idFrom(c *gin.Context) (primitive.ObjectID, bool) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	return id, err == nil
}
func userID(c *gin.Context) (primitive.ObjectID, bool) {
	raw, ok := c.Get("userId")
	if !ok {
		return primitive.NilObjectID, false
	}
	id, err := primitive.ObjectIDFromHex(raw.(string))
	return id, err == nil
}
func view(p *models.Poll, base string) gin.H {
	opts := make([]gin.H, len(p.Options))
	for i, o := range p.Options {
		opts[i] = gin.H{"id": o.ID, "text": o.Text}
	}
	primaryBase := strings.TrimSpace(strings.Split(base, ",")[0])
	if primaryBase == "*" || primaryBase == "" {
		primaryBase = "https://polling.naveenselvan.me"
	}
	return gin.H{"id": p.ID.Hex(), "question": p.Question, "options": opts, "status": p.Status, "expiresAt": p.ExpiresAt, "createdAt": p.CreatedAt, "shareUrl": strings.TrimRight(primaryBase, "/") + "/poll/" + p.ID.Hex()}
}
func (h Handler) Create(c *gin.Context) {
	uid, ok := userID(c)
	if !ok {
		response.Error(c, 401, "unauthorized", "authentication required")
		return
	}
	var in createInput
	if c.ShouldBindJSON(&in) != nil {
		response.Error(c, 400, "validation_error", "invalid request body")
		return
	}
	in.Question = strings.TrimSpace(in.Question)
	if in.Question == "" || len([]rune(in.Question)) > 200 {
		response.Error(c, 400, "validation_error", "question is required and must be at most 200 characters")
		return
	}
	if len(in.Options) < 2 {
		response.Error(c, 400, "validation_error", "at least 2 options are required")
		return
	}
	if len(in.Options) > 6 {
		response.Error(c, 400, "validation_error", "at most 6 options are allowed")
		return
	}
	seen := map[string]bool{}
	opts := make([]models.Option, len(in.Options))
	for i, v := range in.Options {
		v = strings.TrimSpace(v)
		if v == "" {
			response.Error(c, 400, "validation_error", "options cannot be empty")
			return
		}
		k := strings.ToLower(v)
		if seen[k] {
			response.Error(c, 400, "validation_error", "options must be unique")
			return
		}
		seen[k] = true
		opts[i] = models.Option{ID: "opt_" + string(rune('1'+i)), Text: v}
	}
	var expiry *time.Time
	if in.ExpiresAt != nil {
		t, err := time.Parse(time.RFC3339, *in.ExpiresAt)
		if err != nil || !t.After(time.Now()) {
			response.Error(c, 400, "validation_error", "expiresAt must be a future ISO 8601 timestamp")
			return
		}
		expiry = &t
	}
	p := &models.Poll{ID: primitive.NewObjectID(), CreatorID: uid, Question: in.Question, Options: opts, ExpiresAt: expiry}
	if err := h.Repo.Create(c, p); err != nil {
		response.Error(c, 500, "server_error", "something went wrong")
		return
	}
	response.JSON(c, 201, view(p, h.FrontendBaseURL))
}
func (h Handler) Get(c *gin.Context) {
	id, ok := idFrom(c)
	if !ok {
		response.Error(c, 404, "not_found", "poll not found")
		return
	}
	p, err := h.Repo.FindByID(c, id)
	if err != nil {
		response.Error(c, 500, "server_error", "something went wrong")
		return
	}
	if p == nil {
		response.Error(c, 404, "not_found", "poll not found")
		return
	}
	if err = h.Repo.ResolveLazyExpiry(c, p); err != nil {
		response.Error(c, 500, "server_error", "something went wrong")
		return
	}
	token, _ := c.Cookie(h.VoterCookie)
	if token == "" {
		token = c.GetHeader("X-Voter-Token")
	}
	has := false
	if token != "" {
		has, err = h.VoteService.HasVoted(c, p.ID.Hex(), token)
		if err != nil {
			response.Error(c, 503, "dependency_unavailable", "live results are unavailable")
			return
		}
	}
	counts, err := h.VoteService.Counts(c, p.ID.Hex())
	if err != nil {
		response.Error(c, 503, "dependency_unavailable", "live results are unavailable")
		return
	}
	v := view(p, h.FrontendBaseURL)
	v["hasVoted"] = has
	v["results"] = counts
	response.JSON(c, 200, v)
}
func (h Handler) Close(c *gin.Context) {
	uid, ok := userID(c)
	if !ok {
		response.Error(c, 401, "unauthorized", "authentication required")
		return
	}
	id, ok := idFrom(c)
	if !ok {
		response.Error(c, 404, "not_found", "poll not found")
		return
	}
	p, err := h.Repo.FindByID(c, id)
	if err != nil {
		response.Error(c, 500, "server_error", "something went wrong")
		return
	}
	if p == nil {
		response.Error(c, 404, "not_found", "poll not found")
		return
	}
	if p.CreatorID != uid {
		response.Error(c, 403, "forbidden", "not authorized to manage this poll")
		return
	}
	if p.Status == models.StatusOpen {
		if err = h.Repo.SetStatus(c, id, models.StatusClosed); err != nil {
			response.Error(c, 500, "server_error", "something went wrong")
			return
		}
		if err = h.VoteService.PublishClosed(c, id.Hex()); err != nil {
			response.Error(c, 503, "dependency_unavailable", "live results are unavailable")
			return
		}
	}
	response.JSON(c, 200, gin.H{"id": id.Hex(), "status": models.StatusClosed})
}
func (h Handler) ListMine(c *gin.Context) {
	uid, ok := userID(c)
	if !ok {
		response.Error(c, 401, "unauthorized", "authentication required")
		return
	}
	ps, err := h.Repo.FindByCreator(c, uid)
	if err != nil {
		response.Error(c, 500, "server_error", "something went wrong")
		return
	}
	out := make([]gin.H, 0, len(ps))
	for i := range ps {
		if err = h.Repo.ResolveLazyExpiry(c, &ps[i]); err != nil {
			response.Error(c, 500, "server_error", "something went wrong")
			return
		}
		out = append(out, view(&ps[i], h.FrontendBaseURL))
	}
	response.JSON(c, 200, out)
}
func (h Handler) Health(c *gin.Context) { c.Status(http.StatusOK) }
