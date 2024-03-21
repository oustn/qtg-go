package teaui

import (
	"fmt"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/oustn/qtg/internal/api"
	"github.com/oustn/qtg/internal/config"
	"github.com/oustn/qtg/internal/download"
	teaauth "github.com/oustn/qtg/internal/ui/auth"
	teacmd "github.com/oustn/qtg/internal/ui/cmd"
	teacommon "github.com/oustn/qtg/internal/ui/common"
	teacomponents "github.com/oustn/qtg/internal/ui/components"
	teahome "github.com/oustn/qtg/internal/ui/home"
	teasearch "github.com/oustn/qtg/internal/ui/search"
	"strings"
)

var (
	appStyle  = lipgloss.NewStyle().Padding(1, 2)
	logoText  = fmt.Sprintf("%s %s", "🐲", "QTG")
	maxHeight = 32
)

func newStatusBar() teacomponents.StatusBar {
	bar := teacomponents.NewStatusBar(4)
	bar.SetColumns([]teacomponents.Column{
		{
			Content:      "QTG",
			Foreground:   teacommon.StatusBarLogoForegroundColor,
			Background:   teacommon.StatusBarLogoBackgroundColor,
			ContentWidth: 30,
		},
		{
			Content:    "一个有趣的蜻蜓FM下载工具",
			Foreground: teacommon.StatusBarBarForegroundColor,
			Background: teacommon.StatusBarBarBackgroundColor,
			Flex:       true,
		},
		{
			Foreground: teacommon.StatusBarTotalFilesForegroundColor,
			Background: teacommon.StatusBarTotalFilesBackgroundColor,
		},
		{
			Content:    logoText,
			Foreground: teacommon.StatusBarLogoForegroundColor,
			Background: teacommon.StatusBarLogoBackgroundColor,
		},
	})
	return bar
}

func NewRenderer(cfg *config.Config, cfgPath string) tea.Model {
	qtApi := api.InitQingTingApi(cfg.Settings.RefreshToken, cfg.Settings.QingTingId)
	_ = qtApi.Auth()

	helper := teacomponents.NewHelper()
	helper.Styles.TitleStyle = lipgloss.NewStyle().
		Foreground(teacommon.AccentColor).
		Bold(true).
		Italic(true)

	downloader := download.NewDownloader(qtApi, 1, 5)

	r := renderer{
		cfg:        cfg,
		cfgPath:    cfgPath,
		api:        qtApi,
		statusbar:  newStatusBar(),
		helper:     helper,
		banner:     teacomponents.NewBanner(),
		downloader: downloader,
	}

	downloader.Callback = r.updateProgress

	return r
}

type (
	renderer struct {
		route       *teacommon.Route
		cfg         *config.Config
		cfgPath     string
		api         *api.QingTingApi
		height      int
		width       int
		statusbar   teacomponents.StatusBar
		helper      teacomponents.Helper
		banner      teacommon.Component
		view        teacommon.View
		downloader  *download.Downloader
		downloading bool
	}

	rendererCreate struct{}
)

func (r renderer) Init() tea.Cmd {
	return tea.Batch(tea.EnterAltScreen, teacmd.CreateMsg(rendererCreate{}))
}

func (r renderer) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h, v := appStyle.GetFrameSize()
		r.height = msg.Height - v
		r.width = msg.Width - h

		msg.Width = r.width
		msg.Height = r.height

		r.adjustDimensions(r.width, r.height)
	case rendererCreate:
		return r, tea.Batch(teacmd.Auth(r.api.RefreshToken, r.api.QingTingId), teacmd.RedirectHome())
	case teacmd.UpdateStatusBarMsg:
		r.updateStatusBar(msg.Title, msg.Msg)
		return r, nil
	case teacmd.AuthMsg:
		return r, r.auth(msg.Token, msg.Id, msg.Callback)
	case teacmd.AuthSuccessMsg:
		return r, r.updateAuthResult(nil)
	case teacmd.AuthFailMsg:
		return r, r.updateAuthResult(msg.Err)
	case teacmd.UpdateQuitKeyMsg:
		r.helper.SetQQuitEnable(msg.Q)
		r.helper.SetEscQuitEnable(msg.Esc)
		return r, nil
	case teacmd.HelpMsg:
		r.helper.SetKeys(msg.Keys)
		return r, nil
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, r.helper.HelpKey):
			return r, r.helper.ToggleShowAll()
		case key.Matches(msg, r.helper.ForceQuitKey):
			return r, tea.Quit
		case key.Matches(msg, r.helper.QQuitKey, r.helper.EscQuitKey):
			if r.route != &teahome.HomeRoute {
				return r, teacmd.RedirectHome()
			}
			return r, tea.Quit
		}
	case teacmd.RouteMsg:
		return r, r.createView(msg.Route)
	case teacmd.SearchMsg:
		return r, r.handleSearch(msg.Keyword, msg.Type, msg.Page)

	case teacmd.DownloadMsg:
		return r, r.download(msg.Channel)
	case teacmd.TickMsg:
		return r, teacmd.Tick()
	}

	if r.view != nil {
		return r, r.view.Handler(msg)
	}
	return r, nil
}

