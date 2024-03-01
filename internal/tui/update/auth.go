package tuiupdate

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/oustn/qtg/internal/tui/common"
	tuimessage "github.com/oustn/qtg/internal/tui/message"
)

func AuthUpdate(msg tea.Msg, q tuicommon.Q) (tea.Model, tea.Cmd) {
	info := q.Info()

	switch msg := msg.(type) {
	case tuimessage.AuthPageEnter:
		return q, tea.Batch(q.PageEnter()...)
	case tea.KeyMsg:
		switch msg.String() {

		// Set focus to next input
		case "tab", "shift+tab", "enter", "up", "down":
			s := msg.String()
			inputLen := len(info.AuthInputs)
			if s == "enter" && info.AuthFocusIndex == inputLen {
				return q, nil
			}

			// Cycle indexes
			if s == "up" || s == "shift+tab" {
				q.AuthSelectUp()
			} else {
				q.AuthSelectDown()
			}

			cmds := q.AuthAutoFocus()
			return q, tea.Batch(cmds...)
		}
	}

	cmd := q.AuthUpdateInputs(msg)

	return q, cmd
}
