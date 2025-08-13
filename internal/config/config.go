// in file: /internal/config/config.go
package config

import "os"

// Config holds all configuration for the application.
type Config struct {
	TelegramToken       string
	CloudinaryCloudName string
	CloudinaryAPIKey    string
	CloudinaryAPISecret string
}

// New loads configuration from environment variables.
func New() *Config {
	return &Config{
		TelegramToken:       os.Getenv("TELEGRAM_TOKEN"),
		CloudinaryCloudName: os.Getenv("CLOUDINARY_CLOUD_NAME"),
		CloudinaryAPIKey:    os.Getenv("CLOUDINARY_API_KEY"),
		CloudinaryAPISecret: os.Getenv("CLOUDINARY_API_SECRET"),
	}
}
