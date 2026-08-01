package opendota

import "testing"

func TestCalcStreak(t *testing.T) {
	cases := []struct {
		name    string
		matches []RecentMatch
		streak  int
		win     bool
	}{
		{
			name:    "win streak",
			matches: []RecentMatch{{RadiantWin: true, PlayerSlot: 0}},
			streak:  1,
			win:     true,
		},
		{
			name:    "loss on radiant",
			matches: []RecentMatch{{RadiantWin: false, PlayerSlot: 0}},
			streak:  1,
			win:     false,
		},
		{
			name:    "win on dire",
			matches: []RecentMatch{{RadiantWin: false, PlayerSlot: 128}},
			streak:  1,
			win:     true,
		},
		{
			name:    "loss on dire",
			matches: []RecentMatch{{RadiantWin: true, PlayerSlot: 128}},
			streak:  1,
			win:     false,
		},
		{
			name: "streak resets after loss",
			matches: []RecentMatch{
				{RadiantWin: false, PlayerSlot: 0},
				{RadiantWin: true, PlayerSlot: 0},
				{RadiantWin: true, PlayerSlot: 0},
			},
			streak: 1,
			win:    false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			streak, win := CalcStreak(tc.matches)
			if streak != tc.streak || win != tc.win {
				t.Errorf("got %d %v, want %d %v", streak, win, tc.streak, tc.win)
			}
		})
	}
}

func TestCalcMaxStreak(t *testing.T) {
	cases := []struct {
		name    string
		matches []RecentMatch
		maxWin  int
		maxLoss int
	}{
		{"empty", []RecentMatch{}, 0, 0},
		{
			"all wins",
			[]RecentMatch{
				{RadiantWin: true, PlayerSlot: 0},
				{RadiantWin: true, PlayerSlot: 0},
				{RadiantWin: true, PlayerSlot: 0},
			},
			3, 0,
		},
		{
			"all losses",
			[]RecentMatch{
				{RadiantWin: false, PlayerSlot: 0},
				{RadiantWin: false, PlayerSlot: 0},
			},
			0, 2,
		},
		{
			"mixed",
			[]RecentMatch{
				{RadiantWin: true, PlayerSlot: 0},
				{RadiantWin: true, PlayerSlot: 0},
				{RadiantWin: false, PlayerSlot: 0},
				{RadiantWin: false, PlayerSlot: 0},
				{RadiantWin: false, PlayerSlot: 0},
				{RadiantWin: true, PlayerSlot: 0},
			},
			2, 3,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			maxWin, maxLoss := CalcMaxStreak(tc.matches)
			if maxWin != tc.maxWin || maxLoss != tc.maxLoss {
				t.Errorf("got %d %d, want %d %d", maxWin, maxLoss, tc.maxWin, tc.maxLoss)
			}
		})
	}
}
