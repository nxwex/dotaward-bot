package handlers

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/nxwex/dotaward-bot/internal/opendota"
	"github.com/nxwex/dotaward-bot/internal/storage"
	"gopkg.in/telebot.v3"
)

type UserStorage interface {
	Register(telegramID, accountID int64, username string) error
	GetAccountID(telegramID int64) (int64, bool)
	GetAccountIDByUsername(username string) (int64, bool)
}

type DotaClient interface {
	GetRecentMatch(accountID int64) (*opendota.RecentMatch, error)
	GetRecentMatches(accountID int64, limit int) ([]opendota.RecentMatch, error)
	GetProfile(accountID int64) (*opendota.PlayerProfile, error)
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

	if err = h.storage.Register(c.Sender().ID, accountID, c.Sender().Username); err != nil {
		if errors.Is(err, storage.ErrAlreadyRegistered) {
			return c.Reply("Вы уже зарегистрированы!")
		} else {
			log.Printf("error register user: %v", err)
			return c.Reply("Неизвестная ошибка")
		}
	}

	log.Printf("[bot] register: @%s telegram: %d, dota: %d", c.Sender().Username, c.Sender().ID, accountID)
	return c.Send(fmt.Sprintf("Аккаунт %d привязан!", accountID))
}

func (h *Handler) Profile(c telebot.Context) error {
	accountID, ok := h.storage.GetAccountID(c.Sender().ID)
	if !ok {
		return c.Reply("Ты не зарегистрирован. Используй /register <dota_account_id>")
	}
	return h.sendProfile(c, accountID)
}

func (h *Handler) sendProfile(c telebot.Context, accountID int64) error {
	profile, err := h.dotaClient.GetProfile(accountID)
	if err != nil {
		log.Printf("GetProfile error: %v", err)
		return c.Reply("Не удалось получить профиль")
	}

	total := profile.Win + profile.Lose
	winRate := 0.0
	if total > 0 {
		winRate = float64(profile.Win) / float64(total) * 100
	}

	rank := opendota.GetRankName(profile.RankTier)

	msg := fmt.Sprintf(
		"• Игрок: %s\n"+
			"─────────────────\n"+
			"• Ранг: %s\n"+
			"• Матчей: %d\n"+
			"• Побед: %d\n"+
			"• Поражений: %d\n"+
			"• Винрейт: %.1f%%",
		profile.Personaname, rank, total,
		profile.Win, profile.Lose, winRate)

	markup := &telebot.ReplyMarkup{}
	btnLastMatch := markup.Data("Последний матч", "lastmatch", fmt.Sprintf("lastmatch:%d", accountID))
	btnDotabuff := markup.URL("Dotabuff", fmt.Sprintf("https://ru.dotabuff.com/players/%d", accountID))
	btnOpenDota := markup.URL("OpenDota", fmt.Sprintf("https://opendota.com/players/%d", accountID))
	markup.Inline(
		markup.Row(btnLastMatch),
		markup.Row(btnDotabuff, btnOpenDota),
	)

	profilePhoto := &telebot.Photo{
		File:    telebot.FromURL(profile.Avatar),
		Caption: msg,
	}
	return c.Send(profilePhoto, markup)
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

	log.Printf("[bot] lastmatch: %s (account: %d) -> @%s", username, accountID, c.Sender().Username)
	return h.sendLastMatch(c, accountID)
}

func (h *Handler) sendLastMatch(c telebot.Context, accountID int64) error {
	profile, err := h.dotaClient.GetProfile(accountID)
	if err != nil {
		log.Printf("GetProfile error: %v", err)
		return c.Reply("Ошибка получения профиля")
	}

	match, err := h.dotaClient.GetRecentMatch(accountID)
	if err != nil {
		log.Printf("GetRecentMatch error: %v", err)
		return c.Reply(fmt.Sprintf("Не удалось получить матч игрока %s", profile.Personaname))
	}

	isRadiant := match.PlayerSlot < 128
	win := (match.RadiantWin && isRadiant) || (!match.RadiantWin && !isRadiant)

	result := "✅ Победа"
	if !win {
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
		profile.Personaname, result, hero,
		match.Kills, match.Deaths, match.Assists,
		match.GPM, match.XPM,
		match.LastHits,
		minutes, seconds,
		match.MatchID)

	markup := &telebot.ReplyMarkup{}
	btnProfile := markup.Data("Профиль игрока", "profile", fmt.Sprintf("profile:%d", accountID))
	markup.Inline(markup.Row(btnProfile))

	return c.Send(msg, markup)
}

func (h *Handler) Streak(c telebot.Context) error {
	accountID, ok := h.storage.GetAccountID(c.Sender().ID)
	if !ok {
		return c.Reply("Вы не зарегистированы. Используйте /register <dota_account_id>")
	}

	matches, err := h.dotaClient.GetRecentMatches(accountID, 20)
	if err != nil {
		log.Printf("GetRecentMatches error: %v", err)
		return c.Reply("Не удалсоь получить матчи")
	}

	if len(matches) == 0 {
		return c.Reply("Матчи не найдены")
	}

	streak, win := opendota.CalcStreak(matches)

	statusEmoji := "🔥"
	status := "победа"
	if streak > 1 {
		status = "победы"
	}
	if streak >= 5 {
		status = "побед"
	}

	if !win {
		statusEmoji = "☠️"
		status = "поражение"
		if streak > 1 {
			status = "поражения"
		}
		if streak >= 5 {
			status = "поражений"
		}
	}

	return c.Reply(fmt.Sprintf("%s Серия: %d %s подряд", statusEmoji, streak, status))
}

func (h *Handler) Help(c telebot.Context) error {
	msg := "📖 Команды бота:\n\n" +
		"/register <dota_account_id> - привязать свой аккаунт\n" +
		"  └ ID можно найти на opendota.com после входа через Steam\n\n" +
		"/profile - статистика вашего аккаунта\n\n" +
		"/lastmatch - узнать статистику последнего матча\n" +
		"  └ без аргументов - твой матч\n" +
		"  └ @username - матч конкретного пользователя\n" +
		"  └ ответом на сообщение - матч того на кого ответил\n\n" +
		"/streak - серия побед/поражений\n\n" +
		"/help - это сообщение"
	return c.Reply(msg)
}

func (h *Handler) HandleCallBack(c telebot.Context) error {
	data := c.Callback().Data

	parts := strings.Split(data, ":")
	if len(parts) != 2 {
		return c.Respond()
	}

	action, value := parts[0], parts[1]

	switch action {
	case "lastmatch":
		accountID, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return c.Respond()
		}

		_ = c.Respond()
		return h.sendLastMatch(c, accountID)

	case "profile":
		accountID, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return c.Respond()
		}
		_ = c.Respond()
		_ = c.Delete()
		return h.sendProfile(c, accountID)

	default:
		return c.Respond()
	}
}
