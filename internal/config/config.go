package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBUser             string `mapstructure:"DB_USER"`
	DBPassword         string `mapstructure:"DB_PASSWORD"`
	DBHost             string `mapstructure:"DB_HOST"`
	DBPort             string `mapstructure:"DB_PORT"`
	DBName             string `mapstructure:"DB_NAME"`
	ServerPort         string `mapstructure:"SERVER_PORT"`
	JWTSecret          string `mapstructure:"JWT_SECRET"`
	SEEDER_ADMIN_PASS  string `mapstructure:"SEEDER_ADMIN_PASS"`
	GoogleClientID     string `mapstructure:"GOOGLE_CLIENT_ID"`
	GoogleClientSecret string `mapstructure:"GOOGLE_CLIENT_SECRET"`
	GoogleRedirectURL  string `mapstructure:"GOOGLE_REDIRECT_URL"`
	SMTPHost           string `mapstructure:"SMTP_HOST"`
	SMTPPort           int    `mapstructure:"SMTP_PORT"`
	SMTPUser           string `mapstructure:"SMTP_USER"`
	SMTPPassword       string `mapstructure:"SMTP_PASSWORD"`
	SMTPSenderEmail    string `mapstructure:"SMTP_SENDER_EMAIL"`
	FrontendBaseURL    string `mapstructure:"FRONTEND_BASE_URL"`
}

var AppConfig *Config

// LoadConfig loads environment variables from .env file
func LoadConfig() {
	// Load .env file if exists
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	AppConfig = &Config{
		DBUser:             getEnv("DB_USER", "root"),
		DBPassword:         getEnv("DB_PASSWORD", ""),
		DBHost:             getEnv("DB_HOST", "localhost"),
		DBPort:             getEnv("DB_PORT", "3306"),
		DBName:             getEnv("DB_NAME", "terra_db"),
		ServerPort:         getEnv("SERVER_PORT", "8080"),
		JWTSecret:          getEnv("JWT_SECRET", "your-secret-key"),
		SEEDER_ADMIN_PASS:  getEnv("SEEDER_ADMIN_PASS", "password123"),
		GoogleClientID:     getEnv("GOOGLE_CLIENT_ID", ""),
		GoogleClientSecret: getEnv("GOOGLE_CLIENT_SECRET", ""),
		GoogleRedirectURL:  getEnv("GOOGLE_REDIRECT_URL", "http://localhost:8080/auth/google/callback"),
		SMTPHost:           getEnv("SMTP_HOST", "smtp.gmail.com"),
		SMTPPort:           getEnvAsInt("SMTP_PORT", 587),
		SMTPUser:           getEnv("SMTP_USER", ""),
		SMTPPassword:       getEnv("SMTP_PASSWORD", ""),
		SMTPSenderEmail:    getEnv("SMTP_SENDER_EMAIL", ""),
		FrontendBaseURL:    getEnv("FRONTEND_BASE_URL", "http://localhost:3000"),
	}

	if AppConfig.FrontendBaseURL == "" || AppConfig.FrontendBaseURL == "https://ww1.your-frontend.com" || AppConfig.FrontendBaseURL == "https://your-frontend.com" {
		log.Fatalf("Error: FRONTEND_BASE_URL environment variable is not set or using a placeholder value. Please set it in your .env file.")
	}

	if AppConfig.SMTPUser == "" {
		log.Fatal("Error: SMTP_USER environment variable is not set or empty.")
	}
	if AppConfig.SMTPPassword == "" {
		log.Fatal("Error: SMTP_PASSWORD environment variable is not set or empty.")
	}
	if AppConfig.SMTPSenderEmail == "" {
		log.Fatal("Error: SMTP_SENDER_EMAIL environment variable is not set or empty.")
	}
	if AppConfig.SMTPHost == "" {
		log.Fatal("Error: SMTP_HOST environment variable is not set or empty.")
	}

	log.Println("Configuration loaded successfully")
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// getEnvAsInt gets an environment variable as integer or returns a default value
func getEnvAsInt(key string, defaultValue int) int {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	intValue := 0
	_, err := fmt.Sscanf(value, "%d", &intValue)
	if err != nil {
		return defaultValue
	}
	return intValue
}
