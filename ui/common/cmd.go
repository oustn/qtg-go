package uicommon

import (
	tea "github.com/charmbracelet/bubbletea"
	"time"
)

func CreateMessage(d time.Duration, msg tea.Msg) tea.Cmd {
	return tea.Tick(d, func(time.Time) tea.Msg {
		return msg
	})
}

func CreateMessageImmediate(msg tea.Msg) tea.Cmd {
	return CreateMessage(0, msg)
}

func Tick() tea.Cmd {
	return CreateMessage(time.Second, TickMsg{})
}

func GotoMsg(name string) tea.Cmd {
	return CreateMessageImmediate(ReplaceRouteMsg{name})
}
