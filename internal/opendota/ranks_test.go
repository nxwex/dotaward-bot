package opendota

import "testing"

func TestGetRankName(t *testing.T) {
	cases := []struct {
		name     string
		rankTier int
		expected string
	}{
		{"no rank", 0, "Ранг отсутствует"},
		{"unknown rank", 99, "Неизвестно"},
		{"herald 1", 11, "Herald 1"},
		{"legend 5", 55, "Legend 5"},
		{"immortal", 80, "Immortal"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			result := GetRankName(tc.rankTier)
			if result != tc.expected {
				t.Errorf("got %q, want %q", result, tc.expected)
			}
		})
	}
}
