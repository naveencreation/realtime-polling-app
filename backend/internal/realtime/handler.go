package realtime

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"polling-backend/internal/poll"
	"polling-backend/pkg/response"
)

type Handler struct {
	Polls poll.Repository
	Redis *redis.Client
}

func (h Handler) Stream(c *gin.Context) {
	id, ok := primitiveID(c.Param("id"))
	if !ok {
		response.Error(c, 404, "not_found", "poll not found")
		return
	}
	p, err := h.Polls.FindByID(c, id)
	if err != nil {
		response.Error(c, 500, "server_error", "something went wrong")
		return
	}
	if p == nil {
		response.Error(c, 404, "not_found", "poll not found")
		return
	}
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Status(http.StatusOK)
	sub := h.Redis.Subscribe(c, key(id.Hex()))
	defer sub.Close()
	if p.Status == "closed" {
		_, _ = fmt.Fprint(c.Writer, "event: closed\ndata: {\"event\":\"closed\"}\n\n")
		if f, ok := c.Writer.(http.Flusher); ok {
			f.Flush()
		}
		return
	}
	ch := sub.Channel()
	for {
		select {
		case <-c.Request.Context().Done():
			return
		case msg, ok := <-ch:
			if !ok {
				return
			}
			event := "results"
			if msg.Payload == `{"event":"closed"}` {
				event = "closed"
			}
			_, _ = fmt.Fprintf(c.Writer, "event: %s\ndata: %s\n\n", event, msg.Payload)
			if f, ok := c.Writer.(http.Flusher); ok {
				f.Flush()
			}
		}
	}
}
func key(id string) string { return "poll:" + id + ":updates" }
func primitiveID(v string) (primitive.ObjectID, bool) {
	id, err := primitive.ObjectIDFromHex(v)
	return id, err == nil
}
