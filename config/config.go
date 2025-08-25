package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB
var JwtSecret []byte
var JwtExpireTime time.Duration
var RefreshTokenExpire time.Duration
var SmtpHost string
var SmtpPort string
var SmtpUser string
var SmtpPass string

// LoadEnv load .env variables
func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Println("⚠️  No .env file found, using system environment")
	}
}

// InitConfig initialize all configs
func InitConfig() {
	LoadEnv()

	// Database connection
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	// MySQL DSN format
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbUser, dbPass, dbHost, dbPort, dbName)

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("❌ Failed to connect database: ", err)
	}

	// JWT configs
	JwtSecret = []byte(os.Getenv("JWT_SECRET"))

	// Parse JWT expire time from minutes
	jwtExpireMinutes := os.Getenv("JWT_EXPIRE_MINUTES")
	if jwtExpireMinutes == "" {
		jwtExpireMinutes = "15" // default 15 minutes
	}
	JwtExpireTime = time.Duration(parseInt(jwtExpireMinutes)) * time.Minute

	// Parse refresh token expire time from hours
	refreshExpireHours := os.Getenv("JWT_REFRESH_EXPIRE_HOURS")
	if refreshExpireHours == "" {
		refreshExpireHours = "24" // default 24 hours
	}
	RefreshTokenExpire = time.Duration(parseInt(refreshExpireHours)) * time.Hour

	// SMTP configs
	SmtpHost = os.Getenv("SMTP_HOST")
	SmtpPort = os.Getenv("SMTP_PORT")
	SmtpUser = os.Getenv("SMTP_USER")
	SmtpPass = os.Getenv("SMTP_PASS")

	log.Println("✅ Config loaded successfully")
}

// Helper function to parse string to int with default value
func parseInt(s string) int {
	val, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return val
}
