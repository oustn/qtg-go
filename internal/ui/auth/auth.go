package teaauth

import (
	"fmt"
	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	teacmd "github.com/oustn/qtg/internal/ui/cmd"
	teacommon "github.com/oustn/qtg/internal/ui/common"
)

var AuthRoute = teacommon.Route{
	Title: "🔑 认证",

	Name: teacommon.AuthRoute,

	RenderComponent: func() teacommon.View {
		a := auth{
			inputs: make([]textinput.Model, 2),
		}

		var t textinput.Model
		for i := range a.inputs {
			t = textinput.New()
			t.Cursor.SetMode(cursor.CursorBlink)
			t.CharLimit = 32
			t.Prompt = ""
			switch i {
			case 0:
				t.Placeholder = "id xxx"
				t.Focus()
			case 1:
				t.Placeholder = "token xxx"
			}
			a.inputs[i] = t
		}
		return &a
	},
}

type keyMap struct {
	next    key.Binding
	prev    key.Binding
	confirm key.Binding
	save    key.Binding
}

var authKeys = keyMap{
	next:    key.NewBinding(key.WithKeys("down", "tab"), key.WithHelp("↓/tab", "下一行")),
	prev:    key.NewBinding(key.WithKeys("up", "shift+tab"), key.WithHelp("↑/shift+tab", "上一行")),
	confirm: key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "认证")),
	save:    key.NewBinding(key.WithKeys("alt+s"), key.WithHelp("alt+s", "保存")),
}

type auth struct {
	focusIndex int
	inputs     []textinput.Model
	authed     bool
	err        error
}

func (a *auth) Init() tea.Cmd {
	return tea.Batch(textinput.Blink, teacmd.UpdateHelpKeys([]*key.Binding{
		&authKeys.next,
		&authKeys.prev,
		&authKeys.confirm,
		&authKeys.save,
	}), teacmd.DisableQQuit())
}

func (a *auth) Render() string {
	var message string

	if a.err != nil {
		message = teacommon.ErrorTextStyle(a.err.Error())
	} else if a.authed {
		message = teacommon.SuccessTextStyle("验证成功，已保存到配置文件")
	} else {
		message = teacommon.BlurredTextStyle("验证")
	}

	return fmt.Sprintf(
		`%s
%s

%s
%s

%s
`,
		teacommon.FocusedStyle.Width(30).Render("蜻蜓ID"),
		a.inputs[0].View(),
		teacommon.FocusedStyle.Width(30).Render("Token"),
		a.inputs[1].View(),
		message,
	) + "\n"
}

func (a *auth) Handler(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, authKeys.next):
			a.nextInput()
		case key.Matches(msg, authKeys.prev):
			a.prevInput()
		case key.Matches(msg, authKeys.confirm):
			cmd := a.validate()
			if cmd == nil {
				return textinput.Blink
			}
			return cmd
		}
		a.flush()
	}
	return a.updateInputs(msg)
}

func (a *auth) SetSize(width, height int) {
}

func (a *auth) flush() {
	for i := range a.inputs {
		a.inputs[i].Blur()
	}
	a.inputs[a.focusIndex].Focus()
}

func (a *auth) nextInput() {
	a.focusIndex = (a.focusIndex + 1) % len(a.inputs)
}

func (a *auth) prevInput() {
	a.focusIndex--
	// Wrap around
	if a.focusIndex < 0 {
		a.focusIndex = len(a.inputs) - 1
	}
}

func (a *auth) updateInputs(msg tea.Msg) tea.Cmd {
	cmd := make([]tea.Cmd, len(a.inputs))

	// Only text inputs with Focus() set will respond, so it's safe to simply
	// update all of them here without any further logic.
	for i := range a.inputs {
		a.inputs[i], cmd[i] = a.inputs[i].Update(msg)
	}

	return tea.Batch(cmd...)
}

func (a *auth) validate() tea.Cmd {
	id := a.inputs[0].Value()
	token := a.inputs[1].Value()

	if len(id) != 32 {
		a.err = fmt.Errorf("请输入32位的蜻蜓ID")
		a.focusIndex = 0
		a.flush()
		return nil
	}
	if len(token) != 32 {
		a.err = fmt.Errorf("请输入32位的Token")
		a.focusIndex = 1
		a.flush()
		return nil
	}
	return teacmd.Auth(token, id, func(err error) {
		if err != nil {
			a.err = err
		} else {
			a.err = nil
			a.authed = true
		}
	})
}
