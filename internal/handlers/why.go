package handlers

import (
	"context"
	"encoding/json"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/nxwex/dotaward-bot/internal/models"
	"github.com/nxwex/dotaward-bot/internal/opendota"
	"gopkg.in/telebot.v3"
)

func (h *Handler) Why(c telebot.Context) error {
	args := c.Args()
	if len(args) == 0 {
		return c.Reply("Использование: /why <match_id>")
	}

	matchID, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return c.Reply("match_id должен быть числом")
	}

	_ = c.Notify(telebot.Typing)

	match, err := h.dotaClient.GetMatch(matchID)
	if err != nil {
		log.Printf("GetMatch error: %v", err)
		return c.Reply("Не удалось получить матч")
	}

	myDotaID, registered := h.storage.GetDotaID(c.Sender().ID)

	type playerView struct {
		Hero        string         `json:"hero"`
		AccountID   *int64         `json:"account_id,omitempty"`
		IsRadiant   bool           `json:"is_radiant"`
		Kills       int            `json:"kills"`
		Deaths      int            `json:"deaths"`
		Assists     int            `json:"assists"`
		LastHits    int            `json:"last_hits"`
		GPM         int            `json:"gpm"`
		XPM         int            `json:"xpm"`
		NetWorth    int            `json:"net_worth"`
		Items       []int          `json:"items"`
		PurchaseLog any            `json:"purchase_log,omitempty"`
		KillsLog    any            `json:"kills_log,omitempty"`
		KilledBy    map[string]int `json:"killed_by,omitempty"`
		IsMe        bool           `json:"is_me,omitempty"`
	}

	players := make([]playerView, 0, len(match.Players))
	for _, p := range match.Players {
		pv := playerView{
			Hero:      opendota.GetHeroName(p.HeroID),
			AccountID: p.AccountID,
			IsRadiant: p.PlayerSlot < 128,
			Kills:     p.Kills,
			Deaths:    p.Deaths,
			Assists:   p.Assists,
			LastHits:  p.LastHits,
			GPM:       p.GoldPerMin,
			XPM:       p.XpPerMin,
			NetWorth:  p.NetWorth,
			Items:     []int{p.Item0, p.Item1, p.Item2, p.Item3, p.Item4, p.Item5},
		}
		if len(p.PurchaseLog) > 0 {
			pv.PurchaseLog = p.PurchaseLog
		}
		if len(p.KillsLog) > 0 {
			pv.KillsLog = p.KillsLog
		}
		if len(p.KilledBy) > 0 {
			pv.KilledBy = p.KilledBy
		}
		if registered && p.AccountID != nil && *p.AccountID == myDotaID {
			pv.IsMe = true
		}
		players = append(players, pv)
	}

	payload := map[string]any{
		"match_id":    match.MatchID,
		"duration":    match.Duration,
		"radiant_win": match.RadiantWin,
		"parsed":      match.Version != nil,
		"players":     players,
		"teamfights":  match.Teamfights,
	}
	if match.Version == nil {
		payload["warning"] = "Матч не распарсен — нет логов смертей и покупок"
	}

	matchJSON, _ := json.Marshal(payload)

	ctx, cancel := context.WithTimeout(context.Background(), 75*time.Second)
	defer cancel()

	text, err := h.ai.Analyze(ctx, string(matchJSON), "", "")
	if err != nil {
		log.Printf("AI error: %v", err)
		return c.Reply("Ошибка анализа. Попробуйте позже.")
	}
	if len(text) > 4000 {
		text = text[:4000] + "\n..."
	}

	sent, err := c.Bot().Send(c.Chat(), text, &telebot.SendOptions{
		ParseMode: telebot.ModeHTML,
	})
	if err != nil {
		return err
	}

	log.Printf(
		"Analysis message sent: chatID=%d messageID=%d matchID=%d",
		sent.Chat.ID,
		sent.ID,
		matchID,
	)

	err = h.contexts.SaveMatchContext(&models.MatchContext{
		ChatID:          sent.Chat.ID,
		MessageID:       sent.ID,
		ParentMessageID: 0,
		MatchID:         matchID,
		Question:        "",
		Answer:          text,
		MatchData:       string(matchJSON),
		CreatedAt:       time.Now(),
		ExpiresAt:       time.Now().Add(14 * 24 * time.Hour),
	})
	if err != nil {
		log.Printf("SaveMatchContext error: %v", err)
	}

	return nil
}

func (h *Handler) HandleText(c telebot.Context) error {
	msg := c.Message()

	if msg.ReplyTo == nil || msg.ReplyTo.Sender == nil || !msg.ReplyTo.Sender.IsBot {
		return nil
	}

	history, err := h.contexts.GetMatchHistory(msg.Chat.ID, msg.ReplyTo.ID)
	if err != nil {
		log.Printf(
			"GetMatchHistory error: chatID=%d messageID=%d err=%v",
			msg.Chat.ID,
			msg.ReplyTo.ID,
			err,
		)
		return c.Reply("Ошибка получения контекста")
	}

	if len(history) == 0 {
		log.Printf(
			"GetMatchHistory not found: chatID=%d messageID=%d",
			msg.Chat.ID,
			msg.ReplyTo.ID,
		)
		return nil
	}

	current := history[len(history)-1]

	log.Printf(
		"GetMatchHistory ok: chatID=%d messageID=%d matchID=%d history=%d",
		msg.Chat.ID,
		msg.ReplyTo.ID,
		current.MatchID,
		len(history),
	)

	_ = c.Notify(telebot.Typing)

	var conversation strings.Builder

	for _, item := range history {
		if item.Question != "" {
			conversation.WriteString("Пользователь: ")
			conversation.WriteString(item.Question)
			conversation.WriteString("\n")
		}

		conversation.WriteString("Коуч: ")
		conversation.WriteString(item.Answer)
		conversation.WriteString("\n\n")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 75*time.Second)
	defer cancel()

	text, err := h.ai.Analyze(
		ctx,
		current.MatchData,
		conversation.String(),
		msg.Text,
	)
	if err != nil {
		log.Printf("AI follow-up error: %v", err)
		return c.Reply("Ошибка ответа")
	}

	if len(text) > 4000 {
		text = text[:3997] + "..."
	}

	sent, err := c.Bot().Send(c.Chat(), text, &telebot.SendOptions{
		ParseMode: telebot.ModeHTML,
	})
	if err != nil {
		return err
	}

	err = h.contexts.SaveMatchContext(&models.MatchContext{
		ChatID:          sent.Chat.ID,
		MessageID:       sent.ID,
		ParentMessageID: msg.ReplyTo.ID,
		MatchID:         current.MatchID,
		Question:        msg.Text,
		Answer:          text,
		MatchData:       current.MatchData,
		CreatedAt:       time.Now(),
		ExpiresAt:       time.Now().Add(14 * 24 * time.Hour),
	})
	if err != nil {
		log.Printf("SaveMatchContext error: %v", err)
	}

	return nil
}
