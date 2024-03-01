package uihome

import (
	"fmt"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	uicommon "github.com/oustn/qtg/ui/common"
	"strconv"
	"strings"
)

func InitHome() uicommon.View {
	return &Home{
		ticks:     10,
		autoClose: true,
		choice:    uicommon.Routes[0].Path,
		keys: map[string]key.Binding{
			"Up": key.NewBinding(
				key.WithKeys("up", "k"),
				key.WithHelp("↑/k", "选择上一个"),
			),
			"Down": key.NewBinding(
				key.WithKeys("down", "j"),
				key.WithHelp("↓/j", "选择下一个"),
			),
			"Enter": key.NewBinding(
				key.WithKeys("enter"),
				key.WithHelp("↩︎  ", "进入"),
			),
		},
	}
}

type Home struct {
	choice    string
	autoClose bool
	ticks     int
	keys      map[string]key.Binding
}

func (h *Home) KeyBinds() map[string]key.Binding {
	return h.keys
}

func (h *Home) Render() string {
	tpl := "%s\n\n"
	if h.autoClose {
		tpl += "应用将会在 %s 秒后退出\n\n"
	}

	boxes := make([]interface{}, len(uicommon.Routes))
	formats := make([]string, len(uicommon.Routes))

	for i, route := range uicommon.Routes {
		formats[i] = "%s"
		boxes[i] = uicommon.Checkbox(route.Title, h.choice == route.Path)
	}

	choices := fmt.Sprintf(
		strings.Join(formats, "\n"),
		boxes...,
	)

	if h.autoClose {
		return fmt.Sprintf(tpl, choices, uicommon.SuccessTextStyle(strconv.Itoa(h.ticks)))
	} else {
		return fmt.Sprintf(tpl, choices)
	}
}

func (h *Home) ShouldRenderUserInfo() bool {
	return true
}

func (h *Home) Handle(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case uicommon.ViewPayload:
		h.autoClose = !msg.Created
		return uicommon.Tick()
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, h.keys["Up"]):
			h.selectUp()
		case key.Matches(msg, h.keys["Down"]):
			h.selectDown()
		case key.Matches(msg, h.keys["Enter"]):
			return uicommon.GotoMsg(h.choice)
		}
	case uicommon.TickMsg:
		if h.autoClose {
			if h.ticks == 0 {
				return tea.Quit
			}
			h.countdownTick()
			return uicommon.Tick()
		}
	}
	return nil
}

func (h *Home) selectDown() {
	choice := h.indexOfChoice()
	h.autoClose = false
	choice++
	if choice > len(uicommon.Routes)-1 {
		choice = len(uicommon.Routes) - 1
	}
	h.choice = uicommon.Routes[choice].Path
}

func (h *Home) selectUp() {
	choice := h.indexOfChoice()
	h.autoClose = false
	choice--
	if choice < 0 {
		choice = 0
	}
	h.choice = uicommon.Routes[choice].Path
}

func (h *Home) indexOfChoice() int {
	for i, route := range uicommon.Routes {
		if route.Path == h.choice {
			return i
		}
	}
	return -1
}

func (h *Home) countdownTick() {
	h.ticks--
}
