package config

import (
	"os"

	"github.com/joho/godotenv"
)


type DatabaseConfig struct {
	User     string
	Password string
	Name     string
}

type Config struct {
	DB       DatabaseConfig
	BotToken string
}

func ConfigNew() *Config {
	_ = godotenv.Load()
	return &Config{
		DB: DatabaseConfig{
			User:     getEnv("DB_USER", ""),
			Password: getEnv("DB_PASS", ""),
			Name:     getEnv("DB_NAME", ""),
		},
		BotToken: getEnv("BOT_TOKEN", ""),
	}
}

func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
}
