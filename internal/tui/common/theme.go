package tuicommon

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
)

var term = termenv.EnvColorProfile()

func ColorFg(val, color string) string {
	return termenv.String(val).Foreground(term.Color(color)).String()
}

func makeFgStyle(color string) func(string) string {
	return termenv.Style{}.Foreground(term.Color(color)).Styled
}

var (
	SubtitleStyle = makeFgStyle("240")
	BlurredStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render
	FocusedStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	CursorStyle   = FocusedStyle.Copy()
	NoStyle       = lipgloss.NewStyle()
)
