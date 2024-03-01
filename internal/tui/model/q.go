package tuimodel

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/muesli/reflow/indent"
	"github.com/oustn/qtg/internal/config"
	tuicommon "github.com/oustn/qtg/internal/tui/common"
	tuiupdate "github.com/oustn/qtg/internal/tui/update"
	tuiview "github.com/oustn/qtg/internal/tui/view"
)

type Q struct {
	auth    *auth
	menu    *menu
	cfg     *config.Config
	cfgPath string
}

func InitQ(cfg *config.Config, configPath string) Q {
	auth := initAuth(cfg.Settings.RefreshToken, cfg.Settings.QingTingId)

	menu := initMenu()

	return Q{
		auth:    &auth,
		menu:    &menu,
		cfg:     cfg,
		cfgPath: configPath,
	}
}

func (q Q) Init() tea.Cmd {
	return tea.Batch(tuicommon.Created(), tea.EnterAltScreen, q.Auth())
}

func (q Q) View() string {
	s := tuiview.HeaderView(q)

	if q.menu.Home {
		s += tuiview.HomeView(q)
	} else {
		switch q.menu.Choice {
		case tuicommon.AuthPage:
			s += tuiview.AuthView(q)
			break
		case tuicommon.SearchPage:
			s += tuiview.SearchView(q)
		default:
			s += "Hello world"
		}
	}
	return indent.String("\n"+s+"\n\n", 2)
}

func (q Q) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// global message
	if msg, ok := msg.(tea.KeyMsg); ok {
		k := msg.String()
		if k == "q" || k == "esc" || k == "ctrl+c" {
			if q.menu.Home {
				return q, tea.Quit
			}
			q.menu.reset()
			return q, tuicommon.Tick()
		}
	}
	if q.menu.Home {
		return tuiupdate.HomeUpdate(msg, q)
	}
	switch q.menu.Choice {
	case tuicommon.AuthPage:
		return tuiupdate.AuthUpdate(msg, q)

	}
	return q, nil
}

func (q Q) Info() tuicommon.QInfo {
	return tuicommon.QInfo{
		Home:           q.menu.Home,
		Choice:         q.menu.Choice,
		ChoiceTitle:    q.menu.GetChoiceTitle(),
		AuthInputs:     q.auth.Inputs,
		AuthFocusIndex: q.auth.AuthFocusIndex,
	}
}

func (q Q) GetMenuInfo() (choice string, ticks int, autoClose bool) {
	return q.menu.Choice, q.menu.Ticks, q.menu.AutoClose
}

func (q Q) MenuSelectUp() {
	q.menu.SelectUp()
}

func (q Q) MenuSelectDown() {
	q.menu.SelectDown()
}

func (q Q) MenuSelectEnter() tea.Cmd {
	return q.menu.SelectEnter()
}

func (q Q) CountdownTick() {
	q.menu.Countdown()
}

func (q Q) Auth() tea.Cmd {
	return q.auth.Authorization()
}

func (q Q) GetAuthInfo() (bool, bool, string) {
	nickName := ""
	if q.auth.User != nil {
		nickName = q.auth.User.Nickname
	}
	return q.auth.RefreshToken != "" && q.auth.QingTingId != "", q.auth.Authed, nickName
}

func (q Q) AuthSelectDown() {
	q.auth.SelectDown()
}

func (q Q) AuthSelectUp() {
	q.auth.SelectUp()
}

func (q Q) AuthAutoFocus() []tea.Cmd {
	return q.auth.AutoFocus()
}

func (q Q) AuthUpdateInputs(msg tea.Msg) tea.Cmd {
	return q.auth.UpdateInputs(msg)
}

func (q Q) PageEnter() []tea.Cmd {
	if q.menu.Home {
		return nil
	}
	switch q.menu.Choice {
	case tuicommon.AuthPage:
		return q.auth.Reset()
	}
	return nil
}
