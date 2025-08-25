package main

import (
	"log"
	"net/http"
	"os"
	"time"

	authhttp "backend-service-internpro/internal/auth/delivery/http"
	"backend-service-internpro/internal/container"
	"backend-service-internpro/internal/pkg/logger"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humagin"
	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize logger
	logger.InitGlobalLogger(logger.LevelInfo)
	appLogger := logger.Global()

	// Initialize container with all dependencies
	appLogger.Info("initializing application container...")
	c, err := container.NewContainer()
	if err != nil {
		log.Fatal("failed to initialize container:", err)
	}

	// Get port from environment or use default
	port := c.Config.Server.Port

	// Set Gin mode based on environment
	if os.Getenv("GIN_MODE") == "" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Router (Gin) + Huma (OpenAPI runtime)
	r := gin.Default()

	// Add request logging middleware
	r.Use(func(c *gin.Context) {
		start := time.Now()

		appLogger.HTTP().LogRequest(
			c.Request.Method,
			c.Request.URL.Path,
			c.GetHeader("User-Agent"),
			c.ClientIP(),
		)

		c.Next()

		appLogger.HTTP().LogResponse(
			c.Request.Method,
			c.Request.URL.Path,
			c.Writer.Status(),
			time.Since(start),
		)
	})

	api := humagin.New(r, huma.DefaultConfig("Auth API", "1.0.0"))

	// Register routes
	authhttp.New(api, c.AuthService)

	// Health check endpoint
	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"time":    time.Now(),
			"version": "1.0.0",
			"service": "auth-api",
		})
	})

	// Start server
	appLogger.Info("starting server", "port", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		appLogger.ErrorWithErr("server failed to start", err)
		log.Fatal(err)
	}
}
