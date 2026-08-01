package handlers

import (
	"gopkg.in/telebot.v3"
)

func (h *Handler) Setup(b *telebot.Bot) {
	b.Use(Logger)

	b.Handle("/start", h.Start)
	b.Handle("/register", h.Register)
	b.Handle("/profile", h.Profile)
	b.Handle(&telebot.InlineButton{Unique: "profile"}, h.HandleCallBack)
	b.Handle("/lastmatch", h.LastMatch)
	b.Handle(&telebot.InlineButton{Unique: "lastmatch"}, h.HandleCallBack)
	b.Handle("/streak", h.Streak)
	b.Handle("/maxstreak", h.MaxStreak)
	b.Handle("/help", h.Help)
}
