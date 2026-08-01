package config

import (
	"errors"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	BotToken string
	DBPath   string
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

	return &Config{
		BotToken: token,
		DBPath:   dbPath,
	}, nil
}