func (r renderer) View() string {
	var (
		sections        []string
		availableHeight = r.height
	)

	header := r.renderHeader()
	availableHeight -= lipgloss.Height(header)
	sections = append(sections, header)

	statusBar := r.statusbar.Render()
	availableHeight -= lipgloss.Height(statusBar)

	helper := lipgloss.NewStyle().MarginBottom(1).Render(r.helper.Render())
	availableHeight -= lipgloss.Height(helper)

	sections = append(sections, r.renderContent(availableHeight))
	sections = append(sections, helper)
	sections = append(sections, statusBar)

	return appStyle.Render(lipgloss.JoinVertical(lipgloss.Top, sections...))
}

func (r *renderer) renderContent(height int) string {
	h := min(height, maxHeight)
	if r.view != nil {
		r.view.SetSize(r.width, h)
		return lipgloss.NewStyle().Height(h).Render(r.view.Render())
	}
	return lipgloss.NewStyle().Height(min(height, maxHeight)).Render("hello world")
}

func (r *renderer) renderHeader() string {
	var sections []string
	sections = append(sections, r.banner.Render())
	bread := []string{
		teacommon.InfoStyle.Bold(true).Render(teahome.HomeRoute.Title),
	}
	if r.route != nil && r.route != &teahome.HomeRoute {
		bread = []string{
			teacommon.BlurredTextStyle(teahome.HomeRoute.Title),
			teacommon.InfoStyle.Bold(true).Render(r.route.Title),
		}
	}
	sections = append(sections, strings.Join(bread, teacommon.BlurredTextStyle(" > ")))
	return lipgloss.NewStyle().
		Width(r.width).
		MarginBottom(1).
		BorderBottom(true).
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(teacommon.SecondaryTextColor).
		Render(lipgloss.JoinVertical(lipgloss.Top, sections...))
}

func (r *renderer) adjustDimensions(width int, height int) {
	r.statusbar.SetSize(width, height)
	r.helper.SetSize(width, height)
}

func (r *renderer) updateStatusBar(title, msg string) {
	r.statusbar.UpdateColumnContent(0, title)
	r.statusbar.UpdateColumnContent(1, msg)
}

func (r *renderer) resetStatusBar() {
	r.updateStatusBar("QTG", "一个有趣的蜻蜓FM下载工具")
}

func (r *renderer) updateAuthInfo(info string) {
	r.statusbar.UpdateColumnContent(2, info)
}

func (r *renderer) auth(token, id string, callback teacommon.Callback) tea.Cmd {
	if token == "" || id == "" {
		callback(fmt.Errorf("未配置Token"))
		return teacmd.CreateMsg(teacmd.AuthFailMsg{
			Err: fmt.Errorf("未配置Token"),
		})
	}

	err := r.api.Auth(token, id)

	if err != nil {
		callback(err)
		return teacmd.CreateMsg(teacmd.AuthFailMsg{Err: err})
	}
	r.cfg.Settings.RefreshToken = token
	r.cfg.Settings.QingTingId = id
	_ = config.WriteConfig(r.cfg)

	callback(nil)
	return teacmd.CreateMsg(teacmd.AuthSuccessMsg{})
}

func (r *renderer) updateAuthResult(err error) tea.Cmd {
	var msg string
	if err == nil {
		if r.api.User.PrivateInfo.VipInfo.Vip {
			msg = fmt.Sprintf("%s%s", r.api.User.Nickname, lipgloss.NewStyle().Italic(true).Bold(true).Foreground(teacommon.WarnColor).Render(" ✨"))
		} else {
			msg = r.api.User.Nickname
		}
	} else {
		msg = err.Error()
	}
	r.updateAuthInfo(msg)
	return nil
}

func (r *renderer) createView(route string) tea.Cmd {
	switch route {
	case teacommon.HomeRoute:
		r.route = &teahome.HomeRoute
	case teacommon.SearchRoute:
		r.route = &teasearch.SearchRoute
	case teacommon.AuthRoute:
		r.route = &teaauth.AuthRoute
	}
	r.helper.ClearKeys()
	if r.route != nil {
		r.view = r.route.RenderComponent()
		return tea.Batch(teacmd.NotificationRedirect(), r.view.Init())
	}
	return nil
}

func (r *renderer) handleSearch(keyword, searchType string, page int) tea.Cmd {
	if keyword == "" {
		return nil
	}
	result, err := r.api.Search(keyword, searchType, page)
	return teacmd.HandleSearchResult(result, err)
}

func (r *renderer) download(channel teacommon.Channel) tea.Cmd {
	r.downloader.DownloadChannel(channel)
	return teacmd.Tick()
}

func (r *renderer) updateProgress(progress download.ChannelProgress) {
	d := progress.Downloading.Cardinality()
	f := progress.Finished.Cardinality()
	t := progress.Total()
	if d == 0 {
		r.resetStatusBar()
		return
	}
	current := progress.Downloading.ToSlice()[0]
	channel, ok := progress.Channels[current.(string)]
	if !ok {
		r.resetStatusBar()
		return
	}

	pt := channel.Programs.Total()

	desc := ""
	if channel.Programs.Finished != nil {
		// fmt.Println(channel.Programs.Finished.ToSlice(), channel.Programs.Pending.ToSlice(), channel.Programs.Downloading.ToSlice())
		pf := channel.Programs.Finished.Cardinality()
		desc = fmt.Sprintf("%s(%d/%d)", channel.Name, pf, pt)
	}

	r.updateStatusBar(
		fmt.Sprintf("下载中(%d/%d)", f, t),
		desc,
	)
}
