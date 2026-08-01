package opendota

import "fmt"

var ranks = map[int]string{
	1: "Herald",
	2: "Guardian",
	3: "Crusader",
	4: "Archon",
	5: "Legend",
	6: "Ancient",
	7: "Divine",
	8: "Immortal",
}

func GetRankName(rankTier int) string {

	if rankTier == 0 {
		return "Ранг отсутствует"
	}

	rank := rankTier / 10
	stars := rankTier % 10

	name, ok := ranks[rank]
	if !ok {
		return "Неизвестно"
	}

	if rank == 8 {
		return name
	}

	return fmt.Sprintf("%s %d", name, stars)
}
