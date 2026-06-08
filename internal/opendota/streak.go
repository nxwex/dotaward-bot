package opendota

func CalcStreak(matches []RecentMatch) (int, bool) {
	if len(matches) == 0 {
		return 0, false
	}

	first := matches[0]
	isRadiant := first.PlayerSlot < 128
	win := (first.RadiantWin && isRadiant) || (!first.RadiantWin && !isRadiant)

	streak := 1
	for _, m := range matches {
		mIsRadiant := m.PlayerSlot < 128
		mWin := (m.RadiantWin && mIsRadiant) || (!m.RadiantWin && !mIsRadiant)
		if mWin != win {
			break
		}
		streak++
	}

	return streak, win
}
