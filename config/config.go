package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	MongoURI       string
	DatabaseName   string
	ServerPort     string
	AllowedOrigins string
	JWTSecret      string
}

func LoadConfig() *Config {
	// Load .env hanya lokal. Di production (GitHub Actions) ENV langsung dari OS.
	_ = godotenv.Load()

	config := &Config{
		MongoURI:       mustEnv("MONGO_URI"),
		DatabaseName:   getEnv("DATABASE_NAME", "dashboard_db"),
		ServerPort:     getEnv("SERVER_PORT", "8080"),
		AllowedOrigins: getEnv("ALLOWED_ORIGINS", "*"),
		JWTSecret:      mustEnv("JWT_SECRET"),
	}

	return config
}

// wajib ada, kalau tidak ada ENV -> panic biar ketahuan
func mustEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("Missing required environment variable: %s", key)
	}
	return value
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
