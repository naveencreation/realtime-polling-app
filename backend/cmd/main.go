package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"polling-backend/internal/auth"
	"polling-backend/internal/config"
	"polling-backend/internal/db"
	"polling-backend/internal/middleware"
	"polling-backend/internal/poll"
	"polling-backend/internal/realtime"
	"polling-backend/internal/vote"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	if err := run(logger); err != nil {
		logger.Error("startup failed", "error", err)
		os.Exit(1)
	}
}
func run(logger *slog.Logger) error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}
	root, cancel := context.WithCancel(context.Background())
	defer cancel()
	mongoClient, database, err := db.ConnectMongo(root, cfg.MongoURI, cfg.MongoDBName)
	if err != nil {
		return err
	}
	defer mongoClient.Disconnect(context.Background())
	if err = db.EnsureIndexes(root, database); err != nil {
		return err
	}
	redisClient, err := db.ConnectRedis(root, cfg.RedisAddr)
	if err != nil {
		return err
	}
	defer redisClient.Close()
	authService := auth.Service{Secret: []byte(cfg.JWTSecret), Expiry: time.Duration(cfg.JWTExpiryHours) * time.Hour, Cost: cfg.BcryptCost}
	authHandler := auth.Handler{Repo: auth.Repository{Users: database.Collection("users")}, Service: authService, AuthCookie: cfg.AuthCookieName, Secure: cfg.CookieSecure}
	pollRepo := poll.Repository{Collection: database.Collection("polls")}
	voteService := vote.Service{Redis: redisClient}
	pollHandler := poll.Handler{Repo: pollRepo, VoteService: voteService, FrontendBaseURL: cfg.FrontendBaseURL, VoterCookie: cfg.VoterCookieName, Secure: cfg.CookieSecure}
	voteHandler := vote.Handler{Polls: pollRepo, Service: voteService, Audit: database.Collection("votes"), VoterCookie: cfg.VoterCookieName, Secure: cfg.CookieSecure}
	streamHandler := realtime.Handler{Polls: pollRepo, Redis: redisClient}
	router := gin.New()
	router.Use(middleware.Recovery(logger), gin.Logger())
	corsConfig := cors.Config{
		AllowCredentials: true,
		AllowMethods:     []string{"GET", "POST", "PATCH", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Voter-Token"},
		ExposeHeaders:    []string{"Content-Length", "Set-Cookie", "X-Voter-Token"},
	}
	if cfg.FrontendBaseURL == "*" || cfg.FrontendBaseURL == "" {
		corsConfig.AllowOriginFunc = func(origin string) bool { return true }
	} else {
		rawOrigins := strings.Split(cfg.FrontendBaseURL, ",")
		originSet := make(map[string]bool)
		for _, o := range rawOrigins {
			if trimmed := strings.TrimSpace(o); trimmed != "" {
				originSet[trimmed] = true
			}
		}
		corsConfig.AllowOriginFunc = func(origin string) bool {
			if originSet[origin] {
				return true
			}
			if strings.HasSuffix(origin, ".vercel.app") || strings.HasSuffix(origin, ".naveenselvan.me") || strings.Contains(origin, "localhost") {
				return true
			}
			return false
		}
	}
	router.Use(cors.New(corsConfig))
	router.GET("/healthz", func(c *gin.Context) { c.Status(http.StatusOK) })
	router.GET("/readyz", func(c *gin.Context) {
		if redisClient.Ping(c).Err() != nil || database.Client().Ping(c, nil) != nil {
			c.Status(http.StatusServiceUnavailable)
			return
		}
		c.Status(http.StatusOK)
	})
	api := router.Group("/api")
	api.POST("/auth/signup", authHandler.Signup)
	api.POST("/auth/login", authHandler.Login)
	api.POST("/auth/logout", authHandler.Logout)
	protected := api.Group("", auth.RequireAuth(authService, cfg.AuthCookieName))
	protected.GET("/polls/mine", pollHandler.ListMine)
	protected.POST("/polls", pollHandler.Create)
	protected.PATCH("/polls/:id/close", pollHandler.Close)
	api.GET("/polls/:id", pollHandler.Get)
	api.GET("/polls/:id/stream", streamHandler.Stream)
	api.POST("/polls/:id/vote", voteHandler.Vote)
	server := &http.Server{Addr: ":" + cfg.Port, Handler: router, ReadHeaderTimeout: 5 * time.Second, IdleTimeout: 60 * time.Second}
	go func() {
		if e := server.ListenAndServe(); e != nil && !errors.Is(e, http.ErrServerClosed) {
			logger.Error("server failed", "error", e)
		}
	}()
	logger.Info("server ready", "port", cfg.Port)
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig
	logger.Info("shutdown started")
	shutdownCtx, stop := context.WithTimeout(context.Background(), 10*time.Second)
	defer stop()
	if err = server.Shutdown(shutdownCtx); err != nil {
		logger.Error("forced shutdown", "error", err)
		return err
	}
	logger.Info("shutdown complete")
	return nil
}
