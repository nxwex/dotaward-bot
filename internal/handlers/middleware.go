package handlers

import (
	"log"

	"gopkg.in/telebot.v3"
)

func Logger(next telebot.HandlerFunc) telebot.HandlerFunc {
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
}
