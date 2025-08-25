package container

import (
	"time"

	"backend-service-internpro/config"
	"backend-service-internpro/internal/auth/repository"
	"backend-service-internpro/internal/auth/service"
	jwtpkg "backend-service-internpro/internal/pkg/jwt"
	"backend-service-internpro/internal/pkg/migration"

	"gorm.io/gorm"
)

// Container holds all dependencies
type Container struct {
	DB          *gorm.DB
	Config      *Config
	AuthRepo    repository.Repository
	AuthService service.Service
	JWTSecrets  jwtpkg.Secrets
}

// Config holds all configuration values
type Config struct {
	Server ServerConfig
	JWT    JWTConfig
	SMTP   SMTPConfig
}

type ServerConfig struct {
	Port string
}

type JWTConfig struct {
	AccessSecret    []byte
	RefreshSecret   []byte
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
}

type SMTPConfig struct {
	Host string
	Port string
	User string
	Pass string
}

// NewContainer creates and initializes all dependencies
func NewContainer() (*Container, error) {
	// Load configuration
	cfg, err := loadConfig()
	if err != nil {
		return nil, err
	}

	// Initialize database
	db, err := initDatabase()
	if err != nil {
		return nil, err
	}

	// Create JWT secrets
	jwtSecrets := jwtpkg.Secrets{
		Access:  cfg.JWT.AccessSecret,
		Refresh: cfg.JWT.RefreshSecret,
	}

	// Initialize repositories
	authRepo := repository.New(db)

	// Initialize services with configuration
	authService := service.NewWithConfig(authRepo, jwtSecrets, service.Config{
		AccessTTL:  cfg.JWT.AccessTokenTTL,
		RefreshTTL: cfg.JWT.RefreshTokenTTL,
	})

	return &Container{
		DB:          db,
		Config:      cfg,
		AuthRepo:    authRepo,
		AuthService: authService,
		JWTSecrets:  jwtSecrets,
	}, nil
}

func loadConfig() (*Config, error) {
	// Initialize legacy config for now
	config.InitConfig()

	return &Config{
		Server: ServerConfig{
			Port: getEnvWithDefault("APP_PORT", "8080"),
		},
		JWT: JWTConfig{
			AccessSecret:    config.JwtSecret,
			RefreshSecret:   config.JwtSecret, // Could be different
			AccessTokenTTL:  config.JwtExpireTime,
			RefreshTokenTTL: config.RefreshTokenExpire,
		},
		SMTP: SMTPConfig{
			Host: config.SmtpHost,
			Port: config.SmtpPort,
			User: config.SmtpUser,
			Pass: config.SmtpPass,
		},
	}, nil
}

func initDatabase() (*gorm.DB, error) {
	db := config.DB

	// Run migrations
	if err := migration.AutoMigrate(db); err != nil {
		return nil, err
	}

	// Create dummy data if needed (for development)
	if err := migration.CreateDummyData(db); err != nil {
		return nil, err
	}

	return db, nil
}

func getEnvWithDefault(key, defaultValue string) string {
	if value := config.LoadEnvVar(key); value != "" {
		return value
	}
	return defaultValue
}
