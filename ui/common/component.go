package uicommon

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
)

var (
	Dot       = BlurredTextStyle(" • ")
	CheckMark = lipgloss.NewStyle().SetString("✓").
			Foreground(SuccessColor).
			PaddingRight(1).
			String()
	CrossMark = lipgloss.NewStyle().SetString("✗").
			Foreground(ErrorColor).
			PaddingRight(1).
			String()
)

func Checkbox(label string, checked bool) string {
	if checked {
		return FocusedTextStyle("[x] " + label)
	}
	return fmt.Sprintf("[ ] %s", label)
}

func Button(label string) string {
	return fmt.Sprintf("[ %s ]", BlurredTextStyle(label))
}

func FocusButton(label string) string {
	return FocusedTextStyle("[ " + label + " ]")
}
