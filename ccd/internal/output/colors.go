package output

import (
	"os"
)

var (
	Reset   = "\033[0m"
	Red     = "\033[0;31m"
	Green   = "\033[0;32m"
	Yellow  = "\033[1;33m"
	Blue    = "\033[0;34m"
	Magenta = "\033[0;35m"
	Cyan    = "\033[0;36m"
)

var NoColor = false

func init() {
	if os.Getenv("NO_COLOR") != "" {
		DisableColors()
	}
}

func DisableColors() {
	NoColor = true
	Reset = ""
	Red = ""
	Green = ""
	Yellow = ""
	Blue = ""
	Magenta = ""
	Cyan = ""
}

func Colorize(color, text string) string {
	if NoColor {
		return text
	}
	return color + text + Reset
}
