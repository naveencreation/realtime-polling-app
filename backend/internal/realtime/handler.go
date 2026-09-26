package realtime

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"polling-backend/internal/middleware"
	"polling-backend/internal/models"
	"polling-backend/internal/poll"
	"polling-backend/internal/vote"
	"polling-backend/pkg/response"
)

// Handler manages Server-Sent Events (SSE) connections for live poll results.
type Handler struct {
	Polls poll.Repository
	Redis *redis.Client
}

// Stream handles GET /api/polls/:id/stream.
// It establishes a persistent HTTP SSE connection, streaming live tally updates published
// to Redis. It emits periodic keep-alive comments to prevent intermediate proxies from timing out.
func (handler Handler) Stream(ginCtx *gin.Context) {
	logger := middleware.GetLogger(ginCtx)

	pollID, err := primitive.ObjectIDFromHex(ginCtx.Param("id"))
	if err != nil {
		response.Error(ginCtx, http.StatusNotFound, "not_found", "poll not found")
		return
	}

	pollItem, err := handler.Polls.FindByID(ginCtx, pollID)
	if err != nil {
		logger.Error("database error retrieving poll for SSE stream", "poll_id", pollID.Hex(), "error", err)
		response.Error(ginCtx, http.StatusInternalServerError, "server_error", "something went wrong")
		return
	}
	if pollItem == nil {
		response.Error(ginCtx, http.StatusNotFound, "not_found", "poll not found")
		return
	}

	ginCtx.Header("Content-Type", "text/event-stream")
	ginCtx.Header("Cache-Control", "no-cache")
	ginCtx.Header("Connection", "keep-alive")
	ginCtx.Header("X-Accel-Buffering", "no")
	ginCtx.Status(http.StatusOK)

	flusher, isFlusher := ginCtx.Writer.(http.Flusher)
	ginCtx.Writer.WriteHeaderNow()
	_, _ = fmt.Fprint(ginCtx.Writer, ": connected\n\n")
	if isFlusher {
		flusher.Flush()
	}

	// If the poll is already closed, send terminal event and complete the request
	if pollItem.Status == models.StatusClosed {
		_, _ = fmt.Fprint(ginCtx.Writer, "event: closed\ndata: {\"event\":\"closed\"}\n\n")
		if isFlusher {
			flusher.Flush()
		}
		return
	}

	subscription := handler.Redis.Subscribe(ginCtx, vote.UpdatesKey(pollID.Hex()))
	defer subscription.Close()

	messageChannel := subscription.Channel()
	pingTicker := time.NewTicker(20 * time.Second)
	defer pingTicker.Stop()

	for {
		select {
		case <-ginCtx.Request.Context().Done():
			return

		case <-pingTicker.C:
			// Emit periodic SSE comment to keep the connection open across reverse proxies and firewalls
			_, _ = fmt.Fprint(ginCtx.Writer, ": ping\n\n")
			if isFlusher {
				flusher.Flush()
			}

		case pubsubMessage, ok := <-messageChannel:
			if !ok {
				return
			}
			eventName := "results"
			if pubsubMessage.Payload == `{"event":"closed"}` {
				eventName = "closed"
			}
			_, _ = fmt.Fprintf(ginCtx.Writer, "event: %s\ndata: %s\n\n", eventName, pubsubMessage.Payload)
			if isFlusher {
				flusher.Flush()
			}
		}
	}
}
