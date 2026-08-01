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
