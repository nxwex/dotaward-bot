package models

import "time"

type MatchContext struct {
	ID              int64
	ChatID          int64
	MessageID       int
	ParentMessageID int
	MatchID         int64
	Question        string
	Answer          string
	MatchData       string
	CreatedAt       time.Time
	ExpiresAt       time.Time
}
