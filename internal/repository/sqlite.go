package repository

import (
	"database/sql"
	"errors"
	"log"
	"os"
	"strings"
	"time"

	sqlite3 "github.com/mattn/go-sqlite3"
	"github.com/nxwex/dotaward-bot/internal/models"
)

var ErrAlreadyRegistered = errors.New("user already registered")

type DB struct {
	db *sql.DB
}

func New(path string) (*DB, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		log.Printf("database file not found, creating new: %s", path)
	} else {
		log.Printf("database file found: %s", path)
	}

	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		db.Close()
		return nil, err
	}

	db.SetMaxOpenConns(1)

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS users (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	telegram_id INTEGER NOT NULL UNIQUE,
	dota_id INTEGER NOT NULL,
	telegram_name TEXT NOT NULL COLLATE NOCASE,
	created_at DATETIME NOT NULL
);
	CREATE INDEX IF NOT EXISTS idx_users_telegram_name ON users(telegram_name);`)
	if err != nil {
		db.Close()
		return nil, err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS match_contexts (
	id                INTEGER PRIMARY KEY AUTOINCREMENT,
	chat_id           INTEGER NOT NULL,
	message_id        INTEGER NOT NULL,
	parent_message_id INTEGER NOT NULL DEFAULT 0,
	match_id          INTEGER NOT NULL,
	question          TEXT NOT NULL,
	answer            TEXT NOT NULL,
	match_data        TEXT NOT NULL,
	created_at        DATETIME NOT NULL,
	expires_at        DATETIME NOT NULL,
	UNIQUE(chat_id, message_id)
)`)
	if err != nil {
		db.Close()
		return nil, err
	}

	return &DB{db: db}, nil
}

func (d *DB) Close() error {
	return d.db.Close()
}

func (d *DB) Save(u *models.User) (int64, error) {
	result, err := d.db.Exec(
		`INSERT INTO users (dota_id, telegram_id, telegram_name, created_at) VALUES (?, ?, ?, ?)`,
		u.DotaID,
		u.TelegramID,
		u.TelegramName,
		time.Now(),
	)
	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return 0, ErrAlreadyRegistered
		}
		return 0, err
	}

	return result.LastInsertId()
}

func (d *DB) GetDotaID(telegramID int64) (int64, bool) {
	var dotaID int64
	err := d.db.QueryRow(`SELECT dota_id FROM users WHERE telegram_id = ?`, telegramID).Scan(&dotaID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, false
		}
		log.Printf("GetDotaID query error: %v", err)
		return 0, false
	}
	return dotaID, true
}

func (d *DB) GetDotaIDByUsername(username string) (int64, bool) {
	username = strings.TrimPrefix(username, "@")

	var dotaID int64
	err := d.db.QueryRow(`SELECT dota_id FROM users WHERE telegram_name = ? LIMIT 1`, username).Scan(&dotaID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, false
		}
		log.Printf("GetDotaIDByUsername query error: %v", err)
		return 0, false
	}
	return dotaID, true
}

func (d *DB) SaveMatchContext(a *models.MatchContext) error {
	_, err := d.db.Exec(
		`INSERT INTO match_contexts (chat_id, message_id, parent_message_id, match_id, question, answer, match_data, created_at, expires_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		a.ChatID, a.MessageID, a.ParentMessageID, a.MatchID, a.Question, a.Answer, a.MatchData, a.CreatedAt, a.ExpiresAt,
	)
	return err
}

func (d *DB) GetMatchContext(chatID int64, messageID int) (*models.MatchContext, error) {
	a := &models.MatchContext{}

	err := d.db.QueryRow(
		`SELECT id, chat_id, message_id, parent_message_id, match_id, question, answer, match_data, created_at, expires_at
		 FROM match_contexts
		 WHERE chat_id = ? AND message_id = ? AND expires_at > ?`,
		chatID, messageID, time.Now(),
	).Scan(
		&a.ID, &a.ChatID, &a.MessageID, &a.ParentMessageID, &a.MatchID,
		&a.Question, &a.Answer, &a.MatchData, &a.CreatedAt, &a.ExpiresAt,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return a, nil
}

func (d *DB) GetMatchHistory(chatID int64, messageID int) ([]*models.MatchContext, error) {
	var history []*models.MatchContext

	currentID := messageID

	for currentID != 0 {
		a, err := d.GetMatchContext(chatID, currentID)
		if err != nil {
			return nil, err
		}
		if a == nil {
			break
		}

		history = append(history, a)
		currentID = a.ParentMessageID
	}

	for i, j := 0, len(history)-1; i < j; i, j = i+1, j-1 {
		history[i], history[j] = history[j], history[i]
	}

	return history, nil
}
