package uicommon

import "github.com/charmbracelet/lipgloss"

func makeStyle(color lipgloss.TerminalColor) lipgloss.Style {
	return lipgloss.NewStyle().Foreground(color)

}

func makeTextStyle(color lipgloss.TerminalColor) func(...string) string {
	return makeStyle(color).Render
}

var (
	TitleColor         = lipgloss.AdaptiveColor{Dark: "#ffffff", Light: "#24292f"}
	SecondaryTextColor = lipgloss.AdaptiveColor{Dark: "#4d4d4d", Light: "#6e7781"}
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
)
