package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBUser            string `mapstructure:"DB_USER"`
	DBPassword        string `mapstructure:"DB_PASSWORD"`
	DBHost            string `mapstructure:"DB_HOST"`
	DBPort            string `mapstructure:"DB_PORT"`
	DBName            string `mapstructure:"DB_NAME"`
	ServerPort        string `mapstructure:"SERVER_PORT"`
	JWTSecret         string `mapstructure:"JWT_SECRET"`
	SEEDER_ADMIN_PASS string `mapstructure:"SEEDER_ADMIN_PASS"`
	GoogleClientID    string `mapstructure:"GOOGLE_CLIENT_ID"`
	SMTPHost          string `mapstructure:"SMTP_HOST"`
	SMTPPort          int    `mapstructure:"SMTP_PORT"`
	SMTPUser          string `mapstructure:"SMTP_USER"`
	SMTPPassword      string `mapstructure:"SMTP_PASSWORD"`
	SMTPSenderEmail   string `mapstructure:"SMTP_SENDER_EMAIL"`
}

var AppConfig *Config

// LoadConfig loads environment variables from .env file
func LoadConfig() {
	// Load .env file if exists
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	AppConfig = &Config{
		DBUser:            getEnv("DB_USER", "root"),
		DBPassword:        getEnv("DB_PASSWORD", ""),
		DBHost:            getEnv("DB_HOST", "localhost"),
		DBPort:            getEnv("DB_PORT", "3306"),
		DBName:            getEnv("DB_NAME", "terra_db"),
		ServerPort:        getEnv("SERVER_PORT", "8080"),
		JWTSecret:         getEnv("JWT_SECRET", "your-secret-key"),
		SEEDER_ADMIN_PASS: getEnv("SEEDER_ADMIN_PASS", "password123"),
		GoogleClientID:    getEnv("GOOGLE_CLIENT_ID", ""),
		SMTPHost:          getEnv("SMTP_HOST", "smtp.gmail.com"),
		SMTPPort:          getEnvAsInt("SMTP_PORT", 587),
		SMTPUser:          getEnv("SMTP_USER", ""),
		SMTPPassword:      getEnv("SMTP_PASSWORD", ""),
		SMTPSenderEmail:   getEnv("SMTP_SENDER_EMAIL", ""),
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
