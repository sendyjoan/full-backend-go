package main

import (
	"log"
	"net/http"
	"os"
	"time"

	authhttp "backend-service-internpro/internal/auth/delivery/http"
	"backend-service-internpro/internal/container"
	"backend-service-internpro/internal/pkg/logger"
	"backend-service-internpro/internal/pkg/middleware"
	userhttp "backend-service-internpro/internal/user/delivery/http"

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

	// Add middlewares
	r.Use(middleware.RecoveryMiddleware())
	r.Use(middleware.SecurityHeadersMiddleware())
	r.Use(middleware.CORSMiddleware())
	r.Use(middleware.FormDataToJSONMiddleware()) // Add FormData support
	r.Use(middleware.LoggingMiddleware())
	r.Use(middleware.RateLimitMiddleware(time.Second, 100)) // 100 requests per second per IP

	// Configure Huma with detailed OpenAPI documentation
	config := huma.DefaultConfig("SchoolTech Apps API", "1.0.0")
	config.OpenAPI.Info.Description = "Dokumentasi API untuk platform SchoolTech. Ini mencakup endpoint untuk autentikasi, kepentingan internal SchoolTech Indonesia, dan kepentingan produk SchoolTech Indonesia."
	config.OpenAPI.Info.Contact = &huma.Contact{
		Name:  "Tim Developer SchoolTech",
		Email: "dev@schooltech.id",
		URL:   "https://schooltech.id",
	}
	config.OpenAPI.Info.License = &huma.License{
		Name: "MIT",
		URL:  "https://opensource.org/licenses/MIT",
	}
	config.OpenAPI.Info.TermsOfService = "https://schooltech.id/terms"
	config.OpenAPI.Servers = []*huma.Server{
		{
			URL:         "https://api.schooltech.id",
			Description: "Production server",
		},
		{
			URL:         "http://localhost:8080",
			Description: "Local development",
		},
	}

	// Initialize components if not exists
	if config.OpenAPI.Components == nil {
		config.OpenAPI.Components = &huma.Components{}
	}
	if config.OpenAPI.Components.SecuritySchemes == nil {
		config.OpenAPI.Components.SecuritySchemes = make(map[string]*huma.SecurityScheme)
	}

	// Add security schemes
	config.OpenAPI.Components.SecuritySchemes["bearerAuth"] = &huma.SecurityScheme{
		Type:         "http",
		Scheme:       "bearer",
		BearerFormat: "JWT",
	}

	// Add API tags for better organization
	config.OpenAPI.Tags = []*huma.Tag{
		{
			Name:        "Authentication",
			Description: "Endpoint untuk autentikasi pengguna, login, logout, dan manajemen token",
		},
		{
			Name:        "User Management",
			Description: "Endpoint untuk manajemen data pengguna",
		},
		{
			Name:        "School Management",
			Description: "Endpoint untuk manajemen data sekolah",
		},
		{
			Name:        "Teacher Management",
			Description: "Endpoint untuk manajemen data guru",
		},
		{
			Name:        "Student Management",
			Description: "Endpoint untuk manajemen data siswa",
		},
	}

	api := humagin.New(r, config)

	// Register routes
	authhttp.New(api, c.AuthService)
	userhttp.New(api, c.UserService, c.JWTSecrets) // User management routes

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
