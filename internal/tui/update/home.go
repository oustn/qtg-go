package tuiupdate

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/oustn/qtg/internal/tui/common"
	tuimessage "github.com/oustn/qtg/internal/tui/message"
)

func HomeUpdate(msg tea.Msg, q tuicommon.Q) (tea.Model, tea.Cmd) {
	_, tick, autoClose := q.GetMenuInfo()

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "j", "down":
			q.MenuSelectDown()
		case "k", "up":
			q.MenuSelectUp()
		case "enter":
			return q, q.MenuSelectEnter()
		}

	case tuimessage.CreatedMsg, tuimessage.TickMsg:
		if autoClose {
			if tick == 0 {
				return q, tea.Quit
			}
			q.CountdownTick()
			return q, tuicommon.Tick()
		}
	}

	return q, nil
}
