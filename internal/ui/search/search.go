package teasearch

import (
	"fmt"
	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	teacmd "github.com/oustn/qtg/internal/ui/cmd"
	teacommon "github.com/oustn/qtg/internal/ui/common"
	teacomponents "github.com/oustn/qtg/internal/ui/components"
	"math"
)

type title struct {
	input   textinput.Model
	content string
	focused bool
}

func (t *title) Render() string {
	if !t.focused {
		return t.content
	}
	return t.input.View()
}

func (t *title) Handler(msg tea.Msg) tea.Cmd {
	return nil
}

func (t *title) SetSize(width, height int) {
}

func (t *title) SetTitle(title string) {
	t.content = title
}

func (t *title) focus() {
	t.focused = true
	t.input.Focus()
	t.input.SetValue("")
}

func (t *title) blur() {
	t.focused = false
	t.input.Blur()
}

var types = []string{"综合排序", "热度最高", "最近更新"}
var typeKeys = []string{"0", "1", "2"}

type searchKey struct {
	next   key.Binding
	prev   key.Binding
	accept key.Binding
}

type searchItem struct {
	id          string
	title       string
	description string
}

func (s searchItem) Title() string {
	return s.title
}

func (s searchItem) Description() string {
	return s.description
}

func (s searchItem) FilterValue() string {
	return s.title
}

var SearchRoute = teacommon.Route{
	Title: "🔍 搜索",
	Name:  teacommon.SearchRoute,
	RenderComponent: func() teacommon.View {
		items := make([]teacomponents.ListItem, 0)

		searchTitle := title{
			input: textinput.New(),
		}

		searchTitle.input.Placeholder = "搜索..."
		searchTitle.input.Width = 30
		searchTitle.input.Prompt = "🔍 "
		searchTitle.input.Cursor.Style = teacommon.FocusedStyle

		l := teacomponents.NewList(items)
		l.Title = &searchTitle
		l.EmptyTip = fmt.Sprintf(`%s %s %s`,
			teacommon.BlurredTextStyle("听点有趣的~ 键入"),
			teacommon.FocusedStyle.Bold(true).Render("/"),
			teacommon.BlurredTextStyle("开始搜索"),
		)

		return &search{
			list:  l,
			title: &searchTitle,
			keys: searchKey{
				next: key.NewBinding(
					key.WithKeys("ctrl+l"),
					key.WithHelp("ctrl+l", "搜索下一页"),
				),
				prev: key.NewBinding(
					key.WithKeys("ctrl+h"),
					key.WithHelp("ctrl+h", "搜索上一页"),
				),
				accept: key.NewBinding(
					key.WithKeys("ctrl+d"),
					key.WithHelp("ctrl+d", "开始下载"),
				),
			},
		}
	},
}

const pageSize = 30

type search struct {
	list       *teacomponents.List
	title      *title
	focused    bool
	keys       searchKey
	keyword    string
	page       int
	searchType string
	total      int
}

func (s *search) Init() tea.Cmd {
	s.list.UnbindFilterKeys()
	return tea.Batch(
		teacmd.UpdateHelpKeys([]*key.Binding{
			&s.list.ListKeyMap.Filter,
			&s.list.ListKeyMap.ClearFilter,
			&s.list.ListKeyMap.CancelWhileFiltering,
			&s.list.ListKeyMap.AcceptWhileFiltering,
			&s.list.ListKeyMap.CursorUp,
			&s.list.ListKeyMap.CursorDown,
			&s.list.ListKeyMap.PrevPage,
			&s.list.ListKeyMap.NextPage,
			&s.list.ListKeyMap.GoToStart,
			&s.list.ListKeyMap.GoToEnd,
			&s.keys.prev,
			&s.keys.next,
			&s.keys.accept,
		}),
	)
}

func (s *search) Render() string {
	return s.list.Render()
}

func (s *search) Handler(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, s.list.ListKeyMap.Filter):
			return s.focusSearchInput()
		case key.Matches(msg, s.list.ListKeyMap.CancelWhileFiltering):
			return s.blurSearchInput()
		case key.Matches(msg, s.list.ListKeyMap.AcceptWhileFiltering):
			return s.handleAccept()
		case key.Matches(msg, s.keys.next):
			return s.handleSearchNext()
		case key.Matches(msg, s.keys.prev):
			return s.handleSearchPrev()
		case key.Matches(msg, s.keys.accept):
			return s.handleDownload()
		}
	case teacmd.SearchResultMsg:
		return s.HandleSearchResult(msg.Result, msg.Error)

	}
	if s.focused {
		return s.updateInput(msg)
	}
	return s.list.Handler(msg)
}

func (s *search) SetSize(width, height int) {
	s.list.SetSize(width, height)
}

func (s *search) focusSearchInput() tea.Cmd {
	s.title.focus()
	s.list.BindFilterKeys()
	if s.focused {
		s.title.input.SetValue("")
	}
	s.focused = true
	return tea.Batch(cursor.Blink, teacmd.DisableBothQuit())
}

func (s *search) updateInput(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	s.title.input, cmd = s.title.input.Update(msg)
	return cmd
}

func (s *search) blurSearchInput() tea.Cmd {
	if !s.focused {
		return nil
	}
	s.focused = false
	s.title.blur()
	//s.title.input.SetValue("")
	s.list.UnbindFilterKeys()
	return teacmd.EnableBothQuit()
}

func (s *search) handleAccept() tea.Cmd {
	if s.focused {
		return s.handleSearch()
	}
	return s.handleSelect()
}

func (s *search) handleSearch() tea.Cmd {
	keyword := s.title.input.Value()
	if keyword == "" {
		return nil
	}
	s.keyword = keyword
	s.searchType = typeKeys[0]
	s.page = 1
	cmd := s.blurSearchInput()
	s.list.Loading(true)
	return tea.Batch(cmd, s.search())
}

func (s *search) search() tea.Cmd {
	s.list.Loading(true)
	return teacmd.HandleSearch(s.keyword, s.searchType, s.page)
}

func (s *search) handleSelect() tea.Cmd {
	return nil
}

func (s *search) HandleSearchResult(result teacommon.SearchResult, err error) tea.Cmd {
	var items []teacomponents.ListItem

	s.list.Loading(false)

	if err == nil {
		for _, channel := range result.Data {
			items = append(items, channel)
		}
	}

	s.list.SetEntity(items)
	s.total = result.Total
	s.page = result.Page

	s.title.SetTitle(teacommon.BlurredTextStyle(fmt.Sprintf(`"%s"的搜索结果，共%d(%d)条 （%d/%.0f页）`,
		result.Keyword,
		result.Total,
		result.PageSize,
		result.Page,
		math.Ceil(float64(result.Total)/float64(result.PageSize)),
	)))
	return nil
}

func (s *search) handleSearchNext() tea.Cmd {
	if s.focused {
		return nil
	}
	totalPage := int(math.Ceil(float64(s.total) / float64(pageSize)))
	if s.page >= totalPage {
		return nil
	}
	s.page += 1
	return s.search()
}

func (s *search) handleSearchPrev() tea.Cmd {
	if s.focused || s.page == 1 {
		return nil
	}
	s.page -= 1
	s.list.Loading(true)
	return s.search()
}

func (s *search) handleDownload() tea.Cmd {
	item := s.list.GetSelect()
	return teacmd.StartDownload(item.(teacommon.Channel))
}
