package cmd

import "fmt"

const (
	colorRed    = "\033[31m"
	colorYellow = "\033[33m"
	colorGreen  = "\033[32m"
	colorReset  = "\033[0m"
)

func getColorValue(value int) string {
	switch {
	case value > 80:
		return fmt.Sprintf("%s%d%s", colorRed, value, colorReset)
	case value > 50:
		return fmt.Sprintf("%s%d%s", colorYellow, value, colorReset)
	default:
		return fmt.Sprintf("%s%d%s", colorGreen, value, colorReset)
	}
}
