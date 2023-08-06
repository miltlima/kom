package cmd

func getEmoji(cpuUsage, memoryUsage int) string {
	switch {
	case cpuUsage > 80 || memoryUsage > 80:
		return "🟥"
	case cpuUsage > 50 || memoryUsage > 50:
		return "🟨"
	default:
		return "🟩"
	}
}
