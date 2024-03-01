package tuimodel

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/oustn/qtg/internal/api"
	"github.com/oustn/qtg/internal/tui/common"
	tuimessage "github.com/oustn/qtg/internal/tui/message"
)

type auth struct {
	Inputs         []textinput.Model
	InputLen       int
	AuthFocusIndex int
	AuthError      error
	Authed         bool
	RefreshToken   string
	QingTingId     string
	Api            *api.QingTingApi
	Error          error
	User           *api.UserInfo
}

func (a *auth) Authorization() tea.Cmd {
	if a.RefreshToken == "" || a.QingTingId == "" {
		return tuicommon.WrapperMessage(tuimessage.AuthTokenInvalid{})
	}
	err := a.Api.Auth(a.RefreshToken, a.QingTingId)
	if err != nil {
		a.Error = err
		a.User = nil
		a.Authed = false
		return tuicommon.WrapperMessage(tuimessage.AuthFailed{})
	}
	a.Error = nil
	a.Authed = true
	a.User = &a.Api.User
	return tuicommon.WrapperMessage(tuimessage.AuthSuccess{})
}

func (a *auth) SelectDown() {
	a.AuthFocusIndex++
	a.adjustFocusIndex()
}

func (a *auth) SelectUp() {
	a.AuthFocusIndex--
	a.adjustFocusIndex()
}

func (a *auth) adjustFocusIndex() {
	inputLen := len(a.Inputs)

	if a.AuthFocusIndex > inputLen {
		a.AuthFocusIndex = 0
	} else if a.AuthFocusIndex < 0 {
		a.AuthFocusIndex = inputLen
	}
}

func (a *auth) AutoFocus() []tea.Cmd {
	inputLen := len(a.Inputs)

	cmds := make([]tea.Cmd, inputLen)
	for i := 0; i <= inputLen-1; i++ {
		if i == a.AuthFocusIndex {
			// Set focused state
			cmds[i] = a.Inputs[i].Focus()
			a.Inputs[i].PromptStyle = tuicommon.FocusedStyle
			a.Inputs[i].TextStyle = tuicommon.FocusedStyle
			continue
		}
		// Remove focused state
		a.Inputs[i].Blur()
		a.Inputs[i].PromptStyle = tuicommon.NoStyle
		a.Inputs[i].TextStyle = tuicommon.NoStyle
	}
	return cmds
}

func (a *auth) UpdateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(a.Inputs))

	for i := range a.Inputs {
		a.Inputs[i], cmds[i] = a.Inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func (a *auth) Reset() []tea.Cmd {
	a.AuthFocusIndex = 0
	a.Inputs[0].SetValue(a.RefreshToken)
	a.Inputs[1].SetValue(a.QingTingId)
	return a.AutoFocus()
}

func initAuth(refreshToken string, qingtingId string) auth {
	qApi := api.InitQingTingApi()

	auth := auth{
		Inputs:         make([]textinput.Model, 2),
		InputLen:       2,
		AuthFocusIndex: 0,
		AuthError:      nil,
		Authed:         false,
		Api:            qApi,
		RefreshToken:   refreshToken,
		QingTingId:     qingtingId,
	}

	var t textinput.Model
	for i := range auth.Inputs {
		t = textinput.New()
		t.Cursor.Style = tuicommon.CursorStyle
		t.CharLimit = 32

		switch i {
		case 0:
			t.Placeholder = "(refresh_token)"
			t.Focus()
			t.PromptStyle = tuicommon.FocusedStyle
			t.TextStyle = tuicommon.FocusedStyle
			t.CharLimit = 32
			if refreshToken != "" {
				t.SetValue(refreshToken)
			}
		case 1:
			t.Placeholder = "(qingting_id)"
			t.CharLimit = 32
			if qingtingId != "" {
				t.SetValue(qingtingId)
			}
		}

		auth.Inputs[i] = t
	}
	return auth
}
