package tuicommon

import (
	tea "github.com/charmbracelet/bubbletea"
	tuimessage "github.com/oustn/qtg/internal/tui/message"
	"time"
)

func setTimeoutMessage(d time.Duration, msg tea.Msg) tea.Cmd {
	return tea.Tick(d, func(time.Time) tea.Msg {
		return msg
	})
}

func WrapperMessage(msg tea.Msg) tea.Cmd {
	return func() tea.Msg {
		return msg
	}
}

func Tick() tea.Cmd {
	return setTimeoutMessage(time.Second, tuimessage.TickMsg{})
}

func Frame() tea.Cmd {
	return setTimeoutMessage(time.Second/60, tuimessage.FrameMsg{})
}

func Created() tea.Cmd {
	return setTimeoutMessage(time.Second, tuimessage.CreatedMsg{})
}

func EnterAuthPageMsg() tea.Cmd {
	return setTimeoutMessage(time.Second, tuimessage.AuthPageEnter{})
}

func PageEnter(index int) tea.Cmd {
	msg := MenuEnterMsg[index]
	if msg == nil {
		return nil
	}
	return setTimeoutMessage(time.Second/10, msg)
}
