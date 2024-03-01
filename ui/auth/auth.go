package uiauth

import (
	"fmt"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/oustn/qtg/internal/api"
	tuicommon "github.com/oustn/qtg/internal/tui/common"
	uicommon "github.com/oustn/qtg/ui/common"
	"strings"
	"time"
)

var (
	blurredButton = uicommon.Button("认证")
	focusedButton = uicommon.FocusButton("认证")
)

type Auth struct {
	keys       map[string]key.Binding
	inputs     []textinput.Model
	focusIndex int
	message    string
	err        error
	validated  bool
}

func (a *Auth) KeyBinds() map[string]key.Binding {
	return a.keys
}

func InitAuth() uicommon.View {
	model := Auth{
		inputs: make([]textinput.Model, 2),
	}

	var t textinput.Model
	for i := range model.inputs {
		t = textinput.New()
		t.Cursor.Style = uicommon.CursorStyle
		t.CharLimit = 32

		switch i {
		case 0:
			t.Placeholder = "(refresh_token)"
			t.Focus()
			t.PromptStyle = tuicommon.FocusedStyle
			t.CharLimit = 32
		case 1:
			t.Placeholder = "(qingting_id)"
			t.CharLimit = 32
		}

		model.inputs[i] = t
	}

	return &model
}

func (a *Auth) Render() string {
	var b strings.Builder

	b.WriteString(uicommon.Dot + tuicommon.SubtitleStyle("抓取 https://user.qtfm.cn/u2/api/v4/auth 接口获取"))
	b.WriteString("\n\n")
	for i := range a.inputs {
		if i == 0 {
			b.WriteString(uicommon.InfoTextStyle("刷新令牌"))
		} else {
			b.WriteString(uicommon.InfoTextStyle("蜻蜓ID"))
		}
		b.WriteString("\n")
		b.WriteString(a.inputs[i].View())
		if i < len(a.inputs)-1 {
			b.WriteRune('\n')
			b.WriteRune('\n')
		}
	}

	button := &blurredButton
	if a.focusIndex == len(a.inputs) {
		button = &focusedButton
	}
	_, _ = fmt.Fprintf(&b, "\n\n%s\n\n", *button)

	if a.message != "" || a.err != nil {
		if a.err != nil {
			b.WriteString(uicommon.CrossMark + uicommon.ErrorTextStyle(""+a.err.Error()))
		} else {
			b.WriteString(uicommon.CheckMark + uicommon.SuccessTextStyle(""+a.message))
		}
	}

	return b.String()
}

func (a *Auth) ShouldRenderUserInfo() bool {
	return false
}

func (a *Auth) Handle(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case uicommon.ViewPayload:
		token := msg.RefreshToken
		id := msg.QingTingId
		a.initValues(token, id)
		return tea.Batch(a.autoFocus(a.focusIndex)...)
	case tea.KeyMsg:
		switch msg.String() {
		// Set focus to next input
		case "tab", "shift+tab", "enter", "up", "down":
			s := msg.String()
			inputLen := len(a.inputs)
			if s == "enter" && a.focusIndex == inputLen {
				index, err := a.validate()
				if err != nil {
					a.err = err
					a.message = err.Error()

					return tea.Batch(a.autoFocus(index)...)
				}
				user, err := a.auth()
				if err != nil {
					a.err = err
					a.message = err.Error()
					return tea.Batch(a.autoFocus(0)...)
				}
				a.message = "认证成功，欢迎 " + user.Nickname + " 👏"
				return a.AuthSuccess()
			}

			// Cycle indexes
			if s == "up" || s == "shift+tab" {
				a.selectUp()
			} else {
				a.selectDown()
			}

			cmds := a.autoFocus(a.focusIndex)
			return tea.Batch(cmds...)
		}
	}

	cmd := a.updateInputs(msg)

	return cmd
}

func (a *Auth) selectUp() {
	a.focusIndex--
	a.adjustFocusIndex()
}

func (a *Auth) selectDown() {
	a.focusIndex++
	a.adjustFocusIndex()
}

func (a *Auth) adjustFocusIndex() {
	inputLen := len(a.inputs)

	if a.focusIndex > inputLen {
		a.focusIndex = 0
	} else if a.focusIndex < 0 {
		a.focusIndex = inputLen
	}
}

func (a *Auth) autoFocus(index int) []tea.Cmd {
	inputLen := len(a.inputs)

	cmds := make([]tea.Cmd, inputLen)
	for i := 0; i <= inputLen-1; i++ {
		if i == index {
			// Set focused state
			cmds[i] = a.inputs[i].Focus()
			a.inputs[i].PromptStyle = uicommon.CursorStyle
			continue
		}
		// Remove focused state
		a.inputs[i].Blur()
		a.inputs[i].PromptStyle = uicommon.NoStyle
		a.inputs[i].TextStyle = uicommon.NoStyle
	}
	return cmds
}

func (a *Auth) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(a.inputs))

	for i := range a.inputs {
		a.inputs[i], cmds[i] = a.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func (a *Auth) initValues(strings ...string) {
	for i, str := range strings {
		a.inputs[i].SetValue(str)
	}
}

func (a *Auth) validate() (int, error) {
	if a.inputs[0].Value() == "" {
		return 0, fmt.Errorf("刷新令牌不能为空")
	}
	if a.inputs[1].Value() == "" {
		return 1, fmt.Errorf("蜻蜓ID不能为空")
	}
	return -1, nil
}

func (a *Auth) auth() (*uicommon.UserInfo, error) {
	qtApi := api.InitQingTingApi(a.inputs[0].Value(), a.inputs[1].Value())
	err := qtApi.Auth()
	if err != nil {
		return nil, err
	}
	return &qtApi.User, nil
}

func (a *Auth) AuthSuccess() tea.Cmd {
	token := a.inputs[0].Value()
	id := a.inputs[1].Value()
	return uicommon.CreateMessage(time.Second, uicommon.AuthSuccessMsg{Token: token, Id: id})
}
