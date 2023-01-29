package cli

type TermColor string

const (
	ColorReset TermColor = "\033[0m"
	ColorBold  TermColor = "\033[1m"

	ColorBlack       TermColor = "\033[30m"
	ColorRed         TermColor = "\033[31m"
	ColorGreen       TermColor = "\033[32m"
	ColorYellow      TermColor = "\033[33m"
	ColorBlue        TermColor = "\033[34m"
	ColorPurple      TermColor = "\033[35m"
	ColorCyan        TermColor = "\033[36m"
	ColorLightgray   TermColor = "\033[37m"
	ColorDarkgray    TermColor = "\033[90m"
	ColorLightred    TermColor = "\033[91m"
	ColorLightgreen  TermColor = "\033[92m"
	ColorLightyellow TermColor = "\033[93m"
	ColorLightblue   TermColor = "\033[94m"
	ColorLightpurple TermColor = "\033[95m"
	ColorLightcyan   TermColor = "\033[96m"
	ColorWhite       TermColor = "\033[97m"
)

func (c TermColor) StringColored(s string) string {
	return string(c) + s + string(ColorReset)
}
