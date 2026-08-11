package opendota

import "fmt"

type Match struct {
	MatchID    int64         `json:"match_id"`
	Duration   int           `json:"duration"`
	RadiantWin bool          `json:"radiant_win"`
	Version    *int          `json:"version"`
	Players    []MatchPlayer `json:"players"`
	Teamfights []Teamfight   `json:"teamfights,omitempty"`
}

type MatchPlayer struct {
	AccountID   *int64         `json:"account_id"`
	PlayerSlot  int            `json:"player_slot"`
	HeroID      int            `json:"hero_id"`
	Kills       int            `json:"kills"`
	Deaths      int            `json:"deaths"`
	Assists     int            `json:"assists"`
	LastHits    int            `json:"last_hits"`
	GoldPerMin  int            `json:"gold_per_min"`
	XpPerMin    int            `json:"xp_per_min"`
	NetWorth    int            `json:"net_worth"`
	Item0       int            `json:"item_0"`
	Item1       int            `json:"item_1"`
	Item2       int            `json:"item_2"`
	Item3       int            `json:"item_3"`
	Item4       int            `json:"item_4"`
	Item5       int            `json:"item_5"`
	PurchaseLog []LogEntry     `json:"purchase_log,omitempty"`
	KillsLog    []LogEntry     `json:"kills_log,omitempty"`
	KilledBy    map[string]int `json:"killed_by,omitempty"`
}

type LogEntry struct {
	Time int    `json:"time"`
	Key  string `json:"key"`
}

type Teamfight struct {
	Start  int `json:"start"`
	End    int `json:"end"`
	Deaths int `json:"deaths"`
}

func (c *Client) GetMatch(matchID int64) (*Match, error) {
	var m Match
	url := fmt.Sprintf("%s/matches/%d", baseURL, matchID)
	if err := c.get(url, &m); err != nil {
		return nil, err
	}
	if m.MatchID == 0 {
		return nil, fmt.Errorf("match not found")
	}
	return &m, nil
}
