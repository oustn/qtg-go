package teacomponents

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	teacommon "github.com/oustn/qtg/internal/ui/common"
	"strings"
)

type ListItem interface {
	list.Item
	Title() string
	Description() string
}

var keyMap = list.KeyMap{
	// Browsing.
	CursorUp: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "上一项"),
	),
	CursorDown: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "下一项"),
	),
	PrevPage: key.NewBinding(
		key.WithKeys("left", "h", "pgup", "b", "u"),
		key.WithHelp("←/h/pgup", "上一页"),
	),
	NextPage: key.NewBinding(
		key.WithKeys("right", "l", "pgdown", "f", "d"),
		key.WithHelp("→/l/pgdn", "下一页"),
	),
	GoToStart: key.NewBinding(
		key.WithKeys("home", "g"),
		key.WithHelp("g/home", "跳转到开头"),
	),
	GoToEnd: key.NewBinding(
		key.WithKeys("end", "G"),
		key.WithHelp("G/end", "跳转到结尾"),
	),
	Filter: key.NewBinding(
		key.WithKeys("/"),
		key.WithHelp("/", "输入关键字"),
	),
	ClearFilter: key.NewBinding(
		key.WithKeys("/"),
		key.WithHelp("/", "清空关键字"),
	),

	// Filtering.
	CancelWhileFiltering: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "退出搜索"),
	),

	AcceptWhileFiltering: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "搜索"),
	),
}

func NewList(items []ListItem) *List {
	var listItems []list.Item
	//for _, item := range items {
	//	listItems = append(listItems, item)
	//}
	l := List{
		list: list.New(listItems, list.NewDefaultDelegate(), 0, 0),
	}
	l.list.Styles.Title = lipgloss.NewStyle()
	l.list.SetShowHelp(false)
	l.list.DisableQuitKeybindings()
	l.list.SetFilteringEnabled(true)
	l.list.SetShowStatusBar(false)
	l.list.KeyMap = keyMap
	l.ListKeyMap = keyMap
	return &l
}

type List struct {
	list       list.Model
	Items      *[]ListItem
	Title      teacommon.Component
	EmptyTip   string
	ListKeyMap list.KeyMap
}

func (l *List) Render() string {
	if l.Title != nil {
		l.list.Title = l.Title.Render()
	}
	return strings.Replace(l.list.View(), "No items.", l.EmptyTip, 1)
}

func (l *List) Handler(msg tea.Msg) tea.Cmd {
	var commands []tea.Cmd
	switch msg.(type) {
	case tea.KeyMsg:
		// Don't match any of the keys below if we're actively filtering.
		if l.list.FilterState() == list.Filtering {
			break
		}
	}

	newListModel, cmd := l.list.Update(msg)
	l.list = newListModel
	commands = append(commands, cmd)
	return tea.Batch(commands...)
}

func (l *List) SetSize(width, height int) {
	l.list.SetSize(width, height)
}

func (l *List) BindFilterKeys() {
	l.ListKeyMap.Filter.SetEnabled(false)
	l.ListKeyMap.ClearFilter.SetEnabled(true)
	l.ListKeyMap.CancelWhileFiltering.SetEnabled(true)
}

func (l *List) UnbindFilterKeys() {
	l.ListKeyMap.Filter.SetEnabled(true)
	l.ListKeyMap.ClearFilter.SetEnabled(false)
	l.ListKeyMap.CancelWhileFiltering.SetEnabled(false)
}

func (l *List) Loading(b bool) {
	if b {
		l.list.StartSpinner()
	} else {
		l.list.StopSpinner()
	}
}

func (l *List) SetEntity(items []ListItem) {
	var listItems []list.Item
	for _, item := range items {
		listItems = append(listItems, item)
	}
	l.list.SetItems(listItems)
}

func (l *List) GetSelect() ListItem {
	item := l.list.SelectedItem()

	return item.(ListItem)
}
