package opendota

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

const baseURL = "https://api.opendota.com/api"

type Client struct {
	http *http.Client
}

func NewClient() *Client {
	return &Client{
		http: &http.Client{Timeout: 10 * time.Second},
	}
}

type Profile struct {
	AccountID   int64  `json:"account_id"`
	Personaname string `json:"personaname"`
	Avatar      string `json:"avatarfull"`
}

type WinLose struct {
	Win  int `json:"win"`
	Lose int `json:"lose"`
}

type PlayerProfile struct {
	Profile  `json:"profile"`
	RankTier int `json:"rank_tier"`
	WinLose
}

type RecentMatch struct {
	MatchID  int64 `json:"match_id"`
	Kills    int   `json:"kills"`
	Deaths   int   `json:"deaths"`
	Assists  int   `json:"assists"`
	HeroID   int   `json:"hero_id"`
	Win      int   `json:"win"`
	Duration int   `json:"duration"`
	LastHits int   `json:"last_hits"`
	GPM      int   `json:"gold_per_min"`
	XPM      int   `json:"xp_per_min"`
}

func (c *Client) GetRecentMatch(accountID int64) (*RecentMatch, error) {
	url := fmt.Sprintf("%s/players/%d/recentMatches", baseURL, accountID)

	resp, err := c.http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var matches []RecentMatch
	if err := json.NewDecoder(resp.Body).Decode(&matches); err != nil {
		return nil, err
	}

	if len(matches) == 0 {
		return nil, fmt.Errorf("матчи не найдены")
	}

	return &matches[0], nil
}

func (c *Client) GetProfile(accountID int64) (*PlayerProfile, error) {
	resp, err := c.http.Get(fmt.Sprintf("%s/players/%d", baseURL, accountID))
	if err != nil {
		log.Printf("error resp GetProfile: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	var profile PlayerProfile
	if err := json.NewDecoder(resp.Body).Decode(&profile); err != nil {
		return nil, err
	}

	resp2, err := c.http.Get(fmt.Sprintf("%s/players/%d/wl", baseURL, accountID))
	if err != nil {
		log.Printf("error resp2 GetProfile: %v", err)
		return nil, err
	}
	defer resp2.Body.Close()

	if err := json.NewDecoder(resp2.Body).Decode(&profile.WinLose); err != nil {
		return nil, err
	}

	return &profile, nil
}
