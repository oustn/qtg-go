package teacommon

import "github.com/charmbracelet/bubbles/key"

var Keys = []key.Binding{
	key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "帮助"),
	),
	key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "退出"),
	),
	key.NewBinding(
		key.WithKeys("up"),
		key.WithHelp("up", "上一个"),
	),
	key.NewBinding(
		key.WithKeys("down"),
		key.WithHelp("down", "下一个"),
	),
	key.NewBinding(
		key.WithKeys("left"),
		key.WithHelp("left", "左边"),
	),
}
