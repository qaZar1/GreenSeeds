package ws

func calcPercent(current, total int) int {
	if total == 0 {
		return 0
	}
	return int(float64(current) / float64(total) * 100)
}
