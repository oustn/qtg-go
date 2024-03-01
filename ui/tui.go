package ui

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/muesli/reflow/indent"
	"github.com/oustn/qtg/internal/api"
	"github.com/oustn/qtg/internal/config"
	uiauth "github.com/oustn/qtg/ui/auth"
	uicommon "github.com/oustn/qtg/ui/common"
	uihome "github.com/oustn/qtg/ui/home"
	uisearch "github.com/oustn/qtg/ui/search"
)

type Q struct {
	keys       uicommon.KeyMap
	help       help.Model
	config     *config.Config
	configPath string
	api        *api.QingTingApi
	route      map[string]uicommon.InitView
	view       uicommon.View
	init       bool
}

type (
	created struct{}
)

func NewQ(cfg *config.Config, cfgPath string) Q {
	qtApi := api.InitQingTingApi(cfg.Settings.RefreshToken, cfg.Settings.QingTingId)
	_ = qtApi.Auth()

	q := Q{
		keys: uicommon.KeyMap{
			Help: key.NewBinding(
				key.WithKeys("?"),
				key.WithHelp("?", "帮助"),
			),
			Quit: key.NewBinding(
				key.WithKeys("esc", "ctrl+c"),
				key.WithHelp("esc", "返回/退出"),
			),
		},
		help:       help.New(),
		config:     cfg,
		configPath: cfgPath,
		api:        qtApi,
		route:      make(map[string]uicommon.InitView),
		view:       nil,
	}

	q.AddRoute(uicommon.HomeRoute, uihome.InitHome)
	q.AddRoute(uicommon.AuthRoute, uiauth.InitAuth)
	q.AddRoute(uicommon.SearchRoute, uisearch.InitSearch)

	return q
}

func (q Q) Init() tea.Cmd {
	return tea.Batch(tea.EnterAltScreen, uicommon.CreateMessageImmediate(created{}))
}

func (q Q) View() string {
	var s string
	if q.view == nil {
		s += "Hello world"
	} else {
		s += q.renderBanner()
		if q.view.ShouldRenderUserInfo() {
			s += q.renderUserInfo() + "\n\n"
		}
		s += q.view.Render()
	}
	s += q.renderHelper()
	return indent.String("\n"+s+"\n\n", 2)
}

func (q Q) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		q.help.Width = msg.Width

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, q.keys.Help):
			q.help.ShowAll = !q.help.ShowAll
		case key.Matches(msg, q.keys.Quit):
			switch q.view.(type) {
			case *uihome.Home:
				return q, tea.Quit
			default:
				return q.Home()
			}
		}
	}

	switch msg.(type) {
	case created:
		return q.Home()
	case uicommon.ReplaceRouteMsg:
		return q.Replace(msg.(uicommon.ReplaceRouteMsg).Name)
	case uicommon.AuthSuccessMsg:
		return q.AuthSuccess(msg.(uicommon.AuthSuccessMsg))
	}
	if q.view != nil {
		return q, q.view.Handle(msg)
	}
	return q, nil
}

func (q Q) AddRoute(route string, initView uicommon.InitView) {
	q.route[route] = initView
}

func (q Q) Home() (tea.Model, tea.Cmd) {
	return q.Replace(uicommon.HomeRoute)
}

func (q Q) Replace(route string) (tea.Model, tea.Cmd) {
	if viewInit, ok := q.route[route]; ok {
		q.view = viewInit()
		cmd := uicommon.CreateMessageImmediate(uicommon.ViewPayload{
			UserInfo:     &(q.api.User),
			RefreshToken: q.api.RefreshToken,
			QingTingId:   q.api.QingTingId,
			Created:      q.init,
		})
		q.init = true
		q.keys = q.keys.BindDynamicKeys(q.view.KeyBinds())
		return q, cmd
	}
	return q, nil
}

func (q Q) AuthSuccess(msg uicommon.AuthSuccessMsg) (tea.Model, tea.Cmd) {
	if msg.Id != q.api.QingTingId || msg.Token != q.api.RefreshToken {
		q.api.QingTingId = msg.Id
		q.api.RefreshToken = msg.Token
		q.config.Settings.QingTingId = msg.Id
		q.config.Settings.RefreshToken = msg.Token

		_ = q.api.Auth(msg.Token, msg.Id)
		_ = config.WriteConfig(q.config)
	}
	return q.Home()
}

func (q Q) renderBanner() string {
	return uicommon.InfoStyle.
		Italic(true).
		Bold(true).
		Underline(true).
		Render("QTG - 一个有趣的蜻蜓FM下载器") + "\n\n"
}

func (q Q) renderUserInfo() string {
	if q.api.RefreshToken == "" || q.api.QingTingId == "" {
		return uicommon.BlurredTextStyle("未配置")
	}
	if q.api.User == (uicommon.UserInfo{}) {
		return uicommon.ErrorTextStyle("认证失败")
	}
	return uicommon.SuccessTextStyle("你好，" + q.api.User.Nickname + "")
}

func (q Q) renderFooter() {}

func (q Q) renderHelper() string {
	helpView := q.help.View(q.keys)
	return "\n\n" + helpView
}
