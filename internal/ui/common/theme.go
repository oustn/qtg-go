package teacommon

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/lucasb-eyer/go-colorful"
)

func makeStyle(color lipgloss.TerminalColor) lipgloss.Style {
	return lipgloss.NewStyle().Foreground(color)

}

func makeTextStyle(color lipgloss.TerminalColor) func(...string) string {
	return makeStyle(color).Render
}

var (
	White              = lipgloss.Color("#ffffff")
	TitleColor         = lipgloss.AdaptiveColor{Dark: "#ffffff", Light: "#24292f"}
	SecondaryTextColor = lipgloss.AdaptiveColor{Dark: "#4d4d4d", Light: "#57606a"}
	ErrorColor         = lipgloss.AdaptiveColor{Dark: "#F78166", Light: "#cf222e"}
	WarnColor          = lipgloss.AdaptiveColor{Dark: "#E3B341", Light: "#9a6700"}
	SuccessColor       = lipgloss.AdaptiveColor{Dark: "#56D364", Light: "#1a7f37"}
	PrimaryColor       = lipgloss.AdaptiveColor{Dark: "#6CA4F8", Light: "#0969da"}
	AccentColor        = lipgloss.AdaptiveColor{Dark: "#DB61A2", Light: "#8250df"}

	TitleTextStyle   = makeTextStyle(TitleColor)         // title text style
	BlurredStyle     = makeStyle(SecondaryTextColor)     // blurred style
	BlurredTextStyle = makeTextStyle(SecondaryTextColor) // help text style
	FocusedStyle     = makeStyle(AccentColor)            // focused style
	FocusedTextStyle = makeTextStyle(AccentColor)        // focused text style
	SuccessStyle     = makeStyle(SuccessColor)           // error style
	SuccessTextStyle = makeTextStyle(SuccessColor)       // success text style
	ErrorTextStyle   = makeTextStyle(ErrorColor)         // error text style
	NoStyle          = lipgloss.NewStyle()               // no style
	CursorStyle      = FocusedStyle.Copy()               // cursor style

	InfoStyle     = makeStyle(PrimaryColor)                     // primary style
	InfoTextStyle = makeStyle(PrimaryColor).Italic(true).Render // primary style

	StatusBarSelectedFileForegroundColor = lipgloss.AdaptiveColor{Dark: "#ffffff", Light: "#ffffff"}
	StatusBarSelectedFileBackgroundColor = lipgloss.AdaptiveColor{Dark: "#F25D94", Light: "#F25D94"}
	StatusBarBarForegroundColor          = lipgloss.AdaptiveColor{Dark: "#ffffff", Light: "#ffffff"}
	StatusBarBarBackgroundColor          = lipgloss.AdaptiveColor{Dark: "#3c3836", Light: "#3c3836"}
	StatusBarTotalFilesForegroundColor   = lipgloss.AdaptiveColor{Dark: "#ffffff", Light: "#ffffff"}
	StatusBarTotalFilesBackgroundColor   = lipgloss.AdaptiveColor{Dark: "#A550DF", Light: "#A550DF"}
	StatusBarLogoForegroundColor         = lipgloss.AdaptiveColor{Dark: "#ffffff", Light: "#ffffff"}
	StatusBarLogoBackgroundColor         = lipgloss.AdaptiveColor{Dark: "#6124DF", Light: "#6124DF"}
)

func ColorGrid(xSteps, ySteps int) [][]string {
	x0y0, _ := colorful.Hex("#F25D94")
	x1y0, _ := colorful.Hex("#EDFF82")
	x0y1, _ := colorful.Hex("#643AFF")
	x1y1, _ := colorful.Hex("#14F9D5")

	x0 := make([]colorful.Color, ySteps)
	for i := range x0 {
		x0[i] = x0y0.BlendLuv(x0y1, float64(i)/float64(ySteps))
	}

	x1 := make([]colorful.Color, ySteps)
	for i := range x1 {
		x1[i] = x1y0.BlendLuv(x1y1, float64(i)/float64(ySteps))
	}

	grid := make([][]string, ySteps)
	for x := 0; x < ySteps; x++ {
		y0 := x0[x]
		grid[x] = make([]string, xSteps)
		for y := 0; y < xSteps; y++ {
			grid[x][y] = y0.BlendLuv(x1[x], float64(y)/float64(xSteps)).Hex()
		}
	}

	return grid
}
