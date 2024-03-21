package teahome

import (
	"fmt"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	teaauth "github.com/oustn/qtg/internal/ui/auth"
	teacmd "github.com/oustn/qtg/internal/ui/cmd"
	teacommon "github.com/oustn/qtg/internal/ui/common"
	teasearch "github.com/oustn/qtg/internal/ui/search"
	"strings"
)

func Checkbox(label string, checked bool) string {
	if checked {
		return teacommon.FocusedTextStyle("[x] " + label)
	}
	return fmt.Sprintf("[ ] %s", label)
}

type homeKeys struct {
	Up    key.Binding
	Down  key.Binding
	Enter key.Binding
}

var HomeRoute = teacommon.Route{
	Title: "🏠 主页",
	RenderComponent: func() teacommon.View {
		return &home{
			choice: routes[0].Name,
			keys: homeKeys{
				Up: key.NewBinding(
					key.WithKeys("up", "k"),
					key.WithHelp("↑/k", "选择上一个"),
				),
				Down: key.NewBinding(
					key.WithKeys("down", "j"),
					key.WithHelp("↓/j", "选择下一个"),
				),
				Enter: key.NewBinding(
					key.WithKeys("enter"),
					key.WithHelp("↩︎  ", "进入"),
				),
			},
		}
	},
}

var routes = []teacommon.Route{
	teaauth.AuthRoute,
	teasearch.SearchRoute,
}

type home struct {
	choice string
	ticks  int
	keys   homeKeys
}

func (h *home) Init() tea.Cmd {
	return teacmd.UpdateHelpKeys([]*key.Binding{
		&h.keys.Up,
		&h.keys.Down,
		&h.keys.Enter,
	})
}

func (h *home) Render() string {
	tpl := "%s\n\n"
	boxes := make([]interface{}, len(routes))
	formats := make([]string, len(routes))

	for i, route := range routes {
		formats[i] = "%s"
		boxes[i] = Checkbox(route.Title, h.choice == route.Name)
	}
	choices := fmt.Sprintf(
		strings.Join(formats, "\n"),
		boxes...,
	)
	return fmt.Sprintf(tpl, choices)
}

func (h *home) Handler(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case teacmd.ViewCreatedMsg:
		return teacmd.UpdateStatusBar("主页", "")
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, h.keys.Up):
			h.selectUp()
		case key.Matches(msg, h.keys.Down):
			h.selectDown()
		case key.Matches(msg, h.keys.Enter):
			return teacmd.Redirect(h.choice)
		}
	}
	return nil
}

func (h *home) selectDown() {
	choice := h.indexOfChoice()
	choice++
	if choice > len(routes)-1 {
		choice = len(routes) - 1
	}
	h.choice = routes[choice].Name
}

func (h *home) selectUp() {
	choice := h.indexOfChoice()
	choice--
	if choice < 0 {
		choice = 0
	}
	h.choice = routes[choice].Name
}

func (h *home) indexOfChoice() int {
	for i, route := range routes {
		if route.Name == h.choice {
			return i
		}
	}
	return -1
}

func (h *home) SetSize(width, height int) {
}
