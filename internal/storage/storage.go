package storage

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"strings"
	"sync"
)

var ErrAlreadyRegistered = errors.New("user already registered")

type User struct {
	AccountID int64  `json:"account_id"`
	Username  string `json:"username"`
}

type Storage struct {
	mu    sync.RWMutex
	users map[int64]User
	path  string
}

func New(path string) *Storage {
	s := &Storage{
		users: make(map[int64]User),
		path:  path,
	}
	s.load()
	return s
}

func (s *Storage) Register(telegramID, accountID int64, username string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.users[telegramID]; exists {
		return ErrAlreadyRegistered
	}

	s.users[telegramID] = User{AccountID: accountID, Username: username}
	return s.save()
}

func (s *Storage) GetAccountID(telegramID int64) (int64, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	user, ok := s.users[telegramID]
	return user.AccountID, ok
}

func (s *Storage) GetAccountIDByUsername(username string) (int64, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	username = strings.TrimPrefix(username, "@")
	for _, user := range s.users {
		if strings.EqualFold(user.Username, username) {
			return user.AccountID, true
		}
	}
	return 0, false
}

func (s *Storage) save() error {
	data, err := json.Marshal(s.users)
	if err != nil {
		log.Printf("ошибка маршалинга в storage: %v", err)
		return err
	}
	return os.WriteFile(s.path, data, 0644)
}

func (s *Storage) load() {
	data, err := os.ReadFile(s.path)
	if err != nil {
		return
	}
	if err := json.Unmarshal(data, &s.users); err != nil {
		log.Printf("ошибка загрузки storage: %v", err)
	}
	log.Printf("user storage loaded: %d users.", len(s.users))
}
