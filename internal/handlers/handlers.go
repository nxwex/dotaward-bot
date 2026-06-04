package handlers

import (
	"fmt"
	"log"
	"strconv"

	"github.com/nxwex/dotaward-bot/internal/opendota"
	"gopkg.in/telebot.v3"
)

type UserStorage interface {
	Register(telegramID, accountID int64, username string)
	GetAccountID(telegramID int64) (int64, bool)
	GetAccountIDByUsername(username string) (int64, bool)
}

type DotaClient interface {
	GetRecentMatch(accountID int64) (*opendota.RecentMatch, error)
}

type Handler struct {
	storage    UserStorage
	dotaClient DotaClient
}

func New(storage UserStorage, dotaClient DotaClient) *Handler {
	return &Handler{storage: storage, dotaClient: dotaClient}
}

func (h *Handler) Register(c telebot.Context) error {
	args := c.Args()
	if len(args) == 0 {
		return c.Reply("Использование: /register <dota_account_id>")
	}

	accountID, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return c.Reply("account_id должен быть числом")
	}

	h.storage.Register(c.Sender().ID, accountID, c.Sender().Username)

	log.Printf("[bot] register: @%s telegram: %d, dota: %d", c.Sender().Username, c.Sender().ID, accountID)
	return c.Send(fmt.Sprintf("Аккаунт %d привязан!", accountID))
}

func (h *Handler) LastMatch(c telebot.Context) error {
	var accountID int64
	var ok bool
	var username string

	switch {
	case c.Message().ReplyTo != nil:
		accountID, ok = h.storage.GetAccountID(c.Message().ReplyTo.Sender.ID)
		username = c.Message().ReplyTo.Sender.FirstName
	case len(c.Args()) > 0:
		accountID, ok = h.storage.GetAccountIDByUsername(c.Args()[0])
		username = c.Args()[0]
	default:
		accountID, ok = h.storage.GetAccountID(c.Sender().ID)
		username = c.Sender().FirstName
	}

	if !ok {
		return c.Reply("Пользователь не зарегистрирован")
	}

	match, err := h.dotaClient.GetRecentMatch(accountID)
	if err != nil {
		return c.Reply("Не удалось получить матч")
	}

	result := "✅ Победа"
	if match.Win == 0 {
		result = "❌ Поражение"
	}

	hero := opendota.GetHeroName(match.HeroID)
	minutes := match.Duration / 60
	seconds := match.Duration % 60

	msg := fmt.Sprintf(
		"📊 Последний матч • %s\n"+
			"─────────────────\n"+
			"%s\n\n"+
			"🧙 Герой: %s\n"+
			"⚔️ KDA: %d/%d/%d\n"+
			"💰 GPM: %d\n"+
			"⭐ XPM: %d\n"+
			"🎯 Last Hits: %d\n"+
			"⏱️ Длительность: %02d:%02d\n"+
			"🏅 Номер матча: %d",
		username, result, hero,
		match.Kills, match.Deaths, match.Assists,
		match.GPM, match.XPM,
		match.LastHits,
		minutes, seconds,
		match.MatchID)

	log.Printf("[bot] lastmatch: %s (account: %d, match: %d) -> @%s", username, accountID, match.MatchID, c.Sender().Username)
	return c.Reply(msg)
}

func (h *Handler) Help(c telebot.Context) error {
	msg := "📖 Команды бота:\n\n" +
		"/register <dota_account_id> - привязать свой аккаунт\n" +
		"  └ ID можно найти на opendota.com после входа через Steam\n\n" +
		"/lastmatch - узнать статистику последнего матча\n" +
		"  └ без аргументов - твой матч\n" +
		"  └ @username - матч конкретного пользователя\n" +
		"  └ ответом на сообщение - матч того на кого ответил\n\n" +
		"/help - это сообщение"
	return c.Reply(msg)
}
