package teacommon

import tea "github.com/charmbracelet/bubbletea"

type (
	View interface {
		Component
		Init() tea.Cmd
	}

	Component interface {
		Render() string
		Handler(msg tea.Msg) tea.Cmd
		SetSize(width, height int)
	}

	RenderComponent func() View
)

type Route struct {
	Title           string
	Name            string
	RenderComponent RenderComponent
}

const (
	HomeRoute   = "home"
	SearchRoute = "search"
	AuthRoute   = "auth"
)

type Callback func(err error)
