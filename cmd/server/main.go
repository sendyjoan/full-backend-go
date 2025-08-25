package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"backend-service-internpro/config"
	authhttp "backend-service-internpro/internal/auth/delivery/http"
	"backend-service-internpro/internal/auth/repository"
	"backend-service-internpro/internal/auth/service"
	jwtpkg "backend-service-internpro/internal/pkg/jwt"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humagin"
	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize configuration (loads .env and connects to database)
	config.InitConfig()

	// Get port from environment or use default
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	// Router (Gin) + Huma (OpenAPI runtime)
	r := gin.Default()

	api := humagin.New(r, huma.DefaultConfig("Auth API", "1.0.0"))
	// (opsional) “first hit to generate”: OpenAPI disajikan saat /docs atau /openapi.json pertama kali diakses.
	// Huma otomatis serve docs di /docs & spec di /openapi.json.

	// DI
	repo := repository.New(config.DB)
	secrets := jwtpkg.Secrets{Access: config.JwtSecret, Refresh: config.JwtSecret}
	svc := service.New(repo, secrets)
	authhttp.New(api, svc)

	// Health
	r.GET("/healthz", func(c *gin.Context) { c.JSON(200, gin.H{"status": "ok", "time": time.Now()}) })

	// Start
	log.Printf("🚀 Server listening on port %s", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatal(err)
	}
}
