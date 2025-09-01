package container

import (
	"time"

	"backend-service-internpro/config"
	authRepo "backend-service-internpro/internal/auth/repository"
	authService "backend-service-internpro/internal/auth/service"
	jwtpkg "backend-service-internpro/internal/pkg/jwt"
	"backend-service-internpro/internal/pkg/migration"
	rbacRepo "backend-service-internpro/internal/rbac/repository"
	rbacService "backend-service-internpro/internal/rbac/service"
	schoolRepo "backend-service-internpro/internal/school/repository"
	schoolService "backend-service-internpro/internal/school/service"
	userRepo "backend-service-internpro/internal/user/repository"
	userService "backend-service-internpro/internal/user/service"

	"gorm.io/gorm"
)

// Container holds all dependencies
type Container struct {
	DB            *gorm.DB
	Config        *Config
	AuthRepo      authRepo.Repository
	AuthService   authService.Service
	UserRepo      userRepo.Repository
	UserService   userService.Service
	RBACRepo      rbacRepo.Repository
	RBACService   rbacService.Service
	SchoolRepo    schoolRepo.SchoolRepository
	SchoolService schoolService.SchoolService
	JWTSecrets    jwtpkg.Secrets
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
	authRepository := authRepo.New(db)
	userRepository := userRepo.New(db)
	rbacRepository := rbacRepo.NewRepository(db)
	schoolRepository := schoolRepo.NewSchoolRepository(db)

	// Initialize services with configuration
	authSvc := authService.NewWithConfig(authRepository, jwtSecrets, authService.Config{
		AccessTTL:  cfg.JWT.AccessTokenTTL,
		RefreshTTL: cfg.JWT.RefreshTokenTTL,
	})
	userSvc := userService.New(userRepository)
	rbacSvc := rbacService.NewService(rbacRepository)
	schoolSvc := schoolService.NewSchoolService(schoolRepository)

	return &Container{
		DB:            db,
		Config:        cfg,
		AuthRepo:      authRepository,
		AuthService:   authSvc,
		UserRepo:      userRepository,
		UserService:   userSvc,
		RBACRepo:      rbacRepository,
		RBACService:   rbacSvc,
		SchoolRepo:    schoolRepository,
		SchoolService: schoolSvc,
		JWTSecrets:    jwtSecrets,
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

	// Run migrations automatically in development environment
	if err := migration.AutoMigrateIfDevelopment(db); err != nil {
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
