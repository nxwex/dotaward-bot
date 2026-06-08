package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/nxwex/dotaward-bot/internal/handlers"
	"github.com/nxwex/dotaward-bot/internal/opendota"
	"github.com/nxwex/dotaward-bot/internal/storage"
	"gopkg.in/telebot.v3"
)

func main() {
	log.Printf("starting bot...")

	pref := telebot.Settings{
		Token:  os.Getenv("BOT_TOKEN"),
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	}

	b, err := telebot.NewBot(pref)
	if err != nil {
		log.Fatal("error:", err)
	}

	store := storage.New("users.json")
	dotaClient := opendota.NewClient()
	h := handlers.New(store, dotaClient)

	b.Use(func(next telebot.HandlerFunc) telebot.HandlerFunc {
		return func(c telebot.Context) error {
			if c.Callback() != nil {
				log.Printf("[%d] @%s: %s", c.Sender().ID, c.Sender().Username, c.Callback().Data)
			} else {
				log.Printf("[%d] @%s: %s", c.Sender().ID, c.Sender().Username, c.Text())
			}
			err := next(c)
			if err != nil {
				log.Printf("handler error: %v", err)
			}
			return err
		}
	})

	b.Handle("/start", func(c telebot.Context) error {
		return c.Reply(fmt.Sprintf("Привет, %s!\n\nИспользуй /register <account_id>, чтобы привязать аккаунт.\n\nПосмотреть все доступные команды: /help", c.Sender().FirstName))
	})

	b.Handle("/register", h.Register)
	b.Handle("/profile", h.Profile)
	b.Handle(&telebot.InlineButton{Unique: "profile"}, h.HandleCallBack)
	b.Handle("/lastmatch", h.LastMatch)
	b.Handle(&telebot.InlineButton{Unique: "lastmatch"}, h.HandleCallBack)
	b.Handle("/streak", h.Streak)
	b.Handle("/help", h.Help)

	log.Println("bot is running")
	b.Start()
}
