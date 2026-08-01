package models

type User struct {
	DotaID       int64  `json:"dota_id"`
	TelegramID   int64  `json:"telegram_id"`
	TelegramName string `json:"telegram_name"`
}
