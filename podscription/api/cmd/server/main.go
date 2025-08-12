package main

import (
	"fmt"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"podscription-api/controllers"
	"podscription-api/internal/handlers"
	"podscription-api/internal/managers"
	"podscription-api/internal/store"
	"podscription-api/pkg/config"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Setup logger
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetLevel(logrus.InfoLevel)

	// Validate OpenAI API key
	if cfg.OpenAI.APIKey == "" {
		logger.Fatal("OPENAI_API_KEY environment variable is required")
		os.Exit(1)
	}

	logger.WithFields(logrus.Fields{
		"host":         cfg.Server.Host,
		"port":         cfg.Server.Port,
		"openai_model": cfg.OpenAI.Model,
		"store_type":   cfg.Store.Type,
	}).Info("starting podscription API server")

	// Initialize store
	var dataStore store.Store
	switch cfg.Store.Type {
	case "memory":
		dataStore = store.NewMemoryStore(cfg.Store.Path)
	default:
		logger.WithField("store_type", cfg.Store.Type).Fatal("unsupported store type")
		os.Exit(1)
	}

	// Initialize managers
	openAIManager := managers.NewOpenAIManager(cfg.OpenAI)
	sessionManager := managers.NewSessionManager(dataStore, openAIManager, logger)

	// Initialize controllers
	chatController := controllers.NewChatController(sessionManager, logger)

	// Initialize handlers
	chatHandler := handlers.NewChatHandler(chatController, logger)

	// Setup Gin router
	router := setupRouter(chatHandler, logger, cfg)

	// Start server
	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	logger.WithField("address", addr).Info("server starting")
	
	if err := router.Run(addr); err != nil {
		logger.WithError(err).Fatal("failed to start server")
		os.Exit(1)
	}
}

func setupRouter(chatHandler *handlers.ChatHandler, logger *logrus.Logger, cfg *config.Config) *gin.Engine {
	// Set Gin mode
	gin.SetMode(gin.ReleaseMode)
	
	router := gin.New()

	// Middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(corsMiddleware())
	router.Use(loggingMiddleware(logger))

	// Health check
	router.GET("/health", chatHandler.HealthCheck)

	// API routes
	api := router.Group("/api")
	{
		// Chat endpoints
		api.POST("/chat", chatHandler.SendMessage)
		
		// Session endpoints
		api.POST("/sessions", chatHandler.CreateSession)
		api.GET("/sessions", chatHandler.ListSessions)
		api.GET("/sessions/:id", chatHandler.GetSession)
	}

	// Serve static files (for potential future use)
	router.Static("/static", "./static")

	return router
}

func corsMiddleware() gin.HandlerFunc {
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Authorization", "Accept", "X-Requested-With"}
	config.ExposeHeaders = []string{"Content-Length"}
	config.AllowCredentials = true
	
	return cors.New(config)
}

func loggingMiddleware(logger *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Process request
		c.Next()

		// Log request
		latency := time.Since(start)
		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()

		if raw != "" {
			path = path + "?" + raw
		}

		logger.WithFields(logrus.Fields{
			"status_code": statusCode,
			"latency":     latency.Milliseconds(),
			"client_ip":   clientIP,
			"method":      method,
			"path":        path,
		}).Info("request processed")
	}
}