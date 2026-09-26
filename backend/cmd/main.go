package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"syscall"
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

	rootCtx, rootCancel := context.WithCancel(context.Background())
	defer rootCancel()

	// Initialize MongoDB client and ensure required uniqueness indexes
	mongoClient, database, err := db.ConnectMongo(rootCtx, cfg.MongoURI, cfg.MongoDBName)
	if err != nil {
		return err
	}
	defer mongoClient.Disconnect(context.Background())

	if err = db.EnsureIndexes(rootCtx, database); err != nil {
		return err
	}

	// Initialize Redis connection for fast voting tallies and pub/sub
	redisClient, err := db.ConnectRedis(rootCtx, cfg.RedisAddr)
	if err != nil {
		return err
	}
	defer redisClient.Close()

	// Instantiate domain services and handlers
	authService := auth.Service{
		Secret: []byte(cfg.JWTSecret),
		Expiry: time.Duration(cfg.JWTExpiryHours) * time.Hour,
		Cost:   cfg.BcryptCost,
	}

	authHandler := auth.Handler{
		Repo:       auth.Repository{Users: database.Collection("users")},
		Service:    authService,
		AuthCookie: cfg.AuthCookieName,
		Secure:     cfg.CookieSecure,
	}

	pollRepo := poll.Repository{Collection: database.Collection("polls")}
	voteService := vote.Service{Redis: redisClient}

	pollHandler := poll.Handler{
		Repo:            pollRepo,
		VoteService:     voteService,
		Audit:           database.Collection("votes"),
		FrontendBaseURL: cfg.FrontendBaseURL,
		VoterCookie:     cfg.VoterCookieName,
		Secure:          cfg.CookieSecure,
	}

	voteHandler := vote.Handler{
		Polls:       pollRepo,
		Service:     voteService,
		Audit:       database.Collection("votes"),
		VoterCookie: cfg.VoterCookieName,
		Secure:      cfg.CookieSecure,
	}

	streamHandler := realtime.Handler{
		Polls: pollRepo,
		Redis: redisClient,
	}

	router := gin.New()
	router.Use(
		middleware.Recovery(logger),
		middleware.RequestLogger(logger),
	)

	// Configure cross-origin resource sharing for frontend integration
	corsConfig := cors.Config{
		AllowCredentials: true,
		AllowMethods:     []string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Voter-Token", "X-Request-ID"},
		ExposeHeaders:    []string{"Content-Length", "Set-Cookie", "X-Voter-Token", "X-Request-ID"},
	}

	if cfg.FrontendBaseURL == "*" || cfg.FrontendBaseURL == "" {
		corsConfig.AllowOriginFunc = func(origin string) bool { return true }
	} else {
		rawOrigins := strings.Split(cfg.FrontendBaseURL, ",")
		originSet := make(map[string]bool, len(rawOrigins))
		for _, rawOrigin := range rawOrigins {
			if trimmedOrigin := strings.TrimSpace(rawOrigin); trimmedOrigin != "" {
				originSet[trimmedOrigin] = true
			}
		}
		corsConfig.AllowOriginFunc = func(origin string) bool {
			if originSet[origin] {
				return true
			}
			if strings.HasSuffix(origin, ".vercel.app") || strings.HasSuffix(origin, ".naveenselvan.me") {
				return true
			}
			parsedURL, parseErr := url.Parse(origin)
			if parseErr == nil {
				hostname := parsedURL.Hostname()
				if hostname == "localhost" || hostname == "127.0.0.1" {
					return true
				}
			}
			return false
		}
	}
	router.Use(cors.New(corsConfig))

	// Liveness and readiness endpoints for deployment health checks
	router.GET("/healthz", func(ginCtx *gin.Context) { ginCtx.Status(http.StatusOK) })
	router.GET("/readyz", func(ginCtx *gin.Context) {
		if redisClient.Ping(ginCtx).Err() != nil || database.Client().Ping(ginCtx, nil) != nil {
			ginCtx.Status(http.StatusServiceUnavailable)
			return
		}
		ginCtx.Status(http.StatusOK)
	})

	// Public auth endpoints
	api := router.Group("/api")
	api.POST("/auth/signup", authHandler.Signup)
	api.POST("/auth/login", authHandler.Login)
	api.POST("/auth/logout", authHandler.Logout)

	// Protected author routes (requires JWT)
	protected := api.Group("", auth.RequireAuth(authService, cfg.AuthCookieName))
	protected.GET("/polls/mine", pollHandler.ListMine)
	protected.POST("/polls", pollHandler.Create)
	protected.PATCH("/polls/:id/close", pollHandler.Close)
	protected.DELETE("/polls/:id", pollHandler.Delete)

	// Public poll view, vote, and real-time streaming routes
	api.GET("/polls/:id", pollHandler.Get)
	api.GET("/polls/:id/stream", streamHandler.Stream)
	api.POST("/polls/:id/vote", voteHandler.Vote)

	server := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	go func() {
		if listenErr := server.ListenAndServe(); listenErr != nil && !errors.Is(listenErr, http.ErrServerClosed) {
			logger.Error("server failed", "error", listenErr)
		}
	}()
	logger.Info("server ready", "port", cfg.Port)

	// Block until an OS interrupt/termination signal is received
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	<-signalChan

	logger.Info("shutdown started")
	shutdownCtx, stopCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer stopCancel()

	if err = server.Shutdown(shutdownCtx); err != nil {
		logger.Error("forced shutdown", "error", err)
		return err
	}
	logger.Info("shutdown complete")
	return nil
}
