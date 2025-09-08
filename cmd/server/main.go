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
	rbachttp "backend-service-internpro/internal/rbac/delivery/http"
	schoolhttp "backend-service-internpro/internal/school/delivery/http"
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

	// Add middlewares in proper order
	r.Use(middleware.CORSMiddleware()) // CORS first
	r.Use(middleware.RecoveryMiddleware())
	r.Use(middleware.SecurityHeadersMiddleware())
	r.Use(middleware.FormDataToJSONMiddleware()) // Add FormData support
	r.Use(middleware.LoggingMiddleware())
	r.Use(middleware.RateLimitMiddleware(time.Second, 100)) // 100 requests per second per IP

	// Configure Huma with detailed OpenAPI documentation
	config := huma.DefaultConfig("Gapura SchoolTech API", "0.0.1")
	config.OpenAPI.Info.Description = "Dokumentasi API untuk platform SchoolTech. Ini mencakup endpoint untuk autentikasi, kepentingan internal SchoolTech Indonesia, dan kepentingan produk SchoolTech Indonesia."
	config.OpenAPI.Info.Contact = &huma.Contact{
		Name:  "ITDB SchoolTech",
		Email: "itdb@schooltechindonesia.com",
		URL:   "https://schooltechindonesia.com",
	}
	config.OpenAPI.Info.License = &huma.License{
		Name: "MIT",
		URL:  "https://opensource.org/licenses/MIT",
	}
	config.OpenAPI.Info.TermsOfService = "https://schooltechindonesia.com/terms"
	config.OpenAPI.Servers = []*huma.Server{
		{
			URL:         "https://api.schooltechindonesia.com",
			Description: "Production server",
		},
		{
			URL:         "https://staging-api.schooltechindonesia.com",
			Description: "Staging server",
		},
		{
			URL:         "https://testing-api.schooltechindonesia.com",
			Description: "Testing server",
		},
		{
			URL:         "http://localhost:" + port,
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
			Name:        "RBAC - Roles",
			Description: "Endpoint untuk manajemen roles (peran) dalam sistem RBAC",
		},
		{
			Name:        "RBAC - Permissions",
			Description: "Endpoint untuk manajemen permissions (izin) dalam sistem RBAC",
		},
		{
			Name:        "RBAC - Menus",
			Description: "Endpoint untuk manajemen menus dalam sistem RBAC",
		},
		{
			Name:        "RBAC - User Roles",
			Description: "Endpoint untuk assignment dan manajemen roles pengguna",
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
	userhttp.New(api, c.UserService, c.JWTSecrets)     // User management routes
	rbachttp.NewHuma(api, c.RBACService, c.JWTSecrets) // RBAC management routes with Swagger
	schoolhttp.New(api, c.SchoolService, c.JWTSecrets) // School management routes

	// Health check endpoint
	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"time":    time.Now(),
			"version": "0.0.1",
			"service": "Schooltech API Service",
		})
	})

	// CORS test endpoint
	r.GET("/cors-test", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "CORS is working correctly!",
			"origin":  c.Request.Header.Get("Origin"),
			"method":  c.Request.Method,
			"headers": c.Request.Header,
			"time":    time.Now(),
		})
	})

	// Serve static files for CORS testing
	r.Static("/static", "./static")
	r.StaticFile("/test-cors", "./test-cors.html")

	// Start server
	appLogger.Info("starting server", "port", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		appLogger.ErrorWithErr("server failed to start", err)
		log.Fatal(err)
	}
}
