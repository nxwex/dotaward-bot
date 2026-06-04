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
