package uisearch

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	uicommon "github.com/oustn/qtg/ui/common"
)

func InitSearch() uicommon.View {
	return &Search{}
}

type Search struct {
	keys map[string]key.Binding
}

func (s *Search) KeyBinds() map[string]key.Binding {
	return s.keys
}

func (s *Search) Render() string {
	//TODO implement me
	return "Search"
}

func (s *Search) Handle(msg tea.Msg) tea.Cmd {
	//TODO implement me
	return nil
}

func (s *Search) RenderHeader(header string) string {
	return header + "\n\n"
}

func (s *Search) ShouldRenderUserInfo() bool {
	return true
}
