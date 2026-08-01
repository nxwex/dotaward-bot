package opendota

import (
	"encoding/json"
	"fmt"
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
	MatchID    int64 `json:"match_id"`
	PlayerSlot int   `json:"player_slot"`
	RadiantWin bool  `json:"radiant_win"`
	Kills      int   `json:"kills"`
	Deaths     int   `json:"deaths"`
	Assists    int   `json:"assists"`
	HeroID     int   `json:"hero_id"`
	Duration   int   `json:"duration"`
	LastHits   int   `json:"last_hits"`
	GPM        int   `json:"gold_per_min"`
	XPM        int   `json:"xp_per_min"`
}

func (c *Client) GetRecentMatch(accountID int64) (*RecentMatch, error) {
	url := fmt.Sprintf("%s/players/%d/recentMatches", baseURL, accountID)

	var matches []RecentMatch
	if err := c.get(url, &matches); err != nil {
		return nil, err
	}

	if len(matches) == 0 {
		return nil, fmt.Errorf("матчи не найдены")
	}

	return &matches[0], nil
}

func (c *Client) GetRecentMatches(accountID int64, limit int) ([]RecentMatch, error) {
	url := fmt.Sprintf("%s/players/%d/matches", baseURL, accountID)
	if limit > 0 {
		url += fmt.Sprintf("?limit=%d", limit)
	}

	var matches []RecentMatch
	if err := c.get(url, &matches); err != nil {
		return nil, err
	}

	return matches, nil
}

func (c *Client) GetProfile(accountID int64) (*PlayerProfile, error) {
	var profile PlayerProfile

	if err := c.get(fmt.Sprintf("%s/players/%d", baseURL, accountID), &profile); err != nil {
		return nil, err
	}

	if err := c.get(fmt.Sprintf("%s/players/%d/wl", baseURL, accountID), &profile.WinLose); err != nil {
		return nil, err
	}

	return &profile, nil
}

func (c *Client) get(url string, v any) error {
	resp, err := c.http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("opendota api error: %d", resp.StatusCode)
	}

	return json.NewDecoder(resp.Body).Decode(v)
}
