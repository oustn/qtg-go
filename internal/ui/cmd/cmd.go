package teacmd

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	teacommon "github.com/oustn/qtg/internal/ui/common"
	"time"
)

func CreateMsg(msg tea.Msg) tea.Cmd {
	return func() tea.Msg {
		return msg
	}
}

func UpdateStatusBar(title, msg string) tea.Cmd {
	return CreateMsg(UpdateStatusBarMsg{Title: title, Msg: msg})
}

func Auth(token, id string, callbacks ...teacommon.Callback) tea.Cmd {
	callback := func(_ error) {}
	if len(callbacks) > 0 {
		callback = callbacks[0]
	}
	return CreateMsg(AuthMsg{Id: id, Token: token, Callback: callback})
}

func DisableQQuit() tea.Cmd {
	return CreateMsg(UpdateQuitKeyMsg{Q: false, Esc: true})
}
func DisableEscQuit() tea.Cmd {
	return CreateMsg(UpdateQuitKeyMsg{Q: true, Esc: false})
}
func DisableBothQuit() tea.Cmd {
	return CreateMsg(UpdateQuitKeyMsg{})
}

func EnableQQuit() tea.Cmd {
	return CreateMsg(UpdateQuitKeyMsg{Q: true})
}

func EnableEscQuit() tea.Cmd {
	return CreateMsg(UpdateQuitKeyMsg{Esc: true})
}

func EnableBothQuit() tea.Cmd {
	return CreateMsg(UpdateQuitKeyMsg{Q: true, Esc: true})
}

func Redirect(route string) tea.Cmd {
	return CreateMsg(RouteMsg{Route: route})
}

func RedirectHome() tea.Cmd {
	return Redirect(teacommon.HomeRoute)
}

func RedirectSearch() tea.Cmd {
	return Redirect(teacommon.SearchRoute)
}

func RedirectAuth() tea.Cmd {
	return Redirect(teacommon.AuthRoute)
}

func NotificationRedirect() tea.Cmd {
	return CreateMsg(ViewCreatedMsg{})
}

func UpdateHelpKeys(keys []*key.Binding) tea.Cmd {
	return CreateMsg(HelpMsg{keys})
}

func HandleSearch(keyword string, searchType string, page int) tea.Cmd {
	return CreateMsg(SearchMsg{
		Keyword: keyword,
		Type:    searchType,
		Page:    page,
	})
}

func HandleSearchResult(result teacommon.SearchResult, err error) tea.Cmd {
	return CreateMsg(SearchResultMsg{
		Result: result,
		Error:  err,
	})
}

func StartDownload(channel teacommon.Channel) tea.Cmd {
	return CreateMsg(DownloadMsg{
		Channel: channel,
	})
}

func Tick() tea.Cmd {
	return tea.Tick(time.Second, func(time.Time) tea.Msg {
		return TickMsg{}
	})
}
