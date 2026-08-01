package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nxwex/dotaward-bot/config"
	"github.com/nxwex/dotaward-bot/internal/handlers"
	"github.com/nxwex/dotaward-bot/internal/opendota"
	"github.com/nxwex/dotaward-bot/internal/repository"
	"gopkg.in/telebot.v3"
)

func main() {
	log.Printf("starting bot...")

	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	pref := telebot.Settings{
		Token:  cfg.BotToken,
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	}

	b, err := telebot.NewBot(pref)
	if err != nil {
		log.Fatal("error:", err)
	}

	db, err := repository.New(cfg.DBPath)
	if err != nil {
		log.Fatalf("database error: %v", err)
	}
	defer db.Close()

	dotaClient := opendota.NewClient()

	h := handlers.New(db, dotaClient)
	h.Setup(b) // миддлвар + хендлы

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	log.Println("bot is running")
	go b.Start()

	<-quit
	log.Println("shutting down...")
	b.Stop()
}
