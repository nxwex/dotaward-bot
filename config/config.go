package config

import (
	"errors"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	BotToken  string
	DBPath    string
	AIAPIKey  string
	AIBaseURL string
	AIModel   string
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	token := os.Getenv("BOT_TOKEN")
	if token == "" {
		return nil, errors.New("BOT_TOKEN is not set")
	}

	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "users.db"
		log.Println("DB_PATH is not set, using default: users.db")
	}

	aiKey := os.Getenv("AI_API_KEY")
	if aiKey == "" {
		log.Println("AI_API_KEY is not set")
	}

	baseURL := os.Getenv("AI_BASE_URL")
	if baseURL == "" {
		baseURL = "https://api.openai.com/v1"
	}

	model := os.Getenv("AI_MODEL")
	if model == "" {
		model = "gpt-4o-mini"
	}

	return &Config{
		BotToken:  token,
		DBPath:    dbPath,
		AIAPIKey:  aiKey,
		AIBaseURL: baseURL,
		AIModel:   model,
	}, nil
}
