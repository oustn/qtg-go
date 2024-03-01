package tuicommon

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type Menu interface {
}

type QInfo struct {
	Home           bool
	Choice         string
	ChoiceTitle    string
	AuthInputs     []textinput.Model
	AuthFocusIndex int
}

type Q interface {
	Init() tea.Cmd
	Update(tea.Msg) (tea.Model, tea.Cmd)
	View() string
	Info() QInfo

	GetMenuInfo() (choice string, ticks int, autoClose bool)
	MenuSelectUp()
	MenuSelectDown()
	MenuSelectEnter() tea.Cmd
	CountdownTick()

	GetAuthInfo() (bool, bool, string)
	AuthSelectDown()
	AuthSelectUp()
	AuthAutoFocus() []tea.Cmd
	AuthUpdateInputs(msg tea.Msg) tea.Cmd

	PageEnter() []tea.Cmd
}
